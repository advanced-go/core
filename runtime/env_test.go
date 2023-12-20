package runtime

import (
	"fmt"
)

func Example_RuntimeEnv() {
	fmt.Printf("test: IsProdEnvironment() -> %v\n", IsProdEnvironment())
	fmt.Printf("test: IsTestEnvironment() -> %v\n", IsTestEnvironment())
	fmt.Printf("test: IsStageEnvironment() -> %v\n", IsStageEnvironment())
	fmt.Printf("test: IsDebugEnvironment() -> %v\n", IsDebugEnvironment())

	SetProdEnvironment()
	fmt.Printf("test: IsProdEnvironment() -> %v\n", IsProdEnvironment())

	SetTestEnvironment()
	fmt.Printf("test: IsTestEnvironment() -> %v\n", IsTestEnvironment())

	SetStageEnvironment()
	fmt.Printf("test: IsStageEnvironment() -> %v\n", IsStageEnvironment())

	rte = debug
	fmt.Printf("test: IsDebugEnvironment() -> %v\n", IsDebugEnvironment())

	//Output:
	//test: IsProdEnvironment() -> false
	//test: IsTestEnvironment() -> false
	//test: IsStageEnvironment() -> false
	//test: IsDebugEnvironment() -> true
	//test: IsProdEnvironment() -> true
	//test: IsTestEnvironment() -> true
	//test: IsStageEnvironment() -> true
	//test: IsDebugEnvironment() -> true

}
