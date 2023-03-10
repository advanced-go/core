package runtime

import "fmt"

type testStruct struct {
	vers  string
	count int
}

func ExampleIsNil() {
	var i any
	var p *int

	fmt.Printf("test: IsNil(nil) -> %v\n", IsNil(nil))
	fmt.Printf("test: IsNil(i) -> %v\n", IsNil(i))
	fmt.Printf("test: IsNil(pi) -> %v\n", IsNil(p))

	//Output:
	//test: IsNil(nil) -> true
	//test: IsNil(i) -> true
	//test: IsNil(pi) -> true

}

/*
func ExampleIsPointer() {
	var i any
	var s string
	var data = testStruct{}
	var count int
	var bytes []byte

	fmt.Printf("any : %v\n", IsPointer(i))
	fmt.Printf("int : %v\n", IsPointer(count))
	fmt.Printf("int * : %v\n", IsPointer(&count))
	fmt.Printf("string : %v\n", IsPointer(s))
	fmt.Printf("string * : %v\n", IsPointer(&s))
	fmt.Printf("struct : %v\n", IsPointer(data))
	fmt.Printf("struct * : %v\n", IsPointer(&data))
	fmt.Printf("[]byte : %v\n", IsPointer(bytes))
	//fmt.Printf("Struct * : %v\n", IsPointer(&data))

	//Output:
	// any : false
	// int : false
	// int * : true
	// string : false
	// string * : true
	// struct : false
	// struct * : true
	// []byte : false

}


*/
