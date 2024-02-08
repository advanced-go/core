package messaging

import "net/http"

const (
	StartupEvent     = "event:startup"
	ShutdownEvent    = "event:shutdown"
	PingEvent        = "event:ping"
	ReconfigureEvent = "event:reconfigure"

	PauseEvent  = "event:pause"  // disable data channel receive
	ResumeEvent = "event:resume" // enable data channel receive

	XRelatesTo = "x-relates-to"
	XMessageId = "x-message-id"
	XTo        = "x-to"
	XFrom      = "x-from"
	XEvent     = "x-event"
)

// MessageMap - map of messages
type MessageMap map[string]*Message

// MessageHandler - function type to process a Message
type MessageHandler func(msg *Message)

// Message - message payload
type Message struct {
	Header  http.Header
	Status  *Status
	Body    any
	ReplyTo MessageHandler
	Config  map[string]string
}

func (m *Message) To() string {
	return m.Header.Get(XTo)
}
func (m *Message) From() string {
	return m.Header.Get(XFrom)
}

func (m *Message) Event() string {
	return m.Header.Get(XEvent)
}

func (m *Message) RelatesTo() string {
	return m.Header.Get(XRelatesTo)
}

func NewMessage(to, from, event string) *Message {
	m := new(Message)
	m.Header = make(http.Header)
	m.Header.Add(XTo, to)
	m.Header.Add(XFrom, from)
	m.Header.Add(XEvent, event)
	return m
}

func NewMessageWithReply(to, from, event string, replyTo MessageHandler) *Message {
	m := NewMessage(to, from, event)
	m.ReplyTo = replyTo
	return m
}

func NewMessageWithStatus(to, from, event string, status *Status) *Message {
	m := NewMessage(to, from, event)
	m.Status = status
	return m
}

// SendReply - function used by message recipient to reply with a Status
func SendReply(msg *Message, status *Status) {
	if msg == nil || msg.ReplyTo == nil {
		return
	}
	m := NewMessageWithStatus(msg.From(), msg.To(), msg.Event(), status)
	m.Header.Add(XRelatesTo, msg.RelatesTo())
	msg.ReplyTo(m)
}
