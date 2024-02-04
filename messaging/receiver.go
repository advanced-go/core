package messaging

import (
	"time"
)

const (
	ReceiveTimeout = 1
	ReceiveDone    = 0
)

type DoneFunc func(msg *Message) bool

func NewReceiverReplyTo(reply chan *Message) MessageHandler {
	return func(msg *Message) {
		if msg != nil {
			reply <- msg
		}
	}
}

// Receiver - receives reply messages and forward to a function which will return true if the receiving is complete. The interval and
// tries bound the time spent receiving, and an optional status channel can be supplied.
func Receiver(interval time.Duration, tries int, reply <-chan *Message, status chan<- int, done DoneFunc) {
	tick := time.Tick(interval)
	//var msg Message
	reason := ReceiveDone

	if reply == nil || done == nil {
		return
	}
	defer func() {
		if status != nil {
			status <- reason
		}
	}()
	for {
		select {
		case <-tick:
			tries--
			if tries <= 0 {
				reason = ReceiveTimeout
				return
			}
		case msg, open := <-reply:
			if !open {
				return
			}
			if done(msg) {
				return
			}
		default:
		}
	}
}
