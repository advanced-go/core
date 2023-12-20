package exchange

import "fmt"

func echo(key string) string {
	return key
}

func empty(key string) string {
	return ""
}

func Example_Resolver() {
	r := NewResolver("http://localhost:8080", echo)

	key := "test"
	val := r.Resolve(key)
	fmt.Printf("test: Resolver(\"%v\") -> %v\n", key, val)

	key = ""
	val = r.Resolve(key)
	fmt.Printf("test: Resolver(\"%v\") -> %v\n", key, val)

	key = "/"
	val = r.Resolve(key)
	fmt.Printf("test: Resolver(\"%v\") -> %v\n", key, val)

	r = NewResolver("http://localhost:8080", empty)
	key = "/"
	val = r.Resolve(key)
	fmt.Printf("test: Resolver(\"%v\") -> %v\n", key, val)

	//Output:
	//test: Resolver("test") -> test
	//test: Resolver("") ->
	//test: Resolver("/") -> /
	//test: Resolver("/") -> http://localhost:8080/
	
}
