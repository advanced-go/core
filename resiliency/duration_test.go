package resiliency

import (
	"fmt"
	"time"
)

func Example_ExponentialDuration() {
	e := NewExponentialDuration(0.5, 5000)
	for i := 0; i < 15; i++ {
		fmt.Printf("test: Eval() -> %v %v\n", i, int(e.Eval())/int(time.Millisecond))
	}

	//Output:
}
