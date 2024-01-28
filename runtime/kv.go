package runtime

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const (
	stringsMapAdd = PkgPath + ":Add"
	stringsMapGet = PkgPath + ":Get"
	//handlerLookupNIDLocation = PkgPath + ":LookupFromNID"
)

// KV - key, value pair
type KV struct {
	Key, Value string
}

// Pair - key, value pair
type Pair struct {
	Key, Value string
}

// StringsMap - key value pairs of string -> string
type StringsMap struct {
	m *sync.Map
}

func NewStringsMap(h http.Header) *StringsMap {
	m := new(StringsMap)
	m.m = new(sync.Map)
	if h != nil {
		for k, v := range h {
			if len(v) > 0 {
				m.Add(strings.ToLower(k), v[0])
			}
		}
	}
	return m
}

// Add - add a value
func (m *StringsMap) Add(key, val string) Status {
	if len(key) == 0 {
		return NewStatusError(StatusInvalidArgument, stringsMapAdd, errors.New("invalid argument: key is empty"))
	}
	_, ok1 := m.m.Load(key)
	if ok1 {
		return NewStatusError(StatusInvalidArgument, stringsMapAdd, errors.New(fmt.Sprintf("invalid argument: key already exists: [%v]", key)))
	}
	m.m.Store(key, val)
	return StatusOK()
}

// Get - get a value
func (m *StringsMap) Get(key string) (string, Status) {
	v, ok := m.m.Load(key)
	if !ok {
		return "", NewStatusError(StatusInvalidArgument, stringsMapGet, errors.New(fmt.Sprintf("invalid argument: key does not exist: [%v]", key)))
	}
	if val, ok1 := v.(string); ok1 {
		return val, StatusOK()
	}
	return "", NewStatus(StatusInvalidContent)
}
