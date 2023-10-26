package exchange

import "fmt"

func Example_Resolve() {
	var s = ""
	url := Resolve(s)

	fmt.Printf("test: Resolve(%v) -> [%v]\n", s, url)

	s = "http://"
	url = Resolve(s)
	fmt.Printf("test: Resolve(%v) -> [%v]\n", s, url)

	s = "/test/resource?env=dev&cust=1"
	url = Resolve(s)
	fmt.Printf("test: Resolve(%v) -> [%v]\n", s, url)

	s = "https://www.google.com/search?q=testing"
	url = Resolve(s)
	fmt.Printf("test: Resolve(%v) -> [%v]\n", s, url)

	//Output:
	//test: Resolve() -> []
	//test: Resolve(http://) -> [http://]
	//test: Resolve(/test/resource?env=dev&cust=1) -> [http://localhost:8080/test/resource?env=dev&cust=1]
	//test: Resolve(https://www.google.com/search?q=testing) -> [https://www.google.com/search?q=testing]

}
