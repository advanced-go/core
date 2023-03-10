package runtime

import (
	"fmt"
	"log"
)

// OutputHandler - template parameter output handler interface
type OutputHandler interface {
	Write(s string)
}

type NilOutput struct{}

func (NilOutput) Write(s string) {}

type StdOutput struct{}

func (StdOutput) Write(s string) { fmt.Println(s) }

type LogOutput struct{}

func (LogOutput) Write(s string) { log.Println(s) }

type TestOutput struct{}

func (TestOutput) Write(s string) { fmt.Printf("test: Write() -> [%v]\n", s) }
