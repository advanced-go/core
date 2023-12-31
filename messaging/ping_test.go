package messaging

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"time"
)

func ExamplePing() {
	uri1 := "urn:ping:good"
	uri2 := "urn:ping:bad"
	uri3 := "urn:ping:error"
	uri4 := "urn:ping:delay"

	start = time.Now()
	pingDir := any(NewExchange()).(*exchange)

	c := make(chan Message, 16)
	pingDir.Add(newMailbox(uri1, false, c, nil))
	go pingGood(c)
	status := ping[runtime.Output](pingDir, nil, uri1)
	duration := status.Duration()
	fmt.Printf("test: Ping(good) -> [%v] [duration:%v]\n", status, duration)

	c = make(chan Message, 16)
	pingDir.Add(newMailbox(uri2, false, c, nil))
	go pingBad(c)
	status = ping[runtime.Output](pingDir, nil, uri2)
	fmt.Printf("test: Ping(bad) -> [%v] [duration:%v]\n", status, duration)

	c = make(chan Message, 16)
	pingDir.Add(newMailbox(uri3, false, c, nil))
	go pingError(c, errors.New("ping depends error message"))
	status = ping[runtime.Output](pingDir, nil, uri3)
	fmt.Printf("test: Ping(error) -> [%v] [duration:%v]\n", status, duration)

	c = make(chan Message, 16)
	pingDir.Add(newMailbox(uri4, false, c, nil))
	go pingDelay(c)
	status = ping[runtime.Output](pingDir, nil, uri4)
	fmt.Printf("test: Ping(delay) -> [%v] [duration:%v]\n", status, duration)

	//Output:
	//test: Ping(good) -> [OK] [duration:0s]
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
			SendReply(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
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
			SendReply(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
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
				SendReply(msg, runtime.NewStatusError(0, pingLocation, err).SetDuration(time.Since(start)))
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
			SendReply(msg, runtime.NewStatusOK().SetDuration(time.Since(start)))
		default:
		}
	}
}
