package uri

import (
	"fmt"
	"reflect"
)

const (
	passThroughKey = "pass through"
	nobodyHomeKey  = "nobody home"
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

func Example_UriLookup() {
	l := NewLookup()
	k := ""

	v, ok := l.Value(k)
	fmt.Printf("test: Lookup(\"%v\") -> [value:%v] [ok:%v]\n", k, v, ok)

	k = passThroughKey
	v, ok = l.Value(k)
	fmt.Printf("test: Lookup(\"%v\") -> [value:%v] [ok:%v]\n", k, v, ok)

	k = nobodyHomeKey
	v, ok = l.Value(k)
	fmt.Printf("test: Lookup(\"%v\") -> [value:%v] [ok:%v]\n", k, v, ok)

	// set function
	l.SetOverride(map[string]string{passThroughKey: "value-pass-through"})
	k = ""
	v, ok = l.Value(k)
	fmt.Printf("test: Lookup(\"%v\") -> [value:%v] [ok:%v]\n", k, v, ok)

	k = passThroughKey
	v, ok = l.Value(k)
	fmt.Printf("test: Lookup(\"%v\") -> [value:%v] [ok:%v]\n", k, v, ok)

	k = nobodyHomeKey
	v, ok = l.Value(k)
	fmt.Printf("test: Lookup(\"%v\") -> [value:%v] [ok:%v]\n", k, v, ok)

	l.SetOverride(map[string]string{passThroughKey: "value-pass-through", nobodyHomeKey: "value-nobody-home"})
	k = passThroughKey
	v, ok = l.Value(k)
	fmt.Printf("test: Lookup(\"%v\") -> [value:%v] [ok:%v]\n", k, v, ok)

	k = nobodyHomeKey
	v, ok = l.Value(k)
	fmt.Printf("test: Lookup(\"%v\") -> [value:%v] [ok:%v]\n", k, v, ok)

	//Output:
	//test: Lookup("") -> [value:] [ok:false]
	//test: Lookup("pass through") -> [value:] [ok:false]
	//test: Lookup("nobody home") -> [value:] [ok:false]
	//test: Lookup("") -> [value:] [ok:false]
	//test: Lookup("pass through") -> [value:value-pass-through] [ok:true]
	//test: Lookup("nobody home") -> [value:] [ok:false]
	//test: Lookup("pass through") -> [value:value-pass-through] [ok:true]
	//test: Lookup("nobody home") -> [value:value-nobody-home] [ok:true]

}
