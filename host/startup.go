package host

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/messaging"
	"net/http"
	"time"
)

const (
	startupLocation = PkgPath + ":Startup"
)

// ContentMap - slice of any content to be included in a message
type ContentMap map[string]map[string]string

// Startup - templated function to start all registered resources.
func Startup(duration time.Duration, content ContentMap) bool {
	return startup(messaging.HostExchange, duration, content)
}

func startup(ex *messaging.Exchange, duration time.Duration, content ContentMap) bool {
	var failures []string
	var count = ex.Count()

	if count == 0 {
		return true
	}
	cache := messaging.NewMessageCache()
	toSend := createToSend(ex, content, messaging.NewMessageCacheHandler(cache))
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
			return true
		}
		break
	}
	shutdownHost(messaging.NewMessage("", "", messaging.ShutdownEvent))
	if len(failures) > 0 {
		handleErrors(failures, cache)
		return false
	}
	fmt.Printf("error: startup failure [%v]\n", errors.New(fmt.Sprintf("response counts < directory entries [%v] [%v]", cache.Count(), ex.Count())))
	return false
}

func createToSend(ex *messaging.Exchange, cm ContentMap, fn messaging.MessageHandler) messaging.MessageMap {
	m := make(messaging.MessageMap)
	for _, k := range ex.List() {
		msg := messaging.NewMessage(k, startupLocation, messaging.StartupEvent)
		msg.ReplyTo = fn
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

func sendMessages(ex *messaging.Exchange, msgs messaging.MessageMap) {
	for k := range msgs {
		ex.SendCtrl(msgs[k])
	}
}

func handleErrors(failures []string, cache *messaging.MessageCache) {
	for _, uri := range failures {
		msg, ok := cache.Get(uri)
		if !ok {
			continue
		}
		if msg.Status.Error != nil {
			//loc := ""
			//if msg.Status.Location() != nil && len(msg.Status.Location()) > 0 {
			//	loc = msg.Status.Location()[0]
			//}
			fmt.Printf("error: startup failure [%v]\n", msg.Status.Error)
		}
	}
}

func handleStatus(cache *messaging.MessageCache) {
	for _, uri := range cache.Uri() {
		msg, ok := cache.Get(uri)
		if !ok {
			continue
		}
		//if msg.Status != nil {
		fmt.Printf("startup successful: [%v] : %s\n", uri, msg.Status.Duration)
		//}
	}
}
