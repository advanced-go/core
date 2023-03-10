package runtime

import (
	"fmt"
)

func ExamplePrimer() {
	current := true

	result := IfElse[bool](current, current, true)
	fmt.Printf("test: T[bool](true,false,true) -> %v\n", result)

	result = IfElse[bool](!current, current, true)
	fmt.Printf("test: T[bool](false,false,true) -> %v\n", result)

	//Output:
	//test: T[bool](true,false,true) -> true
	//test: T[bool](false,false,true) -> true
}

func ExampleInt() {
	bval := false
	num := 99
	var num2 int

	num2 = IfElse[int](bval, 45, 145)
	if num2 > 0 {

	}

	result := IfElse[int](num < 100, 45, 145)
	fmt.Printf("test: T[int](num<100, 45, 145) -> [cond:%v] [result:%v]\n", num < 100, result)

	result = IfElse[int](num >= 100, 55, 155)
	fmt.Printf("test: T[int](num>=100, 45, 145) -> [cond:%v] [result:%v]\n", num >= 100, result)

	//Output:
	//test: T[int](num<100, 45, 145) -> [cond:true] [result:45]
	//test: T[int](num>=100, 45, 145) -> [cond:false] [result:155]

}

/*
func ExampleContext() {
	ctx := ContextWithContent(context.Background(), "new context")

	result := T[context.Context](ctx != nil, ctx, context.Background())
	fmt.Printf("test: T[context.Context](ctx!=nil,ctx,context.Background())) -> [cond:%v] [background:%v]\n", ctx != nil, result == context.Background())

	result = T[context.Context](ctx == nil, ctx, context.Background())
	fmt.Printf("test: T[context.Context](ctx==nil,ctx,context.Background()) -> [cond:%v] [background:%v]\n", ctx == nil, result == context.Background())

	//Output:
	//test: T[context.Context](ctx!=nil,ctx,context.Background())) -> [cond:true] [background:false]
	//test: T[context.Context](ctx==nil,ctx,context.Background()) -> [cond:false] [background:true]

}


*/
