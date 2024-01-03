package runtimetest

import "fmt"

func stringer() string {
	s := "in stringer()"
	fmt.Printf("%v\n", s)
	return s
}

type stringerFunc func() string

func (f stringerFunc) String() string {
	return f()
}
