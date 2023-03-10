package runtime

// IfElse - templated function implementation of "C" ternary operator : conditional ? [true value] : [false value]
func IfElse[T any](cond bool, true T, false T) T {
	if cond {
		return true
	}
	return false
}
