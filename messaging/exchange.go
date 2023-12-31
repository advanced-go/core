package messaging

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"sort"
	"sync"
)

const (
	exSendCtrlLocation = PkgPath + ":Exchange/SendCtrl"
	exSendDataLocation = PkgPath + ":Exchange/SendData"
	exGetLocation      = PkgPath + ":Exchange/get"
	exAddLocation      = PkgPath + ":Exchange/add"
)

// HostExchange - main exchange
var HostExchange = NewExchange()

func shutdownHost(msg Message) runtime.Status {
	//TO DO: authentication and implementation
	return runtime.StatusOK()
}

// Exchange - exchange directory
type Exchange interface {
	Count() int
	List() []string
	Add(m *Mailbox) runtime.Status
	SendCtrl(msg Message) runtime.Status
	SendData(msg Message) runtime.Status
}

type exchange struct {
	m *sync.Map
}

// NewExchange - create a new exchange
func NewExchange() Exchange {
	e := new(exchange)
	e.m = new(sync.Map)
	return e
}

// Count - number of items in the sync map
func (d *exchange) Count() int {
	count := 0
	d.m.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}

// List - a list of item uri's
func (d *exchange) List() []string {
	var uri []string
	d.m.Range(func(key, value any) bool {
		if str, ok := key.(string); ok {
			uri = append(uri, str)
		}
		return true
	})
	sort.Strings(uri)
	return uri
}

// SendCtrl - send a message to the select item's control channel
func (d *exchange) SendCtrl(msg Message) runtime.Status {
	// TO DO : authenticate shutdown control message
	if msg.Event == ShutdownEvent {
		return runtime.StatusOK()
	}
	mbox, status := d.get(msg.To)
	if !status.OK() {
		return status.AddLocation(exSendCtrlLocation)
	}
	mbox.SendCtrl(msg)
	return runtime.StatusOK()
}

// SendData - send a message to the item's data channel
func (d *exchange) SendData(msg Message) runtime.Status {
	mbox, status := d.get(msg.To)
	if !status.OK() {
		return status.AddLocation(exSendDataLocation)
	}
	mbox.SendData(msg)
	return runtime.StatusOK()
}

// Add - add a mailbox
func (d *exchange) Add(m *Mailbox) runtime.Status {
	if m == nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, exAddLocation, errors.New("invalid argument: mailbox is nil"))
	}
	if len(m.uri) == 0 {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, exAddLocation, errors.New("invalid argument: mailbox uri is empty"))
	}
	if m.ctrl == nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, exAddLocation, errors.New("invalid argument: mailbox command channel is nil"))
	}
	_, ok := d.m.Load(m.uri)
	if ok {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, exAddLocation, errors.New(fmt.Sprintf("invalid argument: exchange mailbox already exists: [%v]", m.uri)))
	}
	d.m.Store(m.uri, m)
	m.unregister = func() {
		d.m.Delete(m.uri)
	}
	return runtime.StatusOK()
}

func (d *exchange) get(uri string) (*Mailbox, runtime.Status) {
	if len(uri) == 0 {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, exGetLocation, errors.New("invalid argument: uri is empty"))
	}
	v, ok1 := d.m.Load(uri)
	if !ok1 {
		return nil, runtime.NewStatusError(http.StatusNotFound, exGetLocation, errors.New(fmt.Sprintf("invalid URI: exchange mailbox not found [%v]", uri)))
	}
	if mbox, ok2 := v.(*Mailbox); ok2 {
		return mbox, runtime.StatusOK()
	}
	return nil, runtime.NewStatusError(runtime.StatusInvalidContent, exGetLocation, errors.New("invalid Mailbox type"))
}

// Shutdown - close an item's mailbox
func (d *exchange) Shutdown(msg Message) runtime.Status {
	// TO DO: add authentication
	return runtime.StatusOK() //d.shutdown(msg.To)
}

/*
func (d *exchange) shutdown(uri string) runtime.Status {
	//d.mu.RLock()
	//defer d.mu.RUnlock()
	//for _, e := range d.m {
	//	if e.ctrl != nil {
	//		e.ctrl <- Message{To: e.uri, Event: core.ShutdownEvent}
	//	}
	//}
	m, status := d.get(uri)
	if !status.OK() {
		return status
	}
	if m.data != nil {
		close(m.data)
	}
	if m.ctrl != nil {
		close(m.ctrl)
	}
	d.m.Delete(uri)
	return runtime.StatusOK()
}
*/
