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

var messageHandler startup.MessageHandler = func(msg startup.Message) {
	switch msg.Event {
	case startup.StartupEvent:
		//clientStartup(msg)
	case startup.ShutdownEvent:
		//ClientShutdown()
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
