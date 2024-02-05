package messaging

import (
	"net/http"
	"time"
)

type DoneFunc func(msg *Message) bool

func NewReceiverReplyTo(reply chan *Message) MessageHandler {
	return func(msg *Message) {
		if msg != nil {
			reply <- msg
		}
	}
}

// Receiver - receives reply messages and forwards to a function which will return true if the receiving is complete. The interval
// bounds the time spent receiving, and result status is sent on the status channel.
func Receiver(interval time.Duration, reply <-chan *Message, result chan<- *Status, done DoneFunc) {
	tick := time.Tick(interval)
	var status *Status
	start := time.Now().UTC()

	if interval <= 0 || reply == nil || result == nil || done == nil {
		return
	}
	defer func() {
		result <- status
	}()
	for {
		select {
		case <-tick:
			status = NewStatusDuration(http.StatusGatewayTimeout, time.Since(start))
			return
		case msg, open := <-reply:
			if !open {
				status = NewStatusDuration(http.StatusInternalServerError, time.Since(start))
				return
			}
			if done(msg) {
				status = NewStatusDuration(http.StatusOK, time.Since(start))
				return
			}
		default:
		}
	}
}

func DrainAndClose(duration time.Duration, c chan *Message) {
	tick := time.Tick(time.Second * 10)
	for {
		select {
		case <-tick:
			close(c)
			return
		case <-c:
			close(c)
			return
		default:
		}
	}
}
