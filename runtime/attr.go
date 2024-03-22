package runtime

import (
	"fmt"
	"strconv"
	"time"
)

type Attr struct {
	Key string
	Val any
}

func (a Attr) String() string {
	if s, ok := a.Val.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", a.Val)
}

type Attributable interface {
	Attributes() []Attr
}

func Attributes(e error) []Attr {
	if e == nil {
		return nil
	}
	if a, ok := any(e).(Attributable); ok {
		return a.Attributes()
	}
	return nil
}

func formatAttrs(attrs []any) string {
	if len(attrs) == 0 || attrs[0] == nil {
		return ""
	}
	result := ""
	value := ""
	name := ""
	for i := 0; i < len(attrs)-1; i += 2 {
		if i != 0 {
			result += ", "
		}
		name = fmt.Sprintf("%v", attrs[i])
		if n, ok := attrs[i+1].(int); ok {
			value = strconv.Itoa(n)
			result += jsonMarkup(name, value, false)
			continue
		}
		if b, ok1 := attrs[i+1].(bool); ok1 {
			if b {
				value = "true"
			} else {
				value = "false"
			}
			result += jsonMarkup(name, value, false)
			continue
		}
		if t, ok2 := attrs[i+1].(time.Time); ok2 {
			result += jsonMarkup(name, FmtRFC3339Millis(t), true)
			continue
		}
		value = fmt.Sprintf("%v", attrs[i+1])
		result += jsonMarkup(name, value, true)
	}
	if (len(attrs) & 0x01) == 1 {
		result += ", "
		name = fmt.Sprintf("%v", attrs[len(attrs)-1])
		result += jsonMarkup(name, "", true)
	}
	return result
}
