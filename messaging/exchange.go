package messaging

import (
	"errors"
	"fmt"
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

// Exchange - exchange directory
type Exchange interface {
	Count() int
	List() []string
	Add(m *Mailbox) error
	SendCtrl(msg Message) error
	SendData(msg Message) error
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
func (d *exchange) SendCtrl(msg Message) error {
	// TO DO : authenticate shutdown control message
	if msg.Event == ShutdownEvent {
		return nil
	}
	mbox := d.get(msg.To)
	if mbox == nil {
		return errors.New(fmt.Sprintf("error: exchange.SendCtrl() failed as the message To is empty or invalid [%v]", msg.To))
	}
	mbox.SendCtrl(msg)
	return nil
}

// SendData - send a message to the item's data channel
func (d *exchange) SendData(msg Message) error {
	mbox := d.get(msg.To)
	if mbox == nil {
		return errors.New(fmt.Sprintf("error: exchange.SendCtrl() failed as the message To is empty or invalid [%v]", msg.To))
	}
	mbox.SendData(msg)
	return nil
}

// Add - add a mailbox
func (d *exchange) Add(m *Mailbox) error {
	if m == nil {
		return errors.New("error: exchange.Add() mailbox is nil")
	}
	if len(m.uri) == 0 {
		return errors.New("error: exchange.Add() mailbox uri is empty")
	}
	if m.ctrl == nil {
		return errors.New("error: exchange.Add() mailbox command channel is nil")
	}
	_, ok := d.m.Load(m.uri)
	if ok {
		return errors.New(fmt.Sprintf("error: exchange.Add() mailbox already exists: [%v]", m.uri))
	}
	d.m.Store(m.uri, m)
	m.unregister = func() {
		d.m.Delete(m.uri)
	}
	return nil
}

func (d *exchange) get(uri string) *Mailbox {
	if len(uri) == 0 {
		return nil
	}
	v, ok1 := d.m.Load(uri)
	if !ok1 {
		return nil
	}
	if mbox, ok2 := v.(*Mailbox); ok2 {
		return mbox
	}
	return nil
}

// Shutdown - close an item's mailbox
func (d *exchange) Shutdown(msg Message) error {
	// TO DO: add authentication
	return nil
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
