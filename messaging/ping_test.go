package messaging

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var pingStart = time.Now()

func ExamplePing() {
	uri1 := "urn:ping:good"
	uri2 := "urn:ping:bad"
	uri3 := "urn:ping:error"
	uri4 := "urn:ping:delay"

	//	start := time.Now()
	pingDir := NewExchange()

	c := make(chan Message, 16)
	pingDir.Add(NewMailboxWithCtrl(uri1, false, c, nil))
	go pingGood(c)
	status := ping(pingDir, nil, uri1)
	//duration := status.Duration()
	fmt.Printf("test: Ping(good) -> [%v] [duration:%v]\n", status, time.Since(pingStart))

	c = make(chan Message, 16)
	pingDir.Add(NewMailboxWithCtrl(uri2, false, c, nil))
	go pingBad(c)
	status = ping(pingDir, nil, uri2)
	fmt.Printf("test: Ping(bad) -> [%v] [duration:%v]\n", status, time.Since(pingStart))

	c = make(chan Message, 16)
	pingDir.Add(NewMailboxWithCtrl(uri3, false, c, nil))
	go pingError(c, errors.New("ping depends error message"))
	status = ping(pingDir, nil, uri3)
	fmt.Printf("test: Ping(error) -> [%v] [duration:%v]\n", status, time.Since(pingStart))

	c = make(chan Message, 16)
	pingDir.Add(NewMailboxWithCtrl(uri4, false, c, nil))
	go pingDelay(c)
	status = ping(pingDir, nil, uri4)
	fmt.Printf("test: Ping(delay) -> [%v] [duration:%v]\n", status, time.Since(pingStart))

	//Output:
	//test: Ping(good) -> [200] [duration:615.1358ms]
	//test: Ping(bad) -> [ping response time out: [urn:ping:bad]] [duration:4.2500588s]
	//test: Ping(error) -> [ping response status not available: [urn:ping:error]] [duration:5.4625108s]
	//test: Ping(delay) -> [200] [duration:6.6710815s]

}

func pingGood(c chan Message) {
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

func pingBad(c chan Message) {
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

func pingError(c chan Message, err error) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			if err != nil {
				time.Sleep(time.Second)
				SendReply(msg, NewStatusDurationError(0, time.Since(pingStart), err))
			}
		default:
		}
	}
}

func pingDelay(c chan Message) {
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
