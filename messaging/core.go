package messaging

import (
	"errors"
	"net/http"
)

const (
	StartupEvent     = "event:startup"
	ShutdownEvent    = "event:shutdown"
	PingEvent        = "event:ping"
	ReconfigureEvent = "event:reconfigure"

	PauseEvent  = "event:pause"  // disable data channel receive
	ResumeEvent = "event:resume" // enable data channel receive

	ContentType       = "Content-Type"
	XRelatesTo        = "x-relates-to"
	XMessageId        = "x-message-id"
	XTo               = "x-to"
	XFrom             = "x-from"
	XEvent            = "x-event"
	ContentTypeStatus = "application/status"
	ContentTypeConfig = "application/config"
)

// MessageMap - map of messages
type MessageMap map[string]*Message

// MessageHandler - function type to process a Message
type MessageHandler func(msg *Message)

// Message - message payload
type Message struct {
	Header  http.Header
	body    any
	ReplyTo MessageHandler
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
	m.SetContent(ContentTypeStatus, status)
	m.body = status
	return m
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

func (m *Message) Status() *Status {
	ct := m.Header.Get(ContentType)
	if ct != ContentTypeStatus || m.body == nil {
		return nil
	}
	if s, ok := m.body.(*Status); ok {
		return s
	}
	return nil //StatusOK()
}

func (m *Message) Config() map[string]string {
	ct := m.Header.Get(ContentType)
	if ct != ContentTypeConfig || m.body == nil {
		return nil
	}
	if m, ok := m.body.(map[string]string); ok {
		return m
	}
	return nil
}

func (m *Message) Content() (string, any, bool) {
	if m.body == nil {
		return "", nil, false
	}
	ct := m.Header.Get(ContentType)
	if len(ct) == 0 {
		return "", nil, false
	}
	return ct, m.body, true
}

func (m *Message) SetContent(contentType string, content any) error {
	if len(contentType) == 0 {
		return errors.New("error: content type is empty")
	}
	if content == nil {
		return errors.New("error: content is nil")
	}
	m.body = content
	m.Header.Add(ContentType, contentType)
	return nil
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
