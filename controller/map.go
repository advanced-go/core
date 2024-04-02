package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Controller2 struct {
	Name string `json:"name"`
	//Route    string        `json:"route"`
	//Method   string        `json:"method"`
	//Uri      string        `json:"uri"`
	DurationS string `json:"duration"`
	Duration  time.Duration
}

// Map - key value pairs of string -> string
type Map struct {
	m *sync.Map
}

func NewEmptyMap() *Map {
	m := new(Map)
	m.m = new(sync.Map)
	return m
}

func NewMap(buf []byte) (*Map, error) {
	var ctrl []Controller2
	err := json.Unmarshal(buf, &ctrl)
	if err != nil {
		return nil, err
	}
	m := NewEmptyMap()
	for _, cfg := range ctrl {
		c := new(Controller2)
		c.Name = cfg.Name

		c.Duration, err = ParseDuration(cfg.DurationS)
		if err != nil {
			return nil, err
		}
		m.Add(c)
	}
	return m, nil
}

// Add - add a controller
func (m *Map) Add(ctrl *Controller2) error {
	if ctrl == nil {
		return errors.New("invalid argument: controller is nil")
	}
	if len(ctrl.Name) == 0 {
		return errors.New("invalid argument: key is empty")
	}
	_, ok1 := m.m.Load(ctrl.Name)
	if ok1 {
		return errors.New(fmt.Sprintf("invalid argument: key already exists: [%v]", ctrl.Name))
	}
	m.m.Store(ctrl.Name, ctrl)
	return nil
}

// Get - get a controller
func (m *Map) Get(key string) (*Controller2, error) {
	v, ok := m.m.Load(key)
	if !ok {
		return nil, errors.New(fmt.Sprintf("invalid argument: key does not exist: [%v]", key))
	}
	if val, ok1 := v.(*Controller2); ok1 {
		return val, nil
	}
	return nil, errors.New(fmt.Sprintf("invalid argument: content invalid: [%v]", key))
}
