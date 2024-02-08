package messaging

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// MessageCache - message cache by uri
type MessageCache struct {
	m    map[string]*Message
	errs []error
	mu   sync.RWMutex
}

// NewMessageCache - create a message cache
func NewMessageCache() *MessageCache {
	c := new(MessageCache)
	c.m = make(map[string]*Message)
	return c
}

// Count - return the count of items
func (r *MessageCache) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, _ = range r.m {
		count++
	}
	return count
}

// Filter - apply a filter against a traversal of all items
func (r *MessageCache) Filter(event string, code int, include bool) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var uri []string
	for u, resp := range r.m {
		if include {
			if resp.Status.Code == code && resp.Event() == event {
				uri = append(uri, u)
			}
		} else {
			if resp.Status.Code != code || resp.Event() != event {
				uri = append(uri, u)
			}
		}
	}
	sort.Strings(uri)
	return uri
}

// Include - filter for items that include a specific event
func (r *MessageCache) Include(event string, status int) []string {
	return r.Filter(event, status, true)
}

// Exclude - filter for items that do not include a specific event
func (r *MessageCache) Exclude(event string, status int) []string {
	return r.Filter(event, status, false)
}

// Add - add a message
func (r *MessageCache) Add(msg *Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if msg.From() == "" {
		err := errors.New("invalid argument: message from is empty")
		r.errs = append(r.errs, err)
		return err
	}
	if _, ok := r.m[msg.From()]; !ok {
		r.m[msg.From()] = msg
		return nil
	}
	err0 := errors.New(fmt.Sprintf("invalid argument: message found [%v]", msg.From()))
	r.errs = append(r.errs, err0)
	return err0
}

// Get - get a message based on a URI
func (r *MessageCache) Get(uri string) (*Message, bool) {
	if uri == "" {
		return nil, false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.m[uri]; ok {
		return r.m[uri], true
	}
	return nil, false //errors.New(fmt.Sprintf("invalid argument: uri not found [%v]", uri))
}

// Uri - list the URI's in the cache
func (r *MessageCache) Uri() []string {
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
func (r *MessageCache) ErrorList() []error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.errs
}

// NewMessageCacheHandler - handler to receive messages into a cache.
func NewMessageCacheHandler(cache *MessageCache) MessageHandler {
	return func(msg *Message) {
		err := cache.Add(msg)
		if err != nil {
			fmt.Printf("error: messaging cache handler cache.Add() %v\n", err)
		}
	}
}
