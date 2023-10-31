package json

import "fmt"

func Example_Marshal() {
	var i = 45
	buf, status := Marshal(i)
	fmt.Printf("test: Marshal(int) [status:%v] %v\n", status, string(buf))

	var str []string
	buf, status = Marshal(str)
	fmt.Printf("test: Marshal([]string) [status:%v] %v\n", status, string(buf))

	var ptr *int
	buf, status = Marshal(ptr)
	fmt.Printf("test: Marshal(*int(nil)) [status:%v] %v\n", status, string(buf))

	ptr = &i
	buf, status = Marshal(ptr)
	fmt.Printf("test: Marshal(*int) [status:%v] %v\n", status, string(buf))

	//Output:
	//test: Marshal(int) [status:OK] 45
	//test: Marshal([]string) [status:OK] null
	//test: Marshal(*int(nil)) [status:OK] null
	//test: Marshal(*int) [status:OK] 45

}

func Example_Unmarshal() {
	var i = 45
	buf, status := Marshal(i)
	if status != nil {
	}

	var j int
	status = Unmarshal(buf, &j)
	fmt.Printf("test: Unmarshal(int) [status:%v] %v\n", status, j)

	var str = []string{"test", "of", "[]string"}
	var str2 []string

	buf, status = Marshal(str)
	status = Unmarshal(buf, &str2)
	fmt.Printf("test: Unmarshal([]string) [status:%v] %v\n", status, str2)

	//fmt.Printf("test: Marshal([]string) [status:%v] %v\n", status, string(buf))

	//Output:
	//test: Unmarshal(int) [status:OK] 45
	//test: Unmarshal([]string) [status:OK] [test of []string]

}

/*
func Example_Unmarshal() {
	var i = 45
	buf, status := Marshal(i)
	if status != nil {
	}

	j, status1 := Unmarshal[int](buf)
	fmt.Printf("test: Unmarshal(int) [status:%v] %v\n", status1, j)

	var str = []string{"test", "of", "[]string"}
	buf, status = Marshal(str)
	strs, status2 := Unmarshal[[]string](buf)
	fmt.Printf("test: Unmarshal([]string) [status:%v] %v\n", status2, strs)

	//fmt.Printf("test: MarshalType([]string) [status:%v] %v\n", status, string(buf))

	//Output:
	//test: Unmarshal(int) [status:OK] 45
	//test: Unmarshal([]string) [status:OK] [test of []string]

}

*/
