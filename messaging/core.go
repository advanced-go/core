package messaging

import (
	"net/http"
	"time"
)

const (
	StartupEvent     = "event:startup"
	ShutdownEvent    = "event:shutdown"
	PingEvent        = "event:ping"
	ReconfigureEvent = "event:reconfigure"

	PauseEvent  = "event:pause"  // disable data channel receive
	ResumeEvent = "event:resume" // enable data channel receive
)

// Status - message status
type Status struct {
	Error    error
	Code     int
	Location string
	Duration time.Duration
}

func NewStatus(code int) Status {
	return Status{Code: code}
}

func StatusOK() Status {
	return Status{Code: http.StatusOK}
}

// MessageMap - map of messages
type MessageMap map[string]Message

// MessageHandler - function type to process a Message
type MessageHandler func(msg Message)

// Message - message payload
type Message struct {
	To        string
	From      string
	Event     string
	RelatesTo string
	Status    Status
	Content   []any
	ReplyTo   MessageHandler
	Config    map[string]string
}

// SendReply - function used by message recipient to reply with a Status
func SendReply(msg Message, status Status) {
	if msg.ReplyTo == nil {
		return
	}
	msg.ReplyTo(Message{
		To:        msg.From,
		From:      msg.To,
		RelatesTo: msg.RelatesTo,
		Event:     msg.Event,
		Status:    status,
		Content:   nil,
		ReplyTo:   nil,
	})
}
