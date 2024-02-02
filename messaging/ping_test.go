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
	pingDir := any(NewExchange()).(*exchange)

	c := make(chan Message, 16)
	pingDir.Add(NewMailboxWithCtrl(uri1, false, c, nil))
	go pingGood(c)
	status := ping(pingDir, nil, uri1)
	//duration := status.Duration()
	fmt.Printf("test: Ping(good) -> [%v] [duration:%v]\n", status.Error, time.Since(pingStart))

	c = make(chan Message, 16)
	pingDir.Add(NewMailboxWithCtrl(uri2, false, c, nil))
	go pingBad(c)
	status = ping(pingDir, nil, uri2)
	fmt.Printf("test: Ping(bad) -> [%v] [duration:%v]\n", status.Error, time.Since(pingStart))

	c = make(chan Message, 16)
	pingDir.Add(NewMailboxWithCtrl(uri3, false, c, nil))
	go pingError(c, errors.New("ping depends error message"))
	status = ping(pingDir, nil, uri3)
	fmt.Printf("test: Ping(error) -> [%v] [duration:%v]\n", status.Error, time.Since(pingStart))

	c = make(chan Message, 16)
	pingDir.Add(NewMailboxWithCtrl(uri4, false, c, nil))
	go pingDelay(c)
	status = ping(pingDir, nil, uri4)
	fmt.Printf("test: Ping(delay) -> [%v] [duration:%v]\n", status.Error == nil, time.Since(pingStart))

	//Output:
	//test: Ping(good) -> [OK] [duration:0s]
	//{ "code":4, "status":"Deadline Exceeded", "request-id":null, "trace" : [ "https://github.com/advanced-go/core/tree/main/messaging#Ping" ], "errors" : [ "ping response time out: [urn:ping:bad]" ] }
	//test: Ping(bad) -> [Deadline Exceeded [ping response time out: [urn:ping:bad]]] [duration:0s]
	//test: Ping(error) -> [Internal Error [ping depends error message]] [duration:0s]
	//test: Ping(delay) -> [OK] [duration:0s]

}

func pingGood(c chan Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			SendReply(msg, Status{Code: http.StatusOK, Duration: time.Since(pingStart)})
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
			SendReply(msg, Status{Code: http.StatusOK, Duration: time.Since(pingStart)})
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
				SendReply(msg, Status{Code: 0, Error: err, Duration: time.Since(pingStart)})
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
			SendReply(msg, Status{Code: http.StatusOK, Duration: time.Since(pingStart)})
		default:
		}
	}
}
