package uri

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	EnvPrefix = "$"
)

// LookupEnv - environment lookup
func LookupEnv(name string) (string, error) {
	if strings.HasPrefix(name, EnvPrefix) {
		return os.Getenv(name[len(EnvPrefix):]), nil
	}
	return "", errors.New(fmt.Sprintf("invalid argument : LookupEnv() template variable is invalid: %v", name))
}
