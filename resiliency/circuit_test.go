package resiliency

import (
	"errors"
	"fmt"
)

func Example_CircuitBreaker() {
	err := errors.New("test error")
	fmt.Printf("test: CircuitBreaker() -> %v\n", err)

	//Output:
}
