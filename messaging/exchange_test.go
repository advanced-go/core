package messaging

import (
	"fmt"
	"time"
)

func Example_Add() {
	uri1 := "urn:test:one"

	testDir := NewExchange()
	m1 := newDefaultMailbox(uri1)

	fmt.Printf("test: Count() -> : %v\n", testDir.Count())
	m0 := testDir.get(uri1)
	fmt.Printf("test: get(%v) -> : [mbox:%v]\n", uri1, m0)

	err := testDir.Add(m1)
	fmt.Printf("test: Add(%v) -> : [err:%v]\n", uri1, err)

	fmt.Printf("test: Count() -> : %v\n", testDir.Count())
	m0 = testDir.get(uri1)
	fmt.Printf("test: get(%v) -> : [mbox:%v]\n", uri1, m0)

	uri2 := "urn:test:two"

	m2 := newDefaultMailbox(uri2)
	err = testDir.Add(m2)
	fmt.Printf("test: Add(%v) -> : [err:%v]\n", uri2, err)
	fmt.Printf("test: Count() -> : %v\n", testDir.Count())
	m0 = testDir.get(uri2)
	fmt.Printf("test: get(%v) -> : [mbox:%v]\n", uri2, m0)

	fmt.Printf("test: List() -> : %v\n", testDir.List())

	//Output:
	//test: Count() -> : 0
	//test: get(urn:test:one) -> : [mbox:<nil>] [status:Not Found [invalid URI: exchange mailbox not found [urn:test:one]]]
	//test: Add(urn:test:one) -> : [status:OK]
	//test: Count() -> : 1
	//test: get(urn:test:one) -> : [mbox:urn:test:one] [status:OK]
	//test: Add(urn:test:two) -> : [status:OK]
	//test: Count() -> : 2
	//test: get(urn:test:two) -> : [mbox:urn:test:two] [status:OK]
	//test: List() -> : [urn:test:one urn:test:two]

}

func Example_SendError() {
	uri := "urn:test"
	testDir := NewExchange()

	fmt.Printf("test: SendCtrl(%v) -> : %v\n", uri, testDir.SendCtrl(NewMessage(uri, "", "")))

	m := NewMailboxWithCtrl(uri, false, nil, nil)
	status := testDir.Add(m)
	fmt.Printf("test: Add(%v) -> : [status:%v]\n", uri, status)

	//Output:
	//test: SendCtrl(urn:test) -> : error: exchange.SendCtrl() failed as the message To is empty or invalid [urn:test]
	//test: Add(urn:test) -> : [status:error: exchange.Add() mailbox command channel is nil]

}

func Example_Send() {
	uri1 := "urn:test-1"
	uri2 := "urn:test-2"
	uri3 := "urn:test-3"
	c := make(chan *Message, 16)
	testDir := NewExchange()

	testDir.Add(NewMailboxWithCtrl(uri1, false, c, nil))
	testDir.Add(NewMailboxWithCtrl(uri2, false, c, nil))
	testDir.Add(NewMailboxWithCtrl(uri3, false, c, nil))

	testDir.SendCtrl(NewMessage(uri1, PkgPath, StartupEvent))
	testDir.SendCtrl(NewMessage(uri2, PkgPath, StartupEvent))
	testDir.SendCtrl(NewMessage(uri3, PkgPath, StartupEvent))

	time.Sleep(time.Second * 1)
	resp1 := <-c
	resp2 := <-c
	resp3 := <-c
	fmt.Printf("test: <- c -> : [%v] [%v] [%v]\n", resp1.To(), resp2.To(), resp3.To())
	close(c)

	//Output:
	//test: <- c -> : [urn:test-1] [urn:test-2] [urn:test-3]

}

func Example_ListCount() {
	testDir := NewExchange()

	testDir.Add(newDefaultMailbox("test:uri1"))
	testDir.Add(newDefaultMailbox("test:uri2"))

	fmt.Printf("test: Count() -> : %v\n", testDir.Count())

	fmt.Printf("test: List() -> : %v\n", testDir.List())

	//Output:
	//test: Count() -> : 2
	//test: List() -> : [test:uri1 test:uri2]

}

func Example_Remove() {
	uri := "urn:test/one"

	m := newDefaultMailbox(uri)
	testDir := NewExchange()

	status := testDir.Add(m)
	fmt.Printf("test: Add(%v) -> : [%v]\n", uri, status)

	status = testDir.SendCtrl(NewMessage(uri, "", PingEvent))
	fmt.Printf("test: SendCtrl(%v) -> : [%v]\n", uri, status)

	m.Close()

	status = testDir.SendCtrl(NewMessage(uri, "", PingEvent))
	fmt.Printf("test: SendCtrl(%v) -> : [%v]\n", uri, status)

	//Output:
	//test: Add(urn:test/one) -> : [<nil>]
	//test: SendCtrl(urn:test/one) -> : [<nil>]
	//test: SendCtrl(urn:test/one) -> : [error: exchange.SendCtrl() failed as the message To is empty or invalid [urn:test/one]]

}
