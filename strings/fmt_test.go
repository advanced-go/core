package strings

import (
	"fmt"
	"time"
)

func _ExampleFmtTimestamp() {
	t := time.Now().UTC()
	s := FmtTimestamp(t)
	fmt.Printf("test: FmtTimestamp() -> [%v]\n", s)

	t2, err := ParseTimestamp(s)
	fmt.Printf("test: ParseTimestamp() -> [%v] [%v]\n", FmtTimestamp(t2), err)

	//Output:

}

func stringFn() string {
	s := "in stringFn()"
	fmt.Printf("%v\n", s)
	return s
}

type stringerFunc func() string

func (f stringerFunc) String() string {
	return f()
}

func Example_StringerFunc() {
	var str fmt.Stringer

	str = stringerFunc(stringFn)
	s := str.String()
	fmt.Printf("test: fmt.StringerFunc() -> %v\n", s)

	//Output:
	//in stringFn()
	//test: fmt.StringerFunc() -> in stringFn()

}
