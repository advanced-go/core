package sql

import (
	"errors"
	"fmt"
)

const (
	QueryNSS  = "query"
	InsertNSS = "insert"
	UpdateNSS = "update"
	DeleteNSS = "delete"
	PingNSS   = "ping"
	StatNSS   = "stat"

	SelectCmd = 0
	InsertCmd = 1
	UpdateCmd = 2
	DeleteCmd = 3

	NullExpectedCount = int64(-1)
)

func BuildUri(nid, region, zone, nss, resource string) string {
	return fmt.Sprintf("urn:%v.%v.%v:%v.%v", nid, region, zone, nss, resource)
}

// BuildQueryUri - build an uri with the Query NSS
func BuildQueryUri(nid, region, zone, resource string) string {
	return BuildUri(nid, region, zone, QueryNSS, resource)
}

// BuildInsertUri - build an uri with the Insert NSS
func BuildInsertUri(nid, region, zone, resource string) string {
	return BuildUri(nid, region, zone, InsertNSS, resource)
}

// BuildUpdateUri - build an uri with the Update NSS
func BuildUpdateUri(nid, region, zone, resource string) string {
	return BuildUri(nid, region, zone, UpdateNSS, resource)
}

// BuildDeleteUri - build an uri with the Delete NSS
func BuildDeleteUri(nid, region, zone, resource string) string {
	return BuildUri(nid, region, zone, DeleteNSS, resource)
}

// Request - contains data needed to build the SQL statement related to the uri
type Request struct {
	ExpectedCount int64
	Cmd           int
	Uri           string
	Template      string
	Values        [][]any
	Attrs         []Attr
	Where         []Attr
	Error         error
}

func (r *Request) Validate() error {
	if r.Uri == "" {
		return errors.New("invalid argument: request Uri is empty")
	}
	if r.Template == "" {
		return errors.New("invalid argument: request template is empty")
	}
	return nil
}

func (r *Request) String() string {
	return r.Template
}

func NewQueryRequest(uri, template string, where []Attr) *Request {
	return &Request{ExpectedCount: NullExpectedCount, Cmd: SelectCmd, Uri: uri, Template: template, Where: where}
}

func NewQueryRequestFromValues(uri, template string, values map[string][]string) *Request {
	return &Request{ExpectedCount: NullExpectedCount, Cmd: SelectCmd, Uri: uri, Template: template, Where: BuildWhere(values)}
}

func NewInsertRequest(uri, template string, values [][]any) *Request {
	return &Request{ExpectedCount: NullExpectedCount, Cmd: InsertCmd, Uri: uri, Template: template, Values: values}
}

func NewUpdateRequest(uri, template string, attrs []Attr, where []Attr) *Request {
	return &Request{ExpectedCount: NullExpectedCount, Cmd: UpdateCmd, Uri: uri, Template: template, Attrs: attrs, Where: where}
}

func NewDeleteRequest(uri, template string, where []Attr) *Request {
	return &Request{ExpectedCount: NullExpectedCount, Cmd: DeleteCmd, Uri: uri, Template: template, Attrs: nil, Where: where}
}

// BuildWhere - build the []Attr based on the URL query parameters
func BuildWhere(values map[string][]string) []Attr {
	if len(values) == 0 {
		return nil
	}
	var where []Attr
	for k, v := range values {
		where = append(where, Attr{Key: k, Val: v[0]})
	}
	return where
}
