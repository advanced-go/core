package runtime

import "fmt"

func Example_StatusOK() {
	ok := StatusOK()
	ok2 := StatusOK()

	fmt.Printf("test: StatusOK() -> [ok:%v] [ok2:%v] [equal:%v]\n", ok == statusOK, ok2 == statusOK, ok == ok2)

	ok = nil
	ok3 := StatusOK()
	ok2 = StatusOK()
	fmt.Printf("test: StatusOK() -> [ok2:%v] [ok3:%v] [equal:%v]\n", ok2 == statusOK, ok3 == statusOK, ok2 == ok3)

	//Output:
	//test: StatusOK() -> [ok:true] [ok2:true] [equal:true]
	//test: StatusOK() -> [ok2:true] [ok3:true] [equal:true]

}
