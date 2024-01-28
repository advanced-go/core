package host

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/messaging"
	"github.com/advanced-go/core/runtime"
	"time"
)

var pingStart time.Time

func ExamplePing() {
	uri1 := "urn:ping:good"
	uri2 := "urn:ping:bad"
	uri3 := "urn:ping:error"
	uri4 := "urn:ping:delay"

	pingStart = time.Now()
	pingDir := messaging.NewExchange() //any(messaging.NewExchange()).(*exchange)

	c := make(chan messaging.Message, 16)
	pingDir.Add(messaging.NewMailboxWithCtrl(uri1, false, c, nil))
	go pingGood(c)
	status := ping[runtime.Output](pingDir, nil, uri1)
	duration := status.Duration()
	fmt.Printf("test: Ping(good) -> [%v] [duration:%v]\n", status, duration)

	c = make(chan messaging.Message, 16)
	pingDir.Add(messaging.NewMailboxWithCtrl(uri2, false, c, nil))
	go pingBad(c)
	status = ping[runtime.Output](pingDir, nil, uri2)
	fmt.Printf("test: Ping(bad) -> [%v] [duration:%v]\n", status, duration)

	c = make(chan messaging.Message, 16)
	pingDir.Add(messaging.NewMailboxWithCtrl(uri3, false, c, nil))
	go pingError(c, errors.New("ping depends error message"))
	status = ping[runtime.Output](pingDir, nil, uri3)
	fmt.Printf("test: Ping(error) -> [%v] [duration:%v]\n", status, duration)

	c = make(chan messaging.Message, 16)
	pingDir.Add(messaging.NewMailboxWithCtrl(uri4, false, c, nil))
	go pingDelay(c)
	status = ping[runtime.Output](pingDir, nil, uri4)
	fmt.Printf("test: Ping(delay) -> [%v] [duration:%v]\n", status, duration)

	//Output:
	//test: Ping(good) -> [OK] [duration:0s]
	//{ "code":4, "status":"Deadline Exceeded", "request-id":null, "trace" : [ "https://github.com/advanced-go/core/tree/main/host#Ping" ], "errors" : [ "ping response time out: [urn:ping:bad]" ] }
	//test: Ping(bad) -> [Deadline Exceeded [ping response time out: [urn:ping:bad]]] [duration:0s]
	//test: Ping(error) -> [Internal Error [ping depends error message]] [duration:0s]
	//test: Ping(delay) -> [OK] [duration:0s]

}

func pingGood(c chan messaging.Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			messaging.SendReply(msg, runtime.NewStatusOK().SetDuration(time.Since(pingStart)))
		default:
		}
	}
}

func pingBad(c chan messaging.Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			time.Sleep(maxWait + time.Second)
			messaging.SendReply(msg, runtime.NewStatusOK().SetDuration(time.Since(pingStart)))
		default:
		}
	}
}

func pingError(c chan messaging.Message, err error) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			if err != nil {
				time.Sleep(time.Second)
				messaging.SendReply(msg, runtime.NewStatusError(0, pingLocation, err).SetDuration(time.Since(pingStart)))
			}
		default:
		}
	}
}

func pingDelay(c chan messaging.Message) {
	for {
		select {
		case msg, open := <-c:
			if !open {
				return
			}
			time.Sleep(time.Second)
			messaging.SendReply(msg, runtime.NewStatusOK().SetDuration(time.Since(pingStart)))
		default:
		}
	}
}
