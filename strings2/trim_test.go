package strings2

import "fmt"

func ExampleTrim() {
	left := "left-spaces"
	ls := "   " + left
	fmt.Printf("test: Trim(\"%v\") -> %v\n", ls, "\""+TrimSpace(ls)+"\"")

	right := "right-spaces"
	rs := right + "   "
	fmt.Printf("test: Trim(\"%v\") -> %v\n", rs, "\""+TrimSpace(rs)+"\"")

	both := "both-spaces"
	bs := "  " + both + "  "
	fmt.Printf("test: Trim(\"%v\") -> %v\n", bs, "\""+TrimSpace(bs)+"\"")

	//Output:
	//test: Trim("   left-spaces") -> "left-spaces"
	//test: Trim("right-spaces   ") -> "right-spaces"
	//test: Trim("  both-spaces  ") -> "both-spaces"

}
