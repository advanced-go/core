package startup2

import (
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"sort"
	"sync"
)

// Entry - and entry in an EntryDirectory
type Entry struct {
	uri     string
	handler runtime.TypeHandlerFn
}

// EntryDirectory - collection of Entry
type EntryDirectory struct {
	m  map[string]*Entry
	mu sync.RWMutex
}

// NewEntryDirectory - create a new directory
func NewEntryDirectory() *EntryDirectory {
	return &EntryDirectory{m: make(map[string]*Entry)}
}

func (d *EntryDirectory) Get(uri string) *Entry {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.m[uri]
}

func (d *EntryDirectory) Add(uri string, handler runtime.TypeHandlerFn) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.m[uri] = &Entry{
		uri:     uri,
		handler: handler,
	}
}

func (d *EntryDirectory) Count() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.m)
}

func (d *EntryDirectory) Uri() []string {
	var uri []string
	d.mu.RLock()
	defer d.mu.RUnlock()
	for key, _ := range d.m {
		uri = append(uri, key)
	}
	sort.Strings(uri)
	return uri
}

func (d *EntryDirectory) Send(msg Message) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if e, ok := d.m[msg.To]; ok {
		if e.handler == nil {
			return errors.New(fmt.Sprintf("handler function is nil: [%v]", msg.To))
		}
		req, _ := http.NewRequest("", "/runtime", nil)
		_, status := e.handler(req, msg)
		if status != nil {
		}
		//e.c <- msg
		return nil
	}
	return errors.New(fmt.Sprintf("entry not found: [%v]", msg.To))
}

func (d *EntryDirectory) Shutdown() {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, e := range d.m {
		if e.handler != nil {
			req, _ := http.NewRequest("", "/runtime", nil)
			_, _ = e.handler(req, Message{To: e.uri, Event: ShutdownEvent})
			//	e.c <- Message{To: e.uri, Event: ShutdownEvent}
		}
	}
}

func (d *EntryDirectory) Empty() {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for key, _ := range d.m {
		//if e.handler != nil {
		//	close(e.hanc)
		//}
		delete(d.m, key)
	}
}
