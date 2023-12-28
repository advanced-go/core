package runtime

type runtimeEnv int

const (
	debug runtimeEnv = iota
	test
	stage
	production
)

var (
	rte = debug
)

func IsProdEnvironment() bool {
	return rte == production
}

func SetProdEnvironment() {
	rte = production
}

func IsTestEnvironment() bool {
	return rte == test
}

func SetTestEnvironment() {
	rte = test
}

func IsStageEnvironment() bool {
	return rte == stage
}

func SetStageEnvironment() {
	rte = stage
}

func IsDebugEnvironment() bool {
	return rte == debug
}

// EnvStr - string representation for the environment
func EnvStr() string {
	switch rte {
	case debug:
		return "debug"
	case test:
		return "test"
	case stage:
		return "stage"
	case production:
		return "prod"
	}
	return "unknown"
}
