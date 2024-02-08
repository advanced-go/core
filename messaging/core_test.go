package messaging

import (
	"fmt"
)

func handler(msg *Message) {
	fmt.Printf(msg.Event())
}

func Example_ReplyTo() {
	msg := NewMessageWithReply("test", "", "startup", handler)
	SendReply(msg, StatusOK())

	msg = NewMessage("test", "", "startup")
	SendReply(msg, StatusOK())

	//Output:
	//startup

}
