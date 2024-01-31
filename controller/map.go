package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"sync"
)

const (
	mapAdd = PkgPath + ":Add"
	mapGet = PkgPath + ":Get"
	mapNew = PkgPath + ":NewMap"

	errorKey = "error"
)

// Map - key value pairs of string -> string
type Map struct {
	m *sync.Map
}

func NewEmptyMap() *Map {
	m := new(Map)
	m.m = new(sync.Map)
	return m
}

func NewMap(buf []byte) (*Map, runtime.Status) {
	var ctrl []Config
	err := json.Unmarshal(buf, &ctrl)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusJsonDecodeError, mapNew, errors.New(fmt.Sprintf("JSON decode error: [%v]", err)))
	}
	m := NewEmptyMap()
	for _, cfg := range ctrl {
		c := new(Controller)
		c.Uri = cfg.Uri
		c.Name = cfg.Name
		c.Route = cfg.Route
		c.Method = cfg.Method
		c.Duration, err = ParseDuration(cfg.Duration)
		if err != nil {
			return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, mapNew, errors.New(fmt.Sprintf("duration configuration error: [%v]", err)))
		}
		m.Add(c.Name, c)
	}
	return m, runtime.StatusOK()
}

// Add - add a controller
func (m *Map) Add(key string, ctrl *Controller) runtime.Status {
	if len(key) == 0 {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, mapAdd, errors.New("invalid argument: key is empty"))
	}
	_, ok1 := m.m.Load(key)
	if ok1 {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, mapAdd, errors.New(fmt.Sprintf("invalid argument: key already exists: [%v]", key)))
	}
	m.m.Store(key, ctrl)
	return runtime.StatusOK()
}

// Get - get a controller
func (m *Map) Get(key string) (*Controller, runtime.Status) {
	v, ok := m.m.Load(key)
	if !ok {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, mapGet, errors.New(fmt.Sprintf("invalid argument: key does not exist: [%v]", key)))
	}
	if val, ok1 := v.(*Controller); ok1 {
		return val, runtime.StatusOK()
	}
	return nil, runtime.NewStatus(runtime.StatusInvalidContent)
}
