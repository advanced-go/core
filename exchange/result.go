package exchange

import "errors"

var okResult = new(Result)

type Result struct {
	Error    error
	Location string
}

func NewResultOK() *Result {
	return okResult
}

func NewResult(err error, loc string) *Result {
	r := new(Result)
	r.Error = err
	r.Location = loc
	return r
}

func (r *Result) OK() bool {
	return r.Error == nil
}

func (r *Result) String() string {
	if r.Error != nil {
		return r.Error.Error()
	}
	return "OK"
}

func testResult() {
	r := new(Result)
	r.Error = errors.New("this is an error string")
	r.Location = "this is the location."
}
