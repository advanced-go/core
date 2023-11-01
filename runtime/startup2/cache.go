package startup2

import (
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"sort"
	"sync"
)

type StatusMessage struct {
	Msg    Message
	Status *runtime.Status
}

// MessageCache - message cache by startup uri
type MessageCache struct {
	m  map[string]StatusMessage
	mu sync.RWMutex
}

// NewMessageCache - create a message cache
func NewMessageCache() *MessageCache {
	return &MessageCache{m: make(map[string]StatusMessage)}
}

func (r *MessageCache) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, _ = range r.m {
		count++
	}
	return count
}

func (r *MessageCache) Filter(event string, code int, include bool) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var uri []string
	for u, resp := range r.m {
		if include {
			if resp.Status != nil && resp.Status.Code() == code && resp.Msg.Event == event {
				uri = append(uri, u)
			}
		} else {
			if resp.Status != nil && resp.Status.Code() != code || resp.Msg.Event != event {
				uri = append(uri, u)
			}
		}
	}
	sort.Strings(uri)
	return uri
}

func (r *MessageCache) Include(event string, status int) []string {
	return r.Filter(event, status, true)
}

func (r *MessageCache) Exclude(event string, status int) []string {
	return r.Filter(event, status, false)
}

func (r *MessageCache) Add(msg Message, status *runtime.Status) error {
	if msg.From == "" {
		return errors.New("invalid argument: message from is empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.m[msg.From]; !ok {
		r.m[msg.From] = StatusMessage{Msg: msg, Status: status}
	}
	return nil
}

func (r *MessageCache) Get(uri string) (StatusMessage, error) {
	if uri == "" {
		return StatusMessage{}, errors.New("invalid argument: uri is empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.m[uri]; ok {
		return r.m[uri], nil
	}
	return StatusMessage{}, errors.New(fmt.Sprintf("invalid argument: uri not found [%v]", uri))
}

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
