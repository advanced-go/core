package log

import (
	"github.com/go-ai-agent/core/runtime/startup"
)

var (
	c = make(chan startup.Message, 1)
)

func init() {
	startup.Register(PkgUri, c)
	go receive()
}

var AccessLogger startup.HttpAccessLogFn = accessLogFn

var messageHandler startup.MessageHandler = func(msg startup.Message) {
	switch msg.Event {
	case startup.StartupEvent:
		AccessLogger = accessLogFn
	case startup.ShutdownEvent:
	case startup.PingEvent:
	}
}

func receive() {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			go messageHandler(msg)
		default:
		}
	}
}
