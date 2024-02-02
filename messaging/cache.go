package messaging

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// MessageCache - message cache by uri
type MessageCache interface {
	Count() int
	Filter(event string, code int, include bool) []string
	Include(event string, status int) []string
	Exclude(event string, status int) []string
	Add(msg Message) error
	Get(uri string) (Message, bool)
	Uri() []string
	ErrorList() []error
}

type messageCache struct {
	m    map[string]Message
	errs []error
	mu   sync.RWMutex
}

// NewMessageCache - create a message cache
func NewMessageCache() MessageCache {
	c := new(messageCache)
	c.m = make(map[string]Message)
	return c
}

// Count - return the count of items
func (r *messageCache) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, _ = range r.m {
		count++
	}
	return count
}

// Filter - apply a filter against a traversal of all items
func (r *messageCache) Filter(event string, code int, include bool) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var uri []string
	for u, resp := range r.m {
		if include {
			if resp.Status.Code == code && resp.Event == event {
				uri = append(uri, u)
			}
		} else {
			if resp.Status.Code != code || resp.Event != event {
				uri = append(uri, u)
			}
		}
	}
	sort.Strings(uri)
	return uri
}

// Include - filter for items that include a specific event
func (r *messageCache) Include(event string, status int) []string {
	return r.Filter(event, status, true)
}

// Exclude - filter for items that do not include a specific event
func (r *messageCache) Exclude(event string, status int) []string {
	return r.Filter(event, status, false)
}

// Add - add a message
func (r *messageCache) Add(msg Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if msg.From == "" {
		err := errors.New("invalid argument: message from is empty")
		r.errs = append(r.errs, err)
		return err
	}
	if _, ok := r.m[msg.From]; !ok {
		r.m[msg.From] = msg
		return nil
	}
	err0 := errors.New(fmt.Sprintf("invalid argument: message found [%v]", msg.From))
	r.errs = append(r.errs, err0)
	return err0
}

// Get - get a message based on a URI
func (r *messageCache) Get(uri string) (Message, bool) {
	if uri == "" {
		return Message{}, false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.m[uri]; ok {
		return r.m[uri], true
	}
	return Message{}, false //errors.New(fmt.Sprintf("invalid argument: uri not found [%v]", uri))
}

// Uri - list the URI's in the cache
func (r *messageCache) Uri() []string {
	var uri []string
	r.mu.RLock()
	defer r.mu.RUnlock()
	for key, _ := range r.m {
		uri = append(uri, key)
	}
	sort.Strings(uri)
	return uri
}

// ErrorList - list of errors
func (r *messageCache) ErrorList() []error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.errs
}

// NewMessageCacheHandler - handler to receive messages into a cache.
func NewMessageCacheHandler(cache MessageCache) MessageHandler {
	return func(msg Message) {
		cache.Add(msg)
	}
}
