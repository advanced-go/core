package access

import (
	"fmt"
	"time"
)

func ExampleFmtTimestamp() {
	t := time.Now().UTC()
	s := FmtTimestamp(t)
	fmt.Printf("test: FmtTimestamp() -> [%v]\n", s)

	//t2, err := ParseTimestamp(s)
	//fmt.Printf("test: ParseTimestamp() -> [%v] [%v]\n", FmtTimestamp(t2), err)

	//Output:

}

func ExampleFmtRFC3339Millis() {
	t := time.Now().UTC()
	s := FmtRFC3339Millis(t)

	fmt.Printf("test: ExampleFmtRFC3339Millis -> [%v]\n", s)

	//Output:

}
