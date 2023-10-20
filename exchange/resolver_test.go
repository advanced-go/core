package exchange

import "fmt"

func Example_ResolveUrl() {
	var s = ""
	url := ResolveUrl(s)

	fmt.Printf("test: ResolveUrl(%v) -> [%v]\n", s, url)

	s = "http://"
	url = ResolveUrl(s)
	fmt.Printf("test: ResolveUrl(%v) -> [%v]\n", s, url)

	s = "/test/resource?env=dev&cust=1"
	url = ResolveUrl(s)
	fmt.Printf("test: ResolveUrl(%v) -> [%v]\n", s, url)

	s = "https://www.google.com/search?q=testing"
	url = ResolveUrl(s)
	fmt.Printf("test: ResolveUrl(%v) -> [%v]\n", s, url)

	//Output:
	//test: ResolveUrl() -> []
	//test: ResolveUrl(http://) -> [http://]
	//test: ResolveUrl(/test/resource?env=dev&cust=1) -> [http://localhost:8080/test/resource?env=dev&cust=1]
	//test: ResolveUrl(https://www.google.com/search?q=testing) -> [https://www.google.com/search?q=testing]
	
}
