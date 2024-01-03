package runtimetest

import "fmt"

func Example_StringerFunc() {
	var str fmt.Stringer

	str = stringerFunc(stringer)
	s := str.String()
	fmt.Printf("test: stringerFunc() -> %v\n", s)

	//Output:
	//in stringer()
	//test: stringerFunc() -> in stringer()

}
