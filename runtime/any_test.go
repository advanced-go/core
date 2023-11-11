package runtime

import (
	"fmt"
	"net/http"
)

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

func Example_TypeName() {
	s := TypeName(nil)
	fmt.Printf("test: TypeName(nil) -> %v\n", s)

	s = TypeName("test data")
	fmt.Printf("test: TypeName(string) -> %v\n", s)

	n := 500
	s = TypeName(n)
	fmt.Printf("test: TypeName(int) -> %v\n", s)

	req, _ := http.NewRequest("patch", "https://www.google.com/search", nil)
	s = TypeName(req)
	fmt.Printf("test: TypeName(http.Request) -> %v\n", s)

	//Output:
	//test: TypeName(nil) -> <nil>
	//test: TypeName(string) -> string
	//test: TypeName(int) -> int
	//test: TypeName(http.Request) -> ptr

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

func Example_Nillable() {
	b := isNillable[Nillable](nil)

	fmt.Printf("test: isNil() -> %v", b)

	//Output:
	//test: isNil() -> true

}

func isNillable[T Nillable](t T) bool {
	return t == nil
}
