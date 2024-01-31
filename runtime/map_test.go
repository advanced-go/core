package runtime

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

const (
	configMapUri = "file://[cwd]/runtimetest/config-map.txt"
)

func ExampleStringsMap_Add() {
	smap := NewStringsMapFromHeader(nil)
	key1 := "key-1"

	status := smap.Add("", "")
	fmt.Printf("test: Add(\"\") -> [status:%v]\n", status)

	status = smap.Add(key1, "value-1")
	fmt.Printf("test: Add(%v) -> [status:%v]\n", key1, status)

	//Output:
	//test: Add("") -> [status:Invalid Argument [invalid argument: key is empty]]
	//test: Add(key-1) -> [status:OK]

}

func ExampleStringsMap_Get() {
	key1 := "key-1"
	key2 := "key-2"
	h := make(http.Header)
	h.Add(key1, "value-1")
	h.Add(key2, "value-2")

	smap := NewStringsMapFromHeader(h)

	val, status := smap.Get("")
	fmt.Printf("test: Get(\"\") -> [val:%v] [status:%v]\n", val, status)

	val, status = smap.Get(key1)
	fmt.Printf("test: Get(%v) -> [val:%v] [status:%v]\n", key1, val, status)

	val, status = smap.Get(key2)
	fmt.Printf("test: Get(%v) -> [val:%v] [status:%v]\n", key2, val, status)

	//Output:
	//test: Get("") -> [val:] [status:Invalid Argument [invalid argument: key does not exist: []]]
	//test: Get(key-1) -> [val:value-1] [status:OK]
	//test: Get(key-2) -> [val:value-2] [status:OK]

}

func ExampleNewStringsMap() {
	uri := configMapUri

	m := NewStringsMap(uri)
	fmt.Printf("test: NewStringsMap(\"%v\") -> [err:%v]\n", uri, m.Error())

	key := "user"
	val, status := m.Get(key)
	fmt.Printf("test: Get(\"%v\") -> [user:%v] [status:%v]\n", key, val, status)

	key = "pswd"
	val, status = m.Get(key)
	fmt.Printf("test: Get(\"%v\") -> [pswd:%v] [status:%v]\n", key, val, status)

	key = "uri"
	val, status = m.Get(key)
	fmt.Printf("test: Get(\"%v\") -> [urir:%v] [status:%v]\n", key, val, status)

	//Output:
	//test: NewStringsMap("file://[cwd]/runtimetest/config-map.txt") -> [err:<nil>]
	//test: Get("user") -> [user:bobs-your-uncle] [status:OK]
	//test: Get("pswd") -> [pswd:let-me-in] [status:OK]
	//test: Get("uri") -> [urir:postgres://{user}:{pswd}@{sub-domain}.{db-name}.cloud.timescale.com:31770/tsdb?sslmode=require] [status:OK]

}

func ExampleValidateMap() {
	m := map[string]string{"database-url": "postgres://{user}:{pswd}@{sub-domain}.{db-name}.cloud.timescale.com:31770/tsdb?sslmode=require", "ping-path": "", "postgres-urn": "urn:postgres", "postgres-pgxsql-uri": "github.com/idiomatic-go/postgresql/pgxsql"}
	errs := ValidateMap(nil, nil)
	fmt.Printf("test: ValidateConfig(nil,nil) -> %v\n", errs)

	errs = ValidateMap(m, errors.New("file I/O error"))
	fmt.Printf("test: ValidateConfig(m,err) -> %v\n", errs)

	errs = ValidateMap(m, nil, "not-found")
	fmt.Printf("test: Validate(m,nil,not-found) -> %v\n", errs)

	errs = ValidateMap(m, nil, "database-url", "ping-path", "postgres-pgxsql-uri")
	fmt.Printf("test: Validate(m,nil,...) -> %v\n", errs)

	//Output:
	//test: ValidateConfig(nil,nil) -> [config map is nil]
	//test: ValidateConfig(m,err) -> [config map read error: file I/O error]
	//test: Validate(m,nil,not-found) -> [[config map error: key does not exist [not-found]]
	//test: Validate(m,nil,...) -> [config map error: value for key does not exist [ping-path]]

}

func ExampleParseLine() {
	key, val, err := parseLine("key : value\r\n")
	fmt.Printf("test: parseLine(cr,lf) -> [key:%v] [val:%v] [error:%v]\n", key, val, err)

	key, val, err = parseLine("key : value\n")
	fmt.Printf("test: parseLine(lf) -> [key:%v] [val:%v] [error:%v]\n", key, val, err)

	//Output:
	//test: parseLine(cr,lf) -> [key:key] [val:value] [error:<nil>]
	//test: parseLine(lf) -> [key:key] [val:value] [error:<nil>]

}

func TestParseLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"BlankLine", args{line: ""}, "", "", false},
		{"LeadingSpace", args{line: " "}, "", "", false},
		{"LeadingSpaces", args{line: "       "}, "", "", false},
		{"NewLine", args{line: "\n"}, "", "", false},

		{"Comment", args{line: comment}, "", "", false},
		{"LeadingSpaceComment", args{line: " " + comment}, "", "", false},
		{"LeadingSpacesComment", args{line: "       " + comment}, "", "", false},

		{"MissingDelimiter", args{line: "missing delimiter"}, "", "", true},

		{"KeyOnly", args{line: "key-only :"}, "key-only", "", false},
		{"KeyValue", args{line: "key  : value\r\n"}, "key", "value", false},
		{"KeyValueLeadingSpaces", args{line: "key:      value"}, "key", "value", false},
		{"KeyValueTrailingSpaces", args{line: "key :value    "}, "key", "value    ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseLine() got = [%v], want [%v]", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseLine() got1 = [%v], want [%v]", got1, tt.want1)
			}
		})
	}
}
