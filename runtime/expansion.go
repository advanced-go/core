package runtime

import (
	"errors"
	"fmt"
	"strings"
)

const (
	BeginDelimiter = "{"
	EndDelimiter   = "}"
)

// Resolver - template parameter name value lookup
type Resolver interface {
	Lookup(name string) (string, error)
}

// Expand - templated function to expand a template string, utilizing a resolver
func Expand[T Resolver](t string) (string, error) {
	var r T

	if t == "" {
		return "", nil
	}
	var buf strings.Builder
	tokens := strings.Split(t, BeginDelimiter)
	if len(tokens) == 1 {
		return t, nil
	}
	for _, s := range tokens {
		sub := strings.Split(s, EndDelimiter)
		if len(sub) > 2 {
			return "", errors.New(fmt.Sprintf("invalid argument : token has multiple end delimiters: %v", s))
		}
		// Check case where no end delimiter is found
		if len(sub) == 1 && sub[0] == s {
			buf.WriteString(s)
			continue
		}
		// Have a valid end delimiter, so lookup the variable
		t, err := r.Lookup(sub[0])
		if err != nil {
			return "", err
		}
		buf.WriteString(t)
		if len(sub) == 2 {
			buf.WriteString(sub[1])
		}
	}
	return buf.String(), nil
}
