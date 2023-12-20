package runtime

import (
	"fmt"
	"reflect"
)

func Example_stringFromType() {
	fn := stringFromType(nil)
	fmt.Printf("test: stringFromType(nil) -> [value:<nil>][%v]\n", fn(""))

	n := 345
	fn = stringFromType(n)
	fmt.Printf("test: stringFromType(nil) -> [value:%v] [%v]\n", n, fn(""))

	fn = stringFromType("")
	fmt.Printf("test: stringFromType(\"\") -> [value:%v] [%v]\n", "", fn(""))

	v := "string constant"
	k := "random"
	fn = stringFromType(v)
	fmt.Printf("test: stringFromType(\"\") -> [value:%v] [%v]\n", v, fn(""))
	fmt.Printf("test: stringFromType(\"%v\") -> [value:%v] [%v]\n", k, v, fn(""))

	k = "key-1"
	v1 := map[string]string{k: "map-value"}
	fn = stringFromType(v1)
	fmt.Printf("test: stringFromType(\"\") -> [value:%v] [%v]\n", v1, fn(""))
	fmt.Printf("test: stringFromType(\"%v\") -> [value:%v] [%v]\n", k, v1, fn(k))

	k = "key-2"
	v2 := func(key string) string {
		switch key {
		case k:
			return "value"
		}
		return ""
	}
	fn = stringFromType(v2)
	fmt.Printf("test: stringFromType(\"test-key\") -> [value%v] [%v]\n", reflect.TypeOf(v2), fn("test-key"))
	fmt.Printf("test: stringFromType(\"%v\") -> [value:%v] [%v]\n", k, reflect.TypeOf(v2), fn(k))

	//Output:
	//test: stringFromType(nil) -> [value:<nil>][error: stringFromType() value parameter is nil]
	//test: stringFromType(nil) -> [value:345] [error: stringFromType() value parameter is an invalid type: int]
	//test: stringFromType("") -> [value:] []
	//test: stringFromType("") -> [value:string constant] [string constant]
	//test: stringFromType("random") -> [value:string constant] [string constant]
	//test: stringFromType("") -> [value:map[key-1:map-value]] []
	//test: stringFromType("key-1") -> [value:map[key-1:map-value]] [map-value]
	//test: stringFromType("test-key") -> [valuefunc(string) string] []
	//test: stringFromType("key-2") -> [value:func(string) string] [value]

}

func Example_LookupFromType_String() {
	fn := LookupFromType[func(string) string](nil)
	fmt.Printf("test: LookupFromType(nil) -> [%v]\n", fn(""))

	v := "string constant"
	k := "random"
	fn = LookupFromType[func(string) string](v)
	fmt.Printf("test: LookupFromType(\"\") -> [value:%v] [%v]\n", v, fn(""))
	fmt.Printf("test: LookupFromType(%v) -> [value:%v] [%v]\n", k, v, fn(""))

	k = "key"
	v1 := map[string]string{k: "map-value"}
	fn = LookupFromType[func(string) string](v1)
	fmt.Printf("test: LookupFromType(\"\") -> [value:%v] [%v]\n", v1, fn(""))
	fmt.Printf("test: LookupFromType(\"%v\") -> [value:%v] [%v]\n", k, v1, fn(k))

	k = "key-1"
	v2 := func(key string) string {
		switch key {
		case k:
			return "value"
		}
		return ""
	}
	fn = LookupFromType[func(string) string](v2)
	fmt.Printf("test: LookupFromType(\"test-key\") -> [value:%v] [%v]\n", reflect.TypeOf(v2), fn("test-key"))
	fmt.Printf("test: LookupFromType(\"%v\") -> [value:%v] [%v]\n", k, reflect.TypeOf(v2), fn(k))

	//Output:
	//test: LookupFromType(nil) -> [error: stringFromType() value parameter is nil]
	//test: LookupFromType("") -> [value:string constant] [string constant]
	//test: LookupFromType(random) -> [value:string constant] [string constant]
	//test: LookupFromType("") -> [value:map[key:map-value]] []
	//test: LookupFromType("key") -> [value:map[key:map-value]] [map-value]
	//test: LookupFromType("test-key") -> [value:func(string) string] []
	//test: LookupFromType("key-1") -> [value:func(string) string] [value]

}

func Example_listFromType() {
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

func Example_LookupFromType_List() {
	fn := LookupFromType[func(string) []string](nil)
	//fmt.Printf("test: LookupFromType(\"\") -> [%v]\n", fn)

	v := "string constant"
	k := "random"
	fn = LookupFromType[func(string) []string](v)
	fmt.Printf("test: LookupFromType(\"\") -> [%v]\n", fn(""))
	fmt.Printf("test: LookupFromType(%v) -> [%v]\n", k, fn(""))

	/*
		k = "key"
		v1 := map[string]string{k: "map-value"}
		fn = LookupFromType[func(string) string](v1)
		fmt.Printf("test: LookupFromType(\"\") -> [%v]\n", fn(""))
		fmt.Printf("test: LookupFromType(%v) -> [%v]\n", v1, fn(k))

		k = "key"
		v2 := func(key string) string {
			switch key {
			case "key":
				return "value"
			}
			return ""
		}
		fn = LookupFromType[func(string) string](v2)
		fmt.Printf("test: LookupFromType(\"test-key\") -> [%v]\n", fn("test-key"))
		fmt.Printf("test: LookupFromType() -> [%v]\n", fn(k))


	*/

	//Output:
	//test: LookupFromType("") -> [[string constant ]]
	//test: LookupFromType(random) -> [[string constant ]]

}
