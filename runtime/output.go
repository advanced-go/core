package runtime

import (
	"fmt"
	"log"
)

// OutputHandler - template parameter output handler interface
type OutputHandler interface {
	Write(t any)
}

type NilOutput struct{}

func (NilOutput) Write(t any) {}

type StdOutput struct{}

func (StdOutput) Write(t any) { fmt.Println(t) }

type LogOutput struct{}

func (LogOutput) Write(t any) { log.Println(t) }

type TestOutput struct{}

func (TestOutput) Write(t any) { fmt.Printf("test: Write() -> [%v]\n", t) }
