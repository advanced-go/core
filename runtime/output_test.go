package runtime

import "fmt"

type OutputHandlerFn func(s string)

func (OutputHandlerFn) Write(s string) {

}

func ExampleOutputHandler() {
	fmt.Printf("test: Output[NilOutputHandler](s)\n")
	Output[NilOutput]("nil output handler")

	fmt.Printf("test: Output[StdOutputHandler](s)\n")
	Output[StdOutput]("std output handler")

	fmt.Printf("test: Output[TestOutputHandler](s)\n")
	Output[TestOutput]("test output handler")

	fmt.Printf("test: Output[LogOutputHandler](s)\n")
	Output[LogOutput]("log output handler")

	//Output:
	//test: Output[NilOutputHandler](s)
	//test: Output[StdOutputHandler](s)
	//std output handler
	//test: Output[TestOutputHandler](s)
	//test: Write() -> [test output handler]
	//test: Output[LogOutputHandler](s)

}

func ExampleOutputHandler2() {
	h := Output2[NilOutput]()
	h("nil output handler")

	h = Output2[StdOutput]()
	h("std output handler")

	//fmt.Printf("test: Output[NilOutputHandler](s)\n")
	//Output[NilOutputHandler]("nil output handler")

	//Output:
	//std output handler

}

func Output[O OutputHandler](s string) {
	var o O
	o.Write(s)
}

func Output2[O OutputHandler]() func(s string) {
	var o O
	return func(s string) {
		//var o O
		o.Write(s)
	}
}
