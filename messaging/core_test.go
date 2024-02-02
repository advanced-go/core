package messaging

import (
	"fmt"
	"net/http"
)

func handler(msg Message) {
	fmt.Printf(msg.Event)
}

func Example_ReplyTo() {
	msg := Message{To: "test", Event: "startup", ReplyTo: handler}
	SendReply(msg, Status{Code: http.StatusOK})

	msg = Message{To: "test", Event: "startup", ReplyTo: nil}
	SendReply(msg, Status{Code: http.StatusOK})

	//Output:
	//startup

}
