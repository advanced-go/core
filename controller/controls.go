package controller

import (
	"errors"
	"fmt"

	"github.com/advanced-go/core/runtime"
	"sync"
)

var (
	ctrlMap = NewControls()
)

// RegisterController - add a controller for a URI
func RegisterController(uri string, ctrl *Controller) *runtime.Status {
	return ctrlMap.register(uri, ctrl)
}

func Lookup(uri string) (*Controller, *runtime.Status) {
	nid, _, ok := UprootUrn(uri)
	if !ok {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New(fmt.Sprintf("invalid argument: path is invalid: [%v]", uri)))
	}
	return ctrlMap.lookupByNID(nid)
}

// controls - key value pairs of a URI -> *Controller
type controls struct {
	m *sync.Map
}

// NewControls - create a new Controls map
func NewControls() *controls {
	p := new(controls)
	p.m = new(sync.Map)
	return p
}

// Register - add a controller
func (p *controls) register(uri string, ctrl *Controller) *runtime.Status {
	if len(uri) == 0 {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New("invalid argument: path is empty"))
	}
	nid, _, ok := UprootUrn(uri)
	if !ok {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New(fmt.Sprintf("invalid argument: path is invalid: [%v]", uri)))
	}
	if ctrl == nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New(fmt.Sprintf("invalid argument: Controller is nil: [%v]", uri)))
	}
	_, ok1 := p.m.Load(nid)
	if ok1 {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New(fmt.Sprintf("invalid argument: Controller already exists: [%v]", uri)))
	}
	p.m.Store(nid, ctrl)
	return runtime.StatusOK()
}

// Lookup - get a Controller using a URI as the key
func (p *controls) lookup(uri string) (*Controller, *runtime.Status) {
	nid, _, ok := UprootUrn(uri)
	if !ok {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New(fmt.Sprintf("invalid argument: path is invalid: [%v]", uri)))
	}
	return p.lookupByNID(nid)
}

// LookupByNID - get a Controller using an NID as a key
func (p *controls) lookupByNID(nid string) (*Controller, *runtime.Status) {
	v, ok := p.m.Load(nid)
	if !ok {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New(fmt.Sprintf("invalid argument: Controller does not exist: [%v]", nid)))
	}
	if ctrl, ok1 := v.(*Controller); ok1 {
		return ctrl, runtime.StatusOK()
	}
	return nil, runtime.NewStatus(runtime.StatusInvalidContent)
}
