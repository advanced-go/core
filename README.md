# core

Core provides functionaliy for common development tasks such as: error handling, HTTP exchange functionality, HTTP handler testing, and access logging. 
The packages and specific implementations are as follows:  

## runtime
[Runtime][runtimepkg] implements environment, request context, status, and error types. The status type is used as a function return value, and provides additional context for an error, such as location, request id, http content and status codes. The error handler types are designed to be used as generic parameters.
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

// ErrorHandler - generic parameter error handler interface
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
[Http2][http2pkg] provides functionality for processing HTTP request/response exchange. 

The Do() function supports reading a response from disk or via a proxy.
~~~
// Do - do a Http exchange with a runtime.Status
func Do(req *http.Request) (resp *http.Response, status runtime.Status)
    // implementation details
}
~~~

Generic deserialization is supported:

~~~
// Deserialize - provide deserialization of a request/response body
func Deserialize[T any](body io.ReadCloser) (T, runtime.Status) {
    // implementation details
}
~~~

Also included is a common HTTP write response function:

~~~
// WriteResponse - write a http.Response, utilizing the content, status, and headers
// Only supports []byte, string, io.Reader, and io.ReaderCloser for content
func WriteResponse[E runtime.ErrorHandler](w http.ResponseWriter, content any, status runtime.Status, headers any) {
    // implementation details
}
~~~

## http2test
[Http2test][http2pkg] provides functionality for testing HTTP calls. 

Reading a response and a request will deserialize the type including the content.
~~~
func ReadResponse(uri *url.URL) (*http.Response, error) {
  // implementation details
}
func ReadRequest(uri *url.URL) (*http.Request, error) {
  // implementation details
}
~~~

Several functions are provided to facilitate automated testing when a Test_httpHandler(t *testing.T) is used. 

## access
[Access][accesspkg] implements access logging via am application configurable access log handler. Logging of internal traffic is supported and can be configured at runtime.

Direct and deferred logging are supported. The threshold attributes provide indication of requests that do not meet an SLO.
~~~
// Log - access logging
func Log(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, thresholdFlags string) {
 // implementation details
}

// LogDeferred - deferred accessing logging
func LogDeferred(traffic string, req *http.Request, threshold int, thresholdFlags string, statusCode func() int) func() {
 // implementation details
}

// Defered internal traffic logging example for an HTTP handler.
func() (status runtime.Status) {
    defer access.LogDeferred(access.InternalTraffic, r, -1, "", access.NewStatusCodeClosure(&status))()
    return httpHandler[runtime.Log](nil, w, r)
}()

~~~


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

