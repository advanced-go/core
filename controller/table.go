package controller

import (
	"errors"
	"fmt"
	"sync"
)

// Configuration - configuration for actuators
type Configuration interface {
	SetUriMatcher(fn UriMatcher)
	AddController(route Route) []error
}

// Controllers - public interface
type Controllers interface {
	LookupUri(urn string, method string) Controller
	LookupByName(name string) Controller
}

// Table - controller table
type Table interface {
	Configuration
	Controllers
}

// IngressTable - table for ingress controllers
//var ingressTable = NewIngressTable()

//func IngressTable() Table {
//	return ingressTable
//}

// EgressTable - table for egress controllers
var egressTable = NewEgressTable()

func CtrlTable() Table {
	return egressTable
}

type table struct {
	egress       bool
	allowDefault bool
	mu           sync.RWMutex
	uriMatch     UriMatcher
	defaultCtrl  *controller
	nilCtrl      *controller
	controllers  map[string]*controller
}

// NewEgressTable - create a new Egress table
func NewEgressTable() Table {
	return newTable(true, true)
}

func newTable(egress, allowDefault bool) *table {
	t := new(table)
	t.egress = egress
	t.allowDefault = allowDefault
	t.uriMatch = func(urn string, method string) (name string, ok bool) {
		return "", true
	}
	t.controllers = make(map[string]*controller, 100)
	t.defaultCtrl = newDefaultController(DefaultControllerName)
	t.nilCtrl = newDefaultController(NilControllerName)
	return t
}

func (t *table) isEgress() bool { return t.egress }

func (t *table) SetUriMatcher(fn UriMatcher) {
	if fn == nil {
		return
	}
	t.mu.Lock()
	t.uriMatch = fn
	t.mu.Unlock()
}

func (t *table) LookupUri(uri, method string) Controller {
	name, ok := t.uriMatch(uri, method)
	if !ok {
		return t.nilCtrl
	}
	if name != "" {
		if r := t.LookupByName(name); r != nil {
			return r
		}
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.defaultCtrl
}

func (t *table) LookupByName(name string) Controller {
	if name == "" {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if r, ok := t.controllers[name]; ok {
		return r
	}
	if t.allowDefault {
		return t.defaultCtrl
	}
	return nil
}

func (t *table) AddController(route Route) []error {
	if IsEmpty(route.Name) {
		return []error{errors.New("invalid argument: route name is empty")}
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	ctrl, errs := newController(route, t)
	if len(errs) > 0 {
		return errs
	}
	err := ctrl.validate(t.egress)
	if err != nil {
		return []error{err}
	}
	if _, ok := t.controllers[route.Name]; ok {
		return []error{errors.New(fmt.Sprintf("invalid argument: route name is a duplicate [%v]", route.Name))}
	}
	t.controllers[route.Name] = ctrl
	return nil
}

func (t *table) exists(name string) bool {
	if name == "" {
		return false
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	if _, ok := t.controllers[name]; ok {
		return true
	}
	return false
}

func (c *controller) validate(egress bool) error {
	return nil
}

func (t *table) update(name string, act *controller) {
	if name == "" || act == nil {
		return
	}
	//t.mu.Lock()
	//defer t.mu.Unlock()
	delete(t.controllers, name)
	t.controllers[name] = act
	//return errors.New(fmt.Sprintf("invalid argument : controller not found [%v]", name))
}

func (t *table) count() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.controllers)
}

func (t *table) isEmpty() bool {
	return t.count() == 0
}

func (t *table) remove(name string) {
	if name == "" {
		return
	}
	t.mu.Lock()
	delete(t.controllers, name)
	t.mu.Unlock()
}
