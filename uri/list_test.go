package uri

import (
	"fmt"
	"reflect"
)

func Example_ListFromType() {
	fn := ListFromType(nil)
	fmt.Printf("test: ListFromType(nil) -> [value:<nil>][%v]\n", fn(""))

	n := 345
	fn = ListFromType(n)
	fmt.Printf("test: ListFromType(nil) -> [value:%v] [%v]\n", n, fn(""))

	fn = ListFromType("")
	fmt.Printf("test: ListFromType(\"\") -> [value:%v] [%v]\n", "", fn(""))

	k := "key-2"
	v := []string{"const-val-1", "const-val-2"}
	fn = ListFromType(v)
	fmt.Printf("test: ListFromType(\"\") -> [value:%v] [%v]\n", v, fn(""))
	fmt.Printf("test: ListFromType(\"%v\") -> [value:%v] [%v]\n", k, v, fn(k))

	k = "map-key"
	v1 := map[string][]string{k: {"map-value"}}
	fn = ListFromType(v1)
	fmt.Printf("test: ListFromType(\"\") -> [value:%v] [%v]\n", v1, fn(""))
	fmt.Printf("test: ListFromType(\"%v\") -> [value:%v] [%v]\n", k, v1, fn(k))

	k = "fn-key"
	v2 := func(key string) []string {
		switch key {
		case k:
			return []string{"fn-value"}
		}
		return []string{}
	}
	fn = ListFromType(v2)
	fmt.Printf("test: ListFromType(\"test-key\") -> [value:%v][%v]\n", reflect.TypeOf(v2), fn("test-key"))
	fmt.Printf("test: ListFromType(\"%v\") -> [value:%v] [%v]\n", k, reflect.TypeOf(v2), fn(k))

	//Output:
	//test: ListFromType(nil) -> [value:<nil>][[error: listFromType() value parameter is nil]]
	//test: ListFromType(nil) -> [value:345] [[error: listFromType() value parameter is an invalid type: int]]
	//test: ListFromType("") -> [value:] [[ ]]
	//test: ListFromType("") -> [value:[const-val-1 const-val-2]] [[const-val-1 const-val-2]]
	//test: ListFromType("key-2") -> [value:[const-val-1 const-val-2]] [[const-val-1 const-val-2]]
	//test: ListFromType("") -> [value:map[map-key:[map-value]]] [[]]
	//test: ListFromType("map-key") -> [value:map[map-key:[map-value]]] [[map-value]]
	//test: ListFromType("test-key") -> [value:func(string) []string][[]]
	//test: ListFromType("fn-key") -> [value:func(string) []string] [[fn-value]]

}
