package runtime

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	stringsMapAdd = PkgPath + ":Add"
	stringsMapGet = PkgPath + ":Get"
	errorKey      = "error"
	comment       = "//"
	delimiter     = ":"
	newLine       = "\n"
	cr            = "\r"
)

// StringsMap - key value pairs of string -> string
type StringsMap struct {
	m *sync.Map
}

func NewEmptyStringsMap() *StringsMap {
	m := new(StringsMap)
	m.m = new(sync.Map)
	return m
}

func NewStringsMap(uri string) *StringsMap {
	buf, err := os.ReadFile(FileName(uri))
	if err == nil {
		return ParseStringsMap(buf)
	}
	m := NewEmptyStringsMap()
	m.Add(errorKey, err.Error())
	return m
}

func NewStringsMapFromHeader(h http.Header) *StringsMap {
	m := NewEmptyStringsMap()
	if h != nil {
		for k, v := range h {
			if len(v) > 0 {
				m.Add(strings.ToLower(k), v[0])
			}
		}
	}
	return m
}

// Error - determine if there is an error
func (m *StringsMap) Error() error {
	v, status := m.Get(errorKey)
	if !status.OK() {
		return nil
	}
	return errors.New(v)
}

// Add - add a value
func (m *StringsMap) Add(key, val string) Status {
	if len(key) == 0 {
		return NewStatusError(StatusInvalidArgument, stringsMapAdd, errors.New("invalid argument: key is empty"))
	}
	_, ok1 := m.m.Load(key)
	if ok1 {
		return NewStatusError(StatusInvalidArgument, stringsMapAdd, errors.New(fmt.Sprintf("invalid argument: key already exists: [%v]", key)))
	}
	m.m.Store(key, val)
	return StatusOK()
}

// Get - get a value
func (m *StringsMap) Get(key string) (string, Status) {
	v, ok := m.m.Load(key)
	if !ok {
		return "", NewStatusError(StatusInvalidArgument, stringsMapGet, errors.New(fmt.Sprintf("invalid argument: key does not exist: [%v]", key)))
	}
	if val, ok1 := v.(string); ok1 {
		return val, StatusOK()
	}
	return "", NewStatus(StatusInvalidContent)
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

// ParseStringsMap - create a strings map from a []byte
func ParseStringsMap(buf []byte) *StringsMap {
	m := NewEmptyStringsMap()
	if len(buf) == 0 {
		m.Add(errorKey, "errors: Unable to parse strings map: buf is empty")
		return m
	}
	r := bytes.NewReader(buf)
	reader := bufio.NewReader(r)
	var line string
	var err error
	var count = 0
	for {
		line, err = reader.ReadString('\n')
		count++
		k, v, err0 := parseLine(line)
		if err0 != nil {
			m.Add(errorKey, fmt.Sprintf("%v : line -> %v number -> %v", err0.Error(), line, count))
			return m
		}
		if len(k) > 0 {
			m.Add(k, v)
		}
		if err == io.EOF {
			break
		} else {
			if err != nil {
				break
			}
		}
	}
	return m
}

func parseLine(line string) (string, string, error) {
	if len(line) == 0 {
		return "", "", nil
	}
	line = strings.TrimLeft(line, " ")
	if len(line) == 0 || strings.HasPrefix(line, comment) {
		return "", "", nil
	}
	if line[:1] == cr || line[:1] == newLine {
		return "", "", nil
	}
	i := strings.Index(line, delimiter)
	if i == -1 {
		return "", "", fmt.Errorf("invalid argument : line does not contain the ':' delimeter : [%v]", line)
	}
	key := line[:i]
	val := line[i+1:]
	index := strings.Index(val, "\r")
	if index != -1 {
		val = val[:index]
	}
	index = strings.Index(val, "\n")
	if index != -1 {
		val = val[:index]
	}
	return strings.TrimSpace(key), strings.TrimLeft(val, " "), nil
}
