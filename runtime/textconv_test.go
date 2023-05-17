package runtime

import (
	"errors"
	"fmt"
	"testing"
)

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

func ExampleParse() {
	line := "key\r\n"
	key, val, err := parseLine(line, false)
	fmt.Printf("test: parseLine(keyCrLf) -> [key:%v] [val:%v] [error:%v]\n", key, val, err)

	line = "key\n"
	key, val, err = parseLine(line, false)
	fmt.Printf("test: parseLine(keyLf) -> [key:%v] [val:%v] [error:%v]\n", key, val, err)

	line = "key : value\r\n"
	key, val, err = parseLine(line, true)
	fmt.Printf("test: parseLine(key : valueCrLf) -> [key:%v] [val:%v] [error:%v]\n", key, val, err)

	line = "key : value\n"
	key, val, err = parseLine(line, true)
	fmt.Printf("test: parseLine(key : valueLf) -> [key:%v] [val:%v] [error:%v]\n", key, val, err)

	//Output:
	//test: parseLine(keyCrLf) -> [key:key] [val:] [error:<nil>]
	//test: parseLine(keyLf) -> [key:key] [val:] [error:<nil>]
	//test: parseLine(key : valueCrLf) -> [key:key] [val:value] [error:<nil>]
	//test: parseLine(key : valueLf) -> [key:key] [val:value] [error:<nil>]

}

func TestTextToMap(t *testing.T) {
	type args struct {
		line  string
		isMap bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"BlankLine", args{line: "", isMap: true}, "", "", false},
		{"LeadingSpace", args{line: " ", isMap: true}, "", "", false},
		{"LeadingSpaces", args{line: "       ", isMap: true}, "", "", false},

		{"Comment", args{line: comment, isMap: true}, "", "", false},
		{"LeadingSpaceComment", args{line: " " + comment, isMap: true}, "", "", false},
		{"LeadingSpacesComment", args{line: "       " + comment, isMap: true}, "", "", false},

		{"MissingDelimiter", args{line: "missing delimiter", isMap: true}, "", "", true},

		{"KeyOnly", args{line: "key-only :", isMap: true}, "key-only", "", false},
		{"KeyValue", args{line: "key  : value\r\n", isMap: true}, "key", "value", false},
		{"KeyValueLeadingSpaces", args{line: "key:      value", isMap: true}, "key", "value", false},
		{"KeyValueTrailingSpaces", args{line: "key :value    ", isMap: true}, "key", "value    ", false},
		{"KeyValueEmptyLine", args{line: "key :value\r\n\r\n    ", isMap: true}, "key", "value", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseLine(tt.args.line, tt.args.isMap)
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

func TestTextToSlice(t *testing.T) {
	type args struct {
		line  string
		isMap bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"BlankLine", args{line: "", isMap: false}, "", "", false},
		{"LeadingSpace", args{line: " ", isMap: false}, "", "", false},
		{"LeadingSpaces", args{line: "       ", isMap: false}, "", "", false},

		{"Comment", args{line: comment, isMap: false}, "", "", false},
		{"LeadingSpaceComment", args{line: " " + comment, isMap: false}, "", "", false},
		{"LeadingSpacesComment", args{line: "       " + comment, isMap: false}, "", "", false},

		{"KeyOnly", args{line: "key-only :", isMap: false}, "key-only :", "", false},
		{"KeyValue", args{line: "key  : value\r\n", isMap: false}, "key  : value", "", false},
		{"KeyValueEmptyLine", args{line: "key :value\r\n\r\n    ", isMap: false}, "key :value", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseLine(tt.args.line, tt.args.isMap)
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
