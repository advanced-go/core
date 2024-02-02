package messaging

import (
	"fmt"
)

func handler(msg Message) {
	fmt.Printf(msg.Event)
}

func Example_ReplyTo() {
	msg := Message{To: "test", Event: "startup", ReplyTo: handler}
	SendReply(msg, StatusOK())

	msg = Message{To: "test", Event: "startup", ReplyTo: nil}
	SendReply(msg, StatusOK())

	//Output:
	//startup

}
