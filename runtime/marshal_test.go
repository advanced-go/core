package runtime

import "fmt"

func Example_MarshalType() {
	var i = 45
	buf, status := MarshalType(i)
	fmt.Printf("test: MarshalType(int) [status:%v] %v\n", status, string(buf))

	var str []string
	buf, status = MarshalType(str)
	fmt.Printf("test: MarshalType([]string) [status:%v] %v\n", status, string(buf))

	var ptr *int
	buf, status = MarshalType(ptr)
	fmt.Printf("test: MarshalType(*int(nil)) [status:%v] %v\n", status, string(buf))

	ptr = &i
	buf, status = MarshalType(ptr)
	fmt.Printf("test: MarshalType(*int) [status:%v] %v\n", status, string(buf))

	//Output:
	//test: MarshalType(int) [status:OK] 45
	//test: MarshalType([]string) [status:OK] null
	//test: MarshalType(*int(nil)) [status:OK] null
	//test: MarshalType(*int) [status:OK] 45

}

func Example_UnmarshalType() {
	var i = 45
	buf, status := MarshalType(i)
	if status != nil {
	}

	j, status1 := UnmarshalType[int](buf)
	fmt.Printf("test: UnmarshalType(int) [status:%v] %v\n", status1, j)

	var str = []string{"test", "of", "[]string"}
	buf, status = MarshalType(str)
	strs, status2 := UnmarshalType[[]string](buf)
	fmt.Printf("test: UnmarshalType([]string) [status:%v] %v\n", status2, strs)

	//fmt.Printf("test: MarshalType([]string) [status:%v] %v\n", status, string(buf))

	/*
		var ptr *int
		buf, status = MarshalType[DebugError](ptr)
		fmt.Printf("test: MarshalType(*int(nil)) [status:%v] %v\n", status, string(buf))

		ptr = &i
		buf, status = MarshalType[DebugError](ptr)
		fmt.Printf("test: MarshalType(*int) [status:%v] %v\n", status, string(buf))


	*/

	//Output:
	//test: UnmarshalType(int) [status:OK] 45
	//test: UnmarshalType([]string) [status:OK] [test of []string]

}
