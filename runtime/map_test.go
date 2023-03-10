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
