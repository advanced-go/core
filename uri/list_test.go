package uri

import (
	"fmt"
	"reflect"
)

func Example_ListFromType() {
	fn := listFromType(nil)
	fmt.Printf("test: listFromType(nil) -> [value:<nil>][%v]\n", fn(""))

	n := 345
	fn = listFromType(n)
	fmt.Printf("test: listFromType(nil) -> [value:%v] [%v]\n", n, fn(""))

	fn = listFromType("")
	fmt.Printf("test: listFromType(\"\") -> [value:%v] [%v]\n", "", fn(""))

	k := "key-2"
	v := []string{"const-val-1", "const-val-2"}
	fn = listFromType(v)
	fmt.Printf("test: listFromType(\"\") -> [value:%v] [%v]\n", v, fn(""))
	fmt.Printf("test: listFromType(\"%v\") -> [value:%v] [%v]\n", k, v, fn(k))

	k = "map-key"
	v1 := map[string][]string{k: {"map-value"}}
	fn = listFromType(v1)
	fmt.Printf("test: listFromType(\"\") -> [value:%v] [%v]\n", v1, fn(""))
	fmt.Printf("test: listFromType(\"%v\") -> [value:%v] [%v]\n", k, v1, fn(k))

	k = "fn-key"
	v2 := func(key string) []string {
		switch key {
		case k:
			return []string{"fn-value"}
		}
		return []string{}
	}
	fn = listFromType(v2)
	fmt.Printf("test: listFromType(\"test-key\") -> [value:%v][%v]\n", reflect.TypeOf(v2), fn("test-key"))
	fmt.Printf("test: listFromType(\"%v\") -> [value:%v] [%v]\n", k, reflect.TypeOf(v2), fn(k))

	//Output:
	//test: listFromType(nil) -> [value:<nil>][[error: listFromType() value parameter is nil]]
	//test: listFromType(nil) -> [value:345] [[error: listFromType() value parameter is an invalid type: int]]
	//test: listFromType("") -> [value:] [[ ]]
	//test: listFromType("") -> [value:[const-val-1 const-val-2]] [[const-val-1 const-val-2]]
	//test: listFromType("key-2") -> [value:[const-val-1 const-val-2]] [[const-val-1 const-val-2]]
	//test: listFromType("") -> [value:map[map-key:[map-value]]] [[]]
	//test: listFromType("map-key") -> [value:map[map-key:[map-value]]] [[map-value]]
	//test: listFromType("test-key") -> [value:func(string) []string][[]]
	//test: listFromType("fn-key") -> [value:func(string) []string] [[fn-value]]

}
