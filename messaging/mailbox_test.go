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

func newDefaultMailbox(uri string) *Mailbox {
	m := new(Mailbox)
	m.uri = uri
	m.ctrl = make(chan *Message, 16)
	return m
}
