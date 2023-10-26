package host

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"time"
)

var msgTest = Message{To: "to-uri", From: "from-uri", Content: []any{
	"text content",
	500,
	Credentials(func() (username, password string, err error) { return "", "", nil }),
	time.Second,
	nil,
	//runtime.Handle[runtime.DebugError](),
	errors.New("this is a content error message"),
	func() bool { return false },
	runtime.NewStatusError(0, "location", errors.New("error message")).SetDuration(time.Second * 2),
	//runtime.HandleWithContext[runtime.DebugError](),
	ControllerApply(func(ctx context.Context, statusCode func() int, uri, requestId, method string) (func(), context.Context, bool) {
		return func() {}, ctx, false
	}),
	Resource{"postgres://username:password@database.cloud.timescale.com/database?sslmode=require"},
}}

func ExampleAccessCredentials() {
	fmt.Printf("test: AccessCredentials(nil) -> %v\n", AccessCredentials(nil) != nil)
	fmt.Printf("test: AccessCredentials(msg) -> %v\n", AccessCredentials(&Message{To: "to-uri"}) != nil)
	fmt.Printf("test: AccessCredentials(msg) -> %v\n", AccessCredentials(&msgTest) != nil)

	//Output:
	//test: AccessCredentials(nil) -> false
	//test: AccessCredentials(msg) -> false
	//test: AccessCredentials(msg) -> true
}

func ExampleAccessResource() {
	fmt.Printf("test: AccessResource(nil) -> %v\n", AccessResource(nil))
	fmt.Printf("test: AccessResource(msg) -> %v\n", AccessResource(&Message{To: "to-uri"}))
	fmt.Printf("test: AccessResource(msg) -> %v\n", AccessResource(&msgTest))

	//Output:
	//test: AccessResource(nil) -> {}
	//test: AccessResource(msg) -> {}
	//test: AccessResource(msg) -> {postgres://username:password@database.cloud.timescale.com/database?sslmode=require}

}

func ExampleAccessControllerApply() {
	fmt.Printf("test: AccessControllerApply(nil) -> [valid:%v]\n", AccessControllerApply(nil) != nil)
	fmt.Printf("test: AccessControllerApply(msg) -> [valid:%v]\n", AccessControllerApply(&Message{To: "to-uri"}) != nil)
	fmt.Printf("test: AccessControllerApply(msg) -> [valid:%v]\n", AccessControllerApply(&msgTest) != nil)

	//Output:
	//test: AccessControllerApply(nil) -> [valid:false]
	//test: AccessControllerApply(msg) -> [valid:false]
	//test: AccessControllerApply(msg) -> [valid:true]

}

func ExampleNewStatusCodeFn() {
	var status *runtime.Status

	fn := NewStatusCode(&status)
	status = runtime.NewStatus(runtime.StatusDeadlineExceeded)
	fmt.Printf("test: NewStatusCode(&status) -> [statusCode:%v]\n", fn())

	//Output:
	//test: NewStatusCode(&status) -> [statusCode:4]

}
