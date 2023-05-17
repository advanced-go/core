package runtime

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	comment       = "//"
	mapDelimiter  = ":"
	pairDelimiter = ","
)

type TextPair struct {
	Key   string
	Value string
}

// ValidateMap - validates a configuration map, iterating through all keys
func ValidateMap(m map[string]string, err error, keys ...string) (errs []error) {
	if m == nil {
		return []error{errors.New("config map is nil")}
	}
	if err != nil {
		errs = append(errs, errors.New(fmt.Sprintf("config map read error: %v", err)))
	}
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if v == "" {
				errs = append(errs, errors.New(fmt.Sprintf("config map error: value for key does not exist [%v]", k)))
			}
		} else {
			errs = append(errs, errors.New(fmt.Sprintf("[config map error: key does not exist [%v]", k)))
		}
	}
	return
}

// TextToList - create a slice from a []byte
func TextToList(buf []byte) []string {
	var list []string

	if len(buf) == 0 {
		return list
	}
	r := bytes.NewReader(buf)
	reader := bufio.NewReader(r)
	var line string
	var err error
	for {
		line, err = reader.ReadString('\n')
		k := parseLine(line)
		if len(k) > 0 {
			list = append(list, k)
		}
		if err == io.EOF {
			break
		} else {
			if err != nil {
				break
			}
		}
	}
	return list
}

func ListToTextPair(s []string) []TextPair {
	var pair []TextPair
	for _, line := range s {
		key, val, err := splitLine(removeCrLf(line), pairDelimiter)
		if err == nil {
			pair = append(pair, TextPair{strings.TrimSpace(key), strings.TrimLeft(val, " ")})
		} else {
			pair = append(pair, TextPair{key, val})
		}
	}
	return pair
}

// TextToMap - create a map from a []byte
func TextToMap(buf []byte) (map[string]string, error) {
	m := make(map[string]string)
	if len(buf) == 0 {
		return m, nil
	}
	r := bytes.NewReader(buf)
	reader := bufio.NewReader(r)
	var line string
	var err error
	for {
		line, err = reader.ReadString('\n')
		k, v, err0 := parseMapLine(line)
		if err0 != nil {
			return m, err0
		}
		if len(k) > 0 {
			m[k] = v
		}
		if err == io.EOF {
			break
		} else {
			if err != nil {
				break
			}
		}
	}
	return m, nil
}

func parseLine(line string) string {
	if isEmpty(line) || isComment(line) {
		return ""
	}
	return removeCrLf(line)
}

func parseMapLine(line string) (string, string, error) {
	if len(line) == 0 {
		return "", "", nil
	}
	line = strings.TrimLeft(line, " ")
	if isEmpty(line) || isComment(line) {
		return "", "", nil
	}
	key, val, err1 := splitLine(line, mapDelimiter)
	if err1 != nil {
		return "", "", fmt.Errorf("invalid argument : line does not contain the ':' delimeter : [%v]", line)
	}
	val = removeCrLf(val)
	return strings.TrimSpace(key), strings.TrimLeft(val, " "), nil
}

func isEmpty(line string) bool {
	return len(line) == 0 || line == "" || line == "\r\n" || line == "\n"
}

func isComment(line string) bool {
	return strings.Index(line, comment) != -1
}

func removeCrLf(s string) string {
	index := strings.Index(s, "\r")
	if index != -1 {
		s = s[:index]
	}
	index = strings.Index(s, "\n")
	if index != -1 {
		s = s[:index]
	}
	return s
}

func splitLine(line, substring string) (string, string, error) {
	var err error
	var key string
	var val string

	i := strings.Index(line, substring)
	if i == -1 {
		err = fmt.Errorf("invalid argument : line does not contain the ':' delimeter : [%v]", line)
		key = line
	} else {
		key = line[:i]
		val = line[i+1:]
	}
	return key, val, err
}
