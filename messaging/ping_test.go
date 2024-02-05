package messaging

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var pingStart = time.Now()

func ExamplePing_Good() {
	uri1 := "urn:ping:good"
	pingDir := NewExchange()

	c := make(chan *Message, 16)
	pingDir.Add(NewMailboxWithCtrl(uri1, false, c, nil))
	go pingGood(c)
	status := ping(nil, pingDir, uri1)
	fmt.Printf("test: Ping(good) -> [%v] [duration:%v]\n", status, status.Duration)

	//Output:
	//test: Ping(good) -> [200] [duration:615.1358ms]

}

func ExamplePing_Bad() {
	uri2 := "urn:ping:bad"
	c := make(chan *Message, 16)

	pingDir := NewExchange()
	pingDir.Add(NewMailboxWithCtrl(uri2, false, c, nil))
	go pingBad(c)
	status := ping(nil, pingDir, uri2)
	fmt.Printf("test: Ping(bad) -> [%v] [duration:%v]\n", status, status.Duration)

	//Output:
	//test: Ping(bad) -> [504] [duration:615.1358ms]

}

func ExamplePing_Error() {
	uri3 := "urn:ping:error"
	pingDir := NewExchange()

	c := make(chan *Message, 16)
	pingDir.Add(NewMailboxWithCtrl(uri3, false, c, nil))
	go pingError(c, errors.New("ping response error"))
	status := ping(nil, pingDir, uri3)
	fmt.Printf("test: Ping(error) -> [%v] [error:%v] [duration:%v]\n", status.Code, status, status.Duration)

	//Output:
	//test: Ping(error) -> [418] [error:ping response error] [duration:5.4625108s]

}

func ExamplePing_Delay() {
	uri4 := "urn:ping:delay"
	pingDir := NewExchange()

	c := make(chan *Message, 16)
	pingDir.Add(NewMailboxWithCtrl(uri4, false, c, nil))
	go pingDelay(c)
	status := ping(nil, pingDir, uri4)
	fmt.Printf("test: Ping(delay) -> [%v] [duration:%v]\n", status, status.Duration)

	//Output:
	//test: Ping(delay) -> [200] [duration:6.6710815s]

}

func pingGood(c chan *Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			SendReply(msg, NewStatusDuration(http.StatusOK, time.Since(pingStart)))
		default:
		}
	}
}

func pingBad(c chan *Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			time.Sleep(maxWait + time.Second)
			SendReply(msg, NewStatusDuration(http.StatusOK, time.Since(pingStart)))
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
				SendReply(msg, NewStatusDurationError(http.StatusTeapot, time.Since(pingStart), err))
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
			time.Sleep(time.Second)
			SendReply(msg, NewStatusDuration(http.StatusOK, time.Since(pingStart)))
		default:
		}
	}
}
