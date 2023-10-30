package runtime

import "fmt"

func Example_NilType() {
	b := isNil[NilType](nil)

	fmt.Printf("test: isNil() -> %v", b)

	//Output:
	//test: isNil() -> true
	
}

func isNil[T NilType](t T) bool {
	return t == nil
}
