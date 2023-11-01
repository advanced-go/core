package startup2

import (
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"time"
)

type messageMap map[string]Message

var (
	directory = NewEntryDirectory()
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

// Run - templated function to start all registered resources.
func Run[E runtime.ErrorHandler](duration time.Duration, content ContentMap) (status *runtime.Status) {
	var e E
	var failures []string
	var count = directory.Count()

	if count == 0 {
		return runtime.NewStatusOK()
	}
	cache := NewMessageCache()
	toSend := createToSend(content)
	sendMessages(toSend)
	for wait := time.Duration(float64(duration) * 0.25); duration >= 0; duration -= wait {
		time.Sleep(wait)
		// Check for completion
		if cache.Count() < count {
			continue
		}
		// Check for failed resources
		failures = cache.Exclude(StartupEvent, http.StatusOK)
		if len(failures) == 0 {
			handleStatus(cache)
			return runtime.NewStatusOK()
		}
		break
	}
	Shutdown()
	if len(failures) > 0 {
		handleErrors[E](failures, cache)
		return runtime.NewStatus(http.StatusInternalServerError)
	}
	return e.Handle("", runLocation, errors.New(fmt.Sprintf("response counts < directory entries [%v] [%v]", cache.Count(), directory.Count()))).SetCode(runtime.StatusDeadlineExceeded)
}

func createToSend(cm ContentMap) messageMap {
	m := make(messageMap)
	for _, k := range directory.Uri() {
		msg := Message{To: k, From: PkgUri, Event: StartupEvent}
		if cm != nil {
			if content, ok := cm[k]; ok {
				msg.Content = append(msg.Content, content...)
			}
		}
		m[k] = msg
	}
	return m
}

func sendMessages(msgs messageMap) {
	for k := range msgs {
		directory.Send(msgs[k])
	}
}

func handleErrors[E runtime.ErrorHandler](failures []string, cache *MessageCache) {
	var e E
	for _, uri := range failures {
		msg, err := cache.Get(uri)
		if err != nil {
			continue
		}
		if msg.Status != nil {
			e.Handle("", msg.Status.Location()[0], msg.Status.Errors()...)
		}
	}
}

func handleStatus(cache *MessageCache) {
	for _, uri := range cache.Uri() {
		msg, err := cache.Get(uri)
		if err != nil {
			continue
		}
		if msg.Status != nil {
			fmt.Printf("startup successful for startup [%v] : %s", uri, msg.Status.Duration())
		}
	}
}
