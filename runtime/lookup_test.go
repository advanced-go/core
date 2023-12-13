package runtime

import "fmt"

func ExampleLookup_Invalid() {
	value := "constant"
	ctx := NewLookupContext(nil, 45)

	result := LookupFromContext(ctx, value)
	fmt.Printf("test: LookupFromContext() -> [value:%v] [result:%v]\n", value, result)

	//Output:
	//test: LookupFromContext() -> [value:constant] [result:constant]

}

func ExampleLookup_String() {
	value := "value-0"
	ctx := NewLookupContext(nil, "value-1")

	result := LookupFromContext(ctx, value)
	fmt.Printf("test: LookupFromContext() -> [value:%v] [result:%v]\n", value, result)

	value = "not-in-lookup"
	result = LookupFromContext(ctx, value)
	fmt.Printf("test: LookupFromContext() -> [value:%v] [result:%v]\n", value, result)

	value = ""
	result = LookupFromContext(ctx, value)
	fmt.Printf("test: LookupFromContext() -> [value:%v] [result:%v]\n", value, result)

	//Output:
	//test: LookupFromContext() -> [value:value-0] [result:value-1]
	//test: LookupFromContext() -> [value:not-in-lookup] [result:value-1]
	//test: LookupFromContext() -> [value:] [result:value-1]

}

func ExampleLookup_Map() {
	value := "value-0"

	lookup := map[string]string{value: "value-1"}
	ctx := NewLookupContext(nil, lookup)

	result := LookupFromContext(ctx, value)
	fmt.Printf("test: LookupFromContext() -> [value:%v] [result:%v]\n", value, result)

	value = "not-in-lookup"
	result = LookupFromContext(ctx, value)
	fmt.Printf("test: LookupFromContext() -> [value:%v] [result:%v]\n", value, result)

	//Output:
	//test: LookupFromContext() -> [value:value-0] [result:value-1]
	//test: LookupFromContext() -> [value:not-in-lookup] [result:not-in-lookup]

}

func lookupFunc(k string) string {
	switch k {
	case "value-0":
		return "value-1"
	}
	return k
}

func ExampleLookup_Func() {
	value := "value-0"
	ctx := NewLookupContext(nil, lookupFunc)

	result := LookupFromContext(ctx, value)
	fmt.Printf("test: LookupFromContext() -> [value:%v] [result:%v]\n", value, result)

	value = "not-in-lookup"
	result = LookupFromContext(ctx, value)
	fmt.Printf("test: LookupFromContext() -> [value:%v] [result:%v]\n", value, result)

	//Output:
	//test: LookupFromContext() -> [value:value-0] [result:value-1]
	//test: LookupFromContext() -> [value:not-in-lookup] [result:not-in-lookup]

}
