package messaging

const (
	StartupEvent     = "event:startup"
	ShutdownEvent    = "event:shutdown"
	PingEvent        = "event:ping"
	ReconfigureEvent = "event:reconfigure"

	PauseEvent  = "event:pause"  // disable data channel receive
	ResumeEvent = "event:resume" // enable data channel receive
)

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
	Status    *Status
	Content   []any
	ReplyTo   MessageHandler
	Config    map[string]string
}

// SendReply - function used by message recipient to reply with a Status
func SendReply(msg Message, status *Status) {
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
