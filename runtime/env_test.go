package runtime

import (
	"fmt"
	"os"
)

func ExampleSetEnv() {
	fmt.Printf("test: Environment() -> %v\n", EnvStr())

	os.Setenv(EnvKey, "sTaGE")
	setEnv()
	fmt.Printf("test: Environment() -> %v\n", EnvStr())

	os.Setenv(EnvKey, "tesT")
	setEnv()
	fmt.Printf("test: Environment() -> %v\n", EnvStr())

	os.Setenv(EnvKey, "prod")
	setEnv()
	fmt.Printf("test: Environment() -> %v\n", EnvStr())

	//Output:
	//test: Environment() -> debug
	//test: Environment() -> stage
	//test: Environment() -> test
	//test: Environment() -> prod

}
func ExampleRuntimeEnv() {
	os.Setenv(EnvKey, "")
	setEnv()

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
