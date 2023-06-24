package runtime

import "fmt"

// OutputHandler2 - template parameter output handler interface
type OutputHandler2[T any] interface {
	Write(t T)
}

type StdOutput2 struct{}

func (StdOutput2) Write(s string) { fmt.Println(s) }
