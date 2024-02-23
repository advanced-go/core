package messaging

import (
	"errors"
)

const (
	newAgentLocation = PkgPath + ":NewDefaultAgent"
)

// Agent - AI Agent
type Agent struct {
	m   *Mailbox
	run func(m *Mailbox)
}

// NewDefaultAgent - create an agent with only a control channel, registered with the HostDirectory,
// and using the default run function.
func NewDefaultAgent(uri string, ctrlHandler Handler, public bool) (*Agent, error) {
	return newDefaultAgent(uri, ctrlHandler, HostExchange, public)
}

func newDefaultAgent(uri string, ctrlHandler Handler, e *Exchange, public bool) (*Agent, error) {
	if len(uri) == 0 {
		return nil, errors.New("error: agent URI is empty")
	}
	if ctrlHandler == nil {
		return nil, errors.New("error: agent controller message handler is nil")
	}
	a := new(Agent)
	a.m = NewMailbox(uri, nil)
	a.m.public = public
	a.run = func(m *Mailbox) {
		DefaultRun(m, ctrlHandler)
	}
	return a, a.Register(e)
}

// Run - run the agent
func (a *Agent) Run() {
	go a.run(a.m)
}

// Shutdown - shutdown the agent
func (a *Agent) Shutdown() {
	a.m.Send(NewControlMessage("", "", ShutdownEvent))
}

// Send - send a message
func (a *Agent) Send(msg *Message) {
	a.m.Send(msg)
}

// SendData - send a message to the data channel
/*
func (a *Agent) SendData(msg *Message) {
	a.m.SendData(msg)
}


*/

// Register - register an agent with a directory
func (a *Agent) Register(e *Exchange) error {
	return e.Add(a.m)
}

// DefaultRun - a simple run function that only handles control messages, and dispatches via a message handler
func DefaultRun(m *Mailbox, ctrlHandler Handler) {
	for {
		select {
		case msg, open := <-m.ctrl:
			if !open {
				return
			}
			switch msg.Event() {
			case ShutdownEvent:
				ctrlHandler(NewMessageWithStatus(ChannelControl, "", "", msg.Event(), StatusOK()))
				m.Close()
				return
			default:
				ctrlHandler(msg)
			}
		default:
		}
	}
}
