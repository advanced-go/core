package runtime

import "fmt"

func Example_NewStatusOK() {
	ok := NewStatusOK()
	ok2 := NewStatusOK()

	fmt.Printf("test: NewStatusOK() -> [ok:%v] [ok2:%v] [equal:%v]\n", ok == status2, ok2 == status2, ok == ok2)

	//Output:
	//test: NewStatusOK() -> [ok:true] [ok2:true] [equal:true]

}
