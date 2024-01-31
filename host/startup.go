package host

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/messaging"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

const (
	startupLocation = PkgPath + ":Startup"
)

// ContentMap - slice of any content to be included in a message
type ContentMap map[string]*runtime.StringsMap

// Startup - templated function to start all registered resources.
func Startup[E runtime.ErrorHandler](duration time.Duration, content ContentMap) (status runtime.Status) {
	return startup[E](messaging.HostExchange, duration, content)
}

func startup[E runtime.ErrorHandler](ex messaging.Exchange, duration time.Duration, content ContentMap) (status runtime.Status) {
	var e E
	var failures []string
	var count = ex.Count()

	if count == 0 {
		return runtime.StatusOK()
	}
	cache := messaging.NewMessageCache()
	toSend := createToSend(ex, content, messaging.NewMessageCacheHandler[E](cache))
	sendMessages(ex, toSend)
	for wait := time.Duration(float64(duration) * 0.25); duration >= 0; duration -= wait {
		time.Sleep(wait)
		// Check for completion
		if cache.Count() < count {
			continue
		}
		// Check for failed resources
		failures = cache.Exclude(messaging.StartupEvent, http.StatusOK)
		if len(failures) == 0 {
			handleStatus(cache)
			return runtime.StatusOK()
		}
		break
	}
	shutdownHost(messaging.Message{Event: messaging.ShutdownEvent})
	if len(failures) > 0 {
		handleErrors[E](failures, cache)
		return runtime.NewStatus(http.StatusInternalServerError)
	}
	return e.Handle(runtime.NewStatusError(runtime.StatusDeadlineExceeded, startupLocation, errors.New(fmt.Sprintf("response counts < directory entries [%v] [%v]", cache.Count(), ex.Count()))), "", "")
}

func createToSend(ex messaging.Exchange, cm ContentMap, fn messaging.MessageHandler) messaging.MessageMap {
	m := make(messaging.MessageMap)
	for _, k := range ex.List() {
		msg := messaging.Message{To: k, From: startupLocation, Event: messaging.StartupEvent, Status: nil, ReplyTo: fn}
		if cm != nil {
			if content, ok := cm[k]; ok {
				//msg.Content = append(msg.Content, content)
				msg.Config = content
			}
		}
		m[k] = msg
	}
	return m
}

func sendMessages(ex messaging.Exchange, msgs messaging.MessageMap) {
	for k := range msgs {
		ex.SendCtrl(msgs[k])
	}
}

func handleErrors[E runtime.ErrorHandler](failures []string, cache messaging.MessageCache) {
	var e E
	for _, uri := range failures {
		msg, err := cache.Get(uri)
		if err != nil {
			continue
		}
		if msg.Status != nil && !msg.Status.OK() {
			loc := ""
			if msg.Status.Location() != nil && len(msg.Status.Location()) > 0 {
				loc = msg.Status.Location()[0]
			}
			e.Handle(runtime.NewStatusError(http.StatusInternalServerError, loc, msg.Status.ErrorList()...), "", "")
		}
	}
}

func handleStatus(cache messaging.MessageCache) {
	for _, uri := range cache.Uri() {
		msg, err := cache.Get(uri)
		if err != nil {
			continue
		}
		if msg.Status != nil {
			fmt.Printf("startup successful: [%v] : %s\n", uri, msg.Status.Duration())
		}
	}
}
