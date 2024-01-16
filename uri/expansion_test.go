package uri

import (
	"errors"
	"fmt"
	"strings"
)

func lookup(name string) (string, error) {
	switch strings.ToLower(name) {
	case "env":
		return "[ENV]", nil
	case "next":
		return "[NEXT]", nil
	case "last":
		return "[LAST]", nil
	}
	return "", errors.New(fmt.Sprintf("invalid argument : template variable is invalid: %v", name))
}

func ExampleExpand_InvalidLookup() {
	// Lookup name not found
	s := "test{not-valid}"
	_, err := Expand(s, lookup)
	fmt.Printf("test: Expand(%v) ->  : %v\n", s, err)

	//Output:
	//test: Expand(test{not-valid}) ->  : invalid argument : template variable is invalid: not-valid

}

func ExampleExpand_InvalidDelimiters() {
	var err error
	// Mismatched delimiters - too many end delimiters
	s := "resources/test-file-name{env}}and{next}{last}.txt"
	_, err = Expand(s, lookup)

	fmt.Printf("Path Input  : %v\n", s)
	fmt.Printf("Path Output : %v\n", err)

	// Mismatched delimiters - too many begin delimiters, this is valid as the extra begin delimiters are skipped
	s = "resources/test-file-name{env}and{next}{{last}.txt"
	path, err0 := Expand(s, lookup)

	fmt.Printf("Path Input  : %v\n", s)
	fmt.Printf("Path Output : %v %v\n", path, err0)

	// Mismatched delimiters - embedded begin delimiter
	s = "resources/test-file-name{env}and{next{}{last}.txt"
	path, err0 = Expand(s, lookup)
	fmt.Printf("Path Input  : %v\n", s)
	fmt.Printf("Path Output : %v %v\n", path, err0)

	//Output:
	//Path Input  : resources/test-file-name{env}}and{next}{last}.txt
	//Path Output : invalid argument : token has multiple end delimiters: env}}and
	//Path Input  : resources/test-file-name{env}and{next}{{last}.txt
	//Path Output : resources/test-file-name[ENV]and[NEXT][LAST].txt <nil>
	//Path Input  : resources/test-file-name{env}and{next{}{last}.txt
	//Path Output :  invalid argument : template variable is invalid:
}

func ExampleExpand_Valid() {
	s := ""

	path, err := Expand(s, lookup)
	fmt.Printf("test: Expand() -> [error:%v] [path:%v]\n", err, path)

	s = "resources/test-file-name-and-ext.txt"
	path, err = Expand(s, lookup)
	fmt.Printf("test: Expand(%v) -> [error:%v] [path:%v]\n", s, err, path)

	s = "resources/test-file-name{env}and{next}{last}.txt"
	path, err = Expand(s, lookup)
	fmt.Printf("test: Expand(%v) -> [error:%v] [path:%v]\n", s, err, path)

	s = "resources/test-file-name_{env}.txt"
	path, err = Expand(s, lookup)
	fmt.Printf("test: Expand(%v) -> [error:%v] [path:%v]\n", s, err, path)

	//Output:
	//test: Expand() -> [error:<nil>] [path:]
	//test: Expand(resources/test-file-name-and-ext.txt) -> [error:<nil>] [path:resources/test-file-name-and-ext.txt]
	//test: Expand(resources/test-file-name{env}and{next}{last}.txt) -> [error:<nil>] [path:resources/test-file-name[ENV]and[NEXT][LAST].txt]
	//test: Expand(resources/test-file-name_{env}.txt) -> [error:<nil>] [path:resources/test-file-name_[ENV].txt]

}

func Example_TemplateToken() {
	t := ""

	v, ok := TemplateToken(t)
	fmt.Printf("test: TemplateToken(\"\") -> [var:%v] [ok:%v]\n", v, ok)

	t = "variable-name"
	v, ok = TemplateToken(t)
	fmt.Printf("test: TemplateToken(\"%v\") -> [var:%v] [ok:%v]\n", t, v, ok)

	t = "{variable-name"
	v, ok = TemplateToken(t)
	fmt.Printf("test: TemplateToken(\"%v\") -> [var:%v] [ok:%v]\n", t, v, ok)

	t = "variable-name}"
	v, ok = TemplateToken(t)
	fmt.Printf("test: TemplateToken(\"%v\") -> [var:%v] [ok:%v]\n", t, v, ok)

	t = "{variable-name}"
	v, ok = TemplateToken(t)
	fmt.Printf("test: TemplateToken(\"%v\") -> [var:%v] [ok:%v]\n", t, v, ok)

	//Output:
	//test: TemplateToken("") -> [var:] [ok:false]
	//test: TemplateToken("variable-name") -> [var:variable-name] [ok:false]
	//test: TemplateToken("{variable-name") -> [var:{variable-name] [ok:false]
	//test: TemplateToken("variable-name}") -> [var:variable-name}] [ok:false]
	//test: TemplateToken("{variable-name}") -> [var:variable-name] [ok:true]
	
}
