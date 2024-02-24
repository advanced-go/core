package messaging

// Mailbox - mailbox struct
type Mailbox struct {
	public     bool
	uri        string
	ctrl       chan *Message
	data       chan *Message
	unregister func()
}

// NewMailboxWithCtrl - create a mailbox
func NewMailboxWithCtrl(uri string, public bool, ctrl, data chan *Message) *Mailbox {
	m := new(Mailbox)
	m.public = public
	m.uri = uri
	m.ctrl = ctrl
	m.data = data
	return m
}

// NewMailbox - create a mailbox
func NewMailbox(uri string, data chan *Message) *Mailbox {
	m := new(Mailbox)
	m.uri = uri
	m.ctrl = make(chan *Message, 16)
	m.data = data
	return m
}

// NewPublicMailbox - create a public mailbox
func NewPublicMailbox(uri string, data chan *Message) *Mailbox {
	m := NewMailbox(uri, data)
	m.public = true
	return m
}

// Uri - return the mailbox uri
func (m *Mailbox) Uri() string {
	return m.uri
}

// String - return the mailbox uri
func (m *Mailbox) String() string {
	return m.uri
}

// Send - send a message
func (m *Mailbox) Send(msg *Message) {
	if msg == nil {
		return
	}
	if msg.Channel == ChannelControl && m.ctrl != nil {
		m.ctrl <- msg
	} else {
		if msg.Channel == ChannelData && m.data != nil {
			m.data <- msg
		}
	}

}

// SendData - send a message to the data channel
/*
func (m *Mailbox) SendData(msg *Message) {
	if m.data != nil {
		m.data <- msg
	}
}


*/

// Close - close the mailbox channels and unregister the mailbox with a Directory
func (m *Mailbox) Close() {
	if m.unregister != nil {
		m.unregister()
	}
	if m.data != nil {
		close(m.data)
	}
	if m.ctrl != nil {
		close(m.ctrl)
	}
}
