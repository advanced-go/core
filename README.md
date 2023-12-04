# core

Core was inspired to capitalize on the Go language for application development. Determining the patterns that need to be employed is critical for writing clear idiomatic Go code. This YouTube video [Edward Muller - Go Anti-Patterns][emuller], does an excellent job of framing idiomatic go. 
[Robert Griesemer - The Evolution of Go][rgriesemer], @ 4:00 minutes, also presents an important analogy between Go modules, packages, closures, and LEGO® bricks. Reviewing the Go standard library packaging structure provides a blueprint for an application architecture, and underscores how essential package design is for idiomatic Go. 

What follows is a description of the packages in Core, highlighting specific patterns and template implementations.  

## runtime
[Runtime][runtimepkg] implements environment, request context, status, and error types. The status type is used as a function return value, and provides error, location, http content and status codes. 

The error handler types are designed to be used as template parameters.
~~~
// Status - used to add additional context to error handling
type Status interface {
    Code() int
    OK() bool
    Http() int

    IsErrors() bool
    Errors() []error
    FirstError() error

    Duration() time.Duration
    SetDuration(duration time.Duration) Status

    RequestId() string
    SetRequestId(requestId any) Status

    Location() []string
    AddLocation(location string) Status

    IsContent() bool
    Content() any
    ContentHeader() http.Header
    ContentString() string
    SetContent(content any, jsonContent bool) Status

    Description() string
    String() string
}

// ErrorHandler - template parameter error handler interface
type ErrorHandler interface {
    Handle(s Status, requestId string, callerLocation string) Status
}
~~~

Context functionality is provied for a request Id, and a ProxyContext used for testing:

~~~
// ContextWithRequestId - creates a new Context with a request id
func ContextWithRequestId(ctx context.Context, requestId string) context.Context {
    // implementation details
}

// ContextWithProxy - create a new Context interface, containing a proxy
func ContextWithProxy(ctx context.Context, proxy any) context.Context {
    // implementation details
}
~~~

## http2
[Http2][http2pkg] provides functionality for processing an Http request/response. Exchange functionality is provied via a templated function, utilizing
template paramters for error processing, deserialization type, and the function for processing the http.Client.Do():

~~~
func Do[E runtime.ErrorHandler, H Exchange, T any](req *http.Request) (resp *http.Response, t T, status *runtime.Status) {
    // implementation details
}
~~~

The deserialization function is also templated:

~~~
// Deserialize - templated function, providing deserialization of a request/response body
func Deserialize[E runtime.ErrorHandler, T any](body io.ReadCloser) (T, *runtime.Status) {
    // implementation details
}
~~~

Testing Http calls is implemented through a proxy design pattern: a context.Context interface that contains an http.Client.Do() call.

~~~
// HttpExchange - interface for Http request/response interaction
type HttpExchange interface {
	Do(req *http.Request) (*http.Response, error)
}
~~~

Exchange also includes a common http write response function:

~~~
// WriteResponse - write a http.Response, utilizing the data, status, and headers for controlling the content
func WriteResponse(w http.ResponseWriter, buf []byte, status *runtime.Status, headers ...string) {
    // implementation details
}
~~~

## access
[Access][accesspkg] implements environment, request context, status, error, and output types. The status type is used extensively as a function return value, and provides error, http, and gRPC status codes. 

## handler
[Handler][handlerpkg] implements environment, request context, status, error, and output types. The status type is used extensively as a function return value, and provides error, http, and gRPC status codes. 

## io2
[Io2][io2pkg] implements environment, request context, status, error, and output types. The status type is used extensively as a function return value, and provides error, http, and gRPC status codes. 

## json2
[Json2][json2pkg] implements environment, request context, status, error, and output types. The status type is used extensively as a function return value, and provides error, http, and gRPC status codes. 

## strings
[Strings][stringspkg] implements environment, request context, status, error, and output types. The status type is used extensively as a function return value, and provides error, http, and gRPC status codes. 

## resiliency
[Resiliency][resiliencypkg] implements environment, request context, status, error, and output types. The status type is used extensively as a function return value, and provides error, http, and gRPC status codes. 




[runtimepkg]: <https://pkg.go.dev/github.com/advanced-go/core/runtime>
[http2pkg]: <https://pkg.go.dev/github.com/advanced-go/core/http2>
[accesspkg]: <https://pkg.go.dev/github.com/advanced-go/core/access>
[handlerpkg]: <https://pkg.go.dev/github.com/advanced-go/core/handler>
[io2pkg]: <https://pkg.go.dev/github.com/advanced-go/core/io2>
[json2pkg]: <https://pkg.go.dev/github.com/advanced-go/core/json2>
[stringspkg]: <https://pkg.go.dev/github.com/advanced-go/core/strings>
[resiliencypkg]: <https://pkg.go.dev/github.com/advanced-go/core/resiliency][=tghtvfcx>

