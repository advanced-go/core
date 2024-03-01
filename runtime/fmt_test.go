package runtime

import (
	"fmt"
	"time"
)

const (
	timeStampString = "2024-03-01T18:23:50.205Z"
)

var testTS time.Time

func init() {
	testTS = time.Date(2024, 3, 1, 18, 23, 50, 205*1e6, time.UTC)

}

func ExampleParseYMD() {
	y, m, d, err := parseYMD(timeStampString)

	fmt.Printf("test: ParseYMD(\"%v\") -> [year:%v] [month:%v] [day:%v] [err:%v]\n", timeStampString, y, m, d, err)

	//Output:
	//test: ParseYMD("2024-03-01T18:23:50.205Z") -> [year:2024] [month:3] [day:1] [err:<nil>]

}

func ExampleParseHMSM() {
	h, m, s, ms, err := parseHMSM(timeStampString)

	fmt.Printf("test: ParseHMSM(\"%v\") -> [hour:%v] [min:%v] [sec:%v] [ms:%v] [err:%v]\n", timeStampString, h, m, s, ms, err)

	//Output:
	//test: ParseHMSM("2024-03-01T18:23:50.205Z") -> [hour:18] [min:23] [sec:50] [ms:205] [err:<nil>]

}

func ExampleParseTimestamp() {
	t2, err := ParseRFC3339Millis(timeStampString)
	s := FmtRFC3339Millis(t2)
	fmt.Printf("test: ParseTimestamp(\"%v\") -> [%v] [err:%v] [equal:%v]\n", timeStampString, err, s, s == timeStampString)

	//Output:
	//test: ParseTimestamp("2024-03-01T18:23:50.205Z") -> [<nil>] [err:2024-03-01T18:23:50.205Z] [equal:true]

}
