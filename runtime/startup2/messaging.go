package startup2

const (
	StartupEvent  = "event:startup"
	ShutdownEvent = "event:shutdown"
	PingEvent     = "event:ping"
)

// Message - message access data
type Message struct {
	To      string
	From    string
	Event   string
	Content []any
}
