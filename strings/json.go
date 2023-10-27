package strings

import (
	"fmt"
	"strings"
)

const (
	markupNull   = "\"%v\":null"
	markupString = "\"%v\":\"%v\""
	markupValue  = "\"%v\":%v"
)

func WriteJsonMarkup(sb *strings.Builder, name, value string, stringValue bool) {
	if sb.Len() == 0 {
		sb.WriteString("{")
	} else {
		sb.WriteString(",")
	}
	if len(value) == 0 {
		sb.WriteString(fmt.Sprintf(markupNull, name))
	} else {
		format := markupString
		if !stringValue {
			format = markupValue
		}
		sb.WriteString(fmt.Sprintf(format, name, value))
	}
}

func JsonMarkup(name, value string, stringValue bool) string {
	if len(value) == 0 {
		return fmt.Sprintf(markupNull, name)
	}
	format := markupString
	if !stringValue {
		format = markupValue
	}
	return fmt.Sprintf(format, name, value)
}
