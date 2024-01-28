package runtime

import "fmt"

type testStruct struct {
	id string
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

func ExampleNillable_FunctionConstraints() {
	s := testGeneric[testStruct](testStruct{})
	fmt.Printf("test: testGeneric() -> %v\n", s)

	// Won't work *struct{} != testStruct{}
	//s = testGeneric[Nillable](testStruct{})
	//fmt.Printf("test: testGeneric() -> %v\n", s)

	// Works - *struct{}, a pointer, can be nil
	s = testGeneric[Nillable](nil)
	fmt.Printf("test: testGeneric() -> %v\n", s)

	//Output:
	//test: testGeneric() -> type testStruct
	//test: testGeneric() -> <nil>

}

type nillableString interface {
	string | Nillable
}

type testStruct2[NS nillableString] struct {
	id NS
}

// testConstraints2 - Get constraints
//type testConstraints2 interface {
//	testStruct2
//}

// if you declare methods on a generic type, you must repeat the type parameter declaration on the receiver, even if the type
// parameters are not used in the method scope — in which case you may use the blank identifier _ to make it obvious:

/*
func (m *Model[T]) Push(item T) {
	m.Data = append(m.Data, item)
}

// not using the type param in this method
func (m *Model[_]) String() string {
	return fmt.Sprint(m.Data)
}
*/

/*
func testGeneric2(t *testStruct2[NS]) string {
	switch any(t).(type) {
	case testStruct:
		return "type testStruct"
	case Nillable:
		return "<nil>"
	default:
		return ""
	}
}

func ExampleNillable_StructConstraints() {
	s := testStruct2[string]{id: "test"}//estGeneric[testStruct](testStruct{})

	s = testStruct2[nillableString]{id: "test"}//estGeneric[testStruct](testStruct{})

	fmt.Printf("test: testGeneric() -> %v\n", s)



	//Output:
	//test: testGeneric() -> type testStruct
	//test: testGeneric() -> <nil>

}


*/
