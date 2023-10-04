package runtime

// IfElseOp - templated function implementation of "C" ternary operator : conditional ? [true value] : [false value]
func IfElseOp[T any](cond bool, trueT T, falseT T) T {
	if cond {
		return trueT
	}
	return falseT
}
