package runtimetest

func stringDefault(key string) string {
	return "default"
}

func Example_LookupString() {
	l := NewLookup[string, func(string) string](stringDefault)

	l.Resolve("")

}

func listDefault(key string) []string {
	return []string{"value-0", "value-1"}
}

func Example_LookupList() {
	l := NewLookup[[]string, func(string) []string](listDefault)

	l.Resolve("")

}

/*

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


*/

/*

func Example_LookupFromType_List() {
	fn := LookupFromType[func(string) []string](nil)
	//fmt.Printf("test: LookupFromType(\"\") -> [%v]\n", fn)

	v := "string constant"
	k := "random"
	fn = LookupFromType[func(string) []string](v)
	fmt.Printf("test: LookupFromType(\"\") -> [%v]\n", fn(""))
	fmt.Printf("test: LookupFromType(%v) -> [%v]\n", k, fn(""))


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


	//Output:
	//test: LookupFromType("") -> [[string constant ]]
	//test: LookupFromType(random) -> [[string constant ]]

}


*/
