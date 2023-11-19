package runtime

import (
	"fmt"
)

type testStruct struct {
	id string
}

func Example_genericConstraints() {
	s := testGeneric[testStruct](testStruct{})
	fmt.Printf("test: testGeneric() -> %v\n", s)

	s = testGeneric[Nillable](nil)
	fmt.Printf("test: testGeneric() -> %v\n", s)

	//Output:
	//test: testGeneric() -> type testStruct
	//test: testGeneric() -> <nil>

}

// testConstraints - Get constraints
type testConstraints interface {
	testStruct | Nillable
}

func testGeneric[T testConstraints](t T) string {
	switch any(t).(type) {
	case testStruct:
		return "type testStruct"
	case Nillable:
		return "<nil>"
	default:
		return ""
	}
}
