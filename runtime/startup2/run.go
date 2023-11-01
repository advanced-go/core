package startup2

import (
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
)

type messageMap map[string]Message

var (
	runLocation = PkgUri + "/Run"
	directory   = NewEntryDirectory()
)

// Register - function to register a startup uri
func Register(uri string, h runtime.TypeHandlerFn) error {
	if uri == "" {
		return errors.New("invalid argument: uri is empty")
	}
	if h == nil {
		return errors.New(fmt.Sprintf("invalid argument: type handler is nil for [%v]", uri))
	}
	registerUnchecked(uri, h)
	return nil
}

func registerUnchecked(uri string, h runtime.TypeHandlerFn) error {
	directory.Add(uri, h)
	return nil
}

// Shutdown - startup shutdown
func Shutdown() {
	directory.Shutdown()
}
