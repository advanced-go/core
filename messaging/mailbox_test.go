package messaging

import (
	"fmt"
)

func Example_NewMailbox() {
	m := NewMailbox("github.com/advanced-go/messaging", nil)
	fmt.Printf("test: NewMailbox() -> %v", m)

	//Output:
	//test: NewMailbox() -> github.com/advanced-go/messaging

}

func newMailbox(uri string, public bool, ctrl, data chan Message) *Mailbox {
	m := new(Mailbox)
	m.public = public
	m.uri = uri
	m.ctrl = ctrl
	m.data = data
	return m
}

func newDefaultMailbox(uri string) *Mailbox {
	m := new(Mailbox)
	m.uri = uri
	m.ctrl = make(chan Message, 16)
	return m
}
