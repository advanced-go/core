package runtime

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	EnvPrefix = "$"
)

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

func LookupEnv(name string) (string, error) {
	if strings.HasPrefix(name, EnvPrefix) {
		return os.Getenv(name[len(EnvPrefix):]), nil
	}
	return "", errors.New(fmt.Sprintf("invalid argument : LookupEnv() template variable is invalid: %v", name))
}
