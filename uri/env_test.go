package uri

import (
	"fmt"
	"os"
)

func ExampleLookupEnv() {
	name := ""

	s, err := LookupEnv(name)
	fmt.Printf("test: LookupEnv(%v) -> [err:%v][%v]\n", name, err, s)

	s, err = LookupEnv("$")
	fmt.Printf("test: LookupEnv(%v) -> [err:%v][%v]\n", name, err, s)

	os.Setenv("RUNTIME", "DEV")
	s, err = LookupEnv("$RUNTIME")
	fmt.Printf("test: LookupEnv(%v) -> [err:%v][%v]\n", name, err, s)

	//Output:
	//test: LookupEnv() -> [err:invalid argument : LookupEnv() template variable is invalid: ][]
	//test: LookupEnv() -> [err:<nil>][]
	//test: LookupEnv() -> [err:<nil>][DEV]

}
