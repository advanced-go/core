package messaging

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	maxWait = timeout + time.Millisecond*100
)

var pingStart = time.Now()

func ExamplePing_Good() {
	uri1 := "urn:ping:good"
	pingDir := NewExchange()

	c := make(chan *Message, 16)
	pingDir.Add(NewMailboxWithCtrl(uri1, false, c, nil))
	go pingGood(c)
	status := ping(nil, pingDir, uri1)
	fmt.Printf("test: Ping(good) -> [%v] [timeout:%v] [duration:%v]\n", status, timeout, status.Duration)

	//Output:
	//test: Ping(good) -> [OK] [timeout:3s] [duration:0s]

}

func ExamplePing_Timeout() {
	uri2 := "urn:ping:timeout"
	c := make(chan *Message, 16)

	pingDir := NewExchange()
	pingDir.Add(NewMailboxWithCtrl(uri2, false, c, nil))
	go pingTimeout(c)
	status := ping(nil, pingDir, uri2)
	fmt.Printf("test: Ping(timeout) -> [%v] [timeout:%v] [duration:%v]\n", status, timeout, status.Duration)

	//Output:
	//test: Ping(timeout) -> [Timeout] [timeout:3s]

}

func ExamplePing_Error() {
	uri3 := "urn:ping:error"
	pingDir := NewExchange()

	c := make(chan *Message, 16)
	pingDir.Add(NewMailboxWithCtrl(uri3, false, c, nil))
	go pingError(c, errors.New("ping response error"))
	status := ping(nil, pingDir, uri3)
	fmt.Printf("test: Ping(error) -> [%v] [error:%v] [timeout:%v] [duration:%v]\n", status.Code, status.Error(), timeout, status.Duration)

	//Output:
	//recovered in messaging.NewReceiverReplyTo() : send on closed channel
	//test: Ping(error) -> [418] [error:ping response error] [timeout:3s] [duration:1.0151556s]

}

func ExamplePing_Delay() {
	uri4 := "urn:ping:delay"
	pingDir := NewExchange()

	c := make(chan *Message, 16)
	pingDir.Add(NewMailboxWithCtrl(uri4, false, c, nil))
	go pingDelay(c)
	status := ping(nil, pingDir, uri4)
	fmt.Printf("test: Ping(delay) -> [%v] [timeout:%v] [duration:%v]\n", status, timeout, status.Duration)

	//Output:
	//test: Ping(delay) -> [OK] [timeout:3s]

}

func pingGood(c chan *Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			SendReply(msg, StatusOK())
		default:
		}
	}
}

func pingTimeout(c chan *Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			time.Sleep(maxWait)
			SendReply(msg, StatusOK())
		default:
		}
	}
}

func pingError(c chan *Message, err error) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			if err != nil {
				time.Sleep(time.Second)
				SendReply(msg, NewStatusError(http.StatusTeapot, errors.New("ping response error"), nil))
			}
		default:
		}
	}
}

func pingDelay(c chan *Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			time.Sleep(timeout / 2)
			SendReply(msg, StatusOK())
		default:
		}
	}
}
