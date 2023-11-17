package runtime

import (
	"fmt"
	"net/http"
	"time"
)

// https://grpc.github.io/grpc/core/md_doc_statuscodes.html

const (
	NilDuration     = time.Duration(-1)
	contentType     = "Content-type"
	contentTypeJson = "application/json"
	contentLocation = "Content-Location"
)

const (
	StatusInvalidContent  = int(90) // Content is not available, is nil, or is of the wrong type, usually found via unmarshalling
	StatusIOError         = int(91) // I/O operation failed
	StatusJsonDecodeError = int(92) // Json decoding failed
	StatusJsonEncodeError = int(93) // Json decoding failed
	StatusNotProvided     = int(94) // No status is available
	StatusRateLimited     = int(95) // Rate limited
	StatusNotStarted      = int(96) // Not started
	StatusHaveContent     = int(97) // Content is available

	/*
		StatusOK                 = codes.OK                 // Not an error; returned on success.
		StatusCancelled          = codes.Canceled           // The operation was cancelled, typically by the caller.
		StatusUnknown            = codes.Unknown            // Unknown error. For example, this error may be returned when a Status value received from another address space belongs to an error space that is not known in this address space. Also errors raised by APIs that do not return enough error information may be converted to this error.
	*/

	StatusInvalidArgument  = 3 //codes.InvalidArgument    // The client specified an invalid argument. Note that this differs from FAILED_PRECONDITION. INVALID_ARGUMENT indicates arguments that are problematic regardless of the state of the system (e.g., a malformed file name).
	StatusDeadlineExceeded = 4 //codes.DeadlineExceeded   // The deadline expired before the operation could complete. For operations that change the state of the system, this error may be returned even if the operation has completed successfully. For example, a successful response from a server could have been delayed long

	/*	StatusNotFound           = codes.NotFound           // Some requested entity (e.g., file or directory) was not found. Note to server developers: if a request is denied for an entire class of users, such as gradual feature rollout or undocumented allowlist, NOT_FOUND may be used. If a request is denied for some users within a class of users, such as user-based access control, PERMISSION_DENIED must be used.
		StatusAlreadyExists      = codes.AlreadyExists      // The entity that a client attempted to create (e.g., file or directory) already exists.
		StatusPermissionDenied   = codes.PermissionDenied   // The caller does not have permission to execute the specified operation. PERMISSION_DENIED must not be used for rejections caused by exhausting some startup (use RESOURCE_EXHAUSTED instead for those errors). PERMISSION_DENIED must not be used if the caller can not be identified (use UNAUTHENTICATED instead for those errors). This error code does not imply the request is valid or the requested entity exists or satisfies other pre-conditions.
		StatusResourceExhausted  = codes.ResourceExhausted  // Some startup has been exhausted, perhaps a per-user quota, or perhaps the entire file system is out of space.
		StatusFailedPrecondition = codes.FailedPrecondition // The operation was rejected because the system is not in a state required for the operation's execution. For example, the directory to be deleted is non-empty, an rmdir operation is applied to a non-directory, etc. Service implementors can use the following guidelines to decide between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE: (a) Use UNAVAILABLE if the client can retry just the failing call. (b) Use ABORTED if the client should retry at a higher level (e.g., when a client-specified test-and-set fails, indicating the client should restart a read-modify-write sequence). (c) Use FAILED_PRECONDITION if the client should not retry until the system state has been explicitly fixed. E.g., if an "rmdir" fails because the directory is non-empty, FAILED_PRECONDITION should be returned since the client should not retry unless the files are deleted from the directory.
		StatusAborted            = codes.Aborted            // The operation was aborted, typically due to a concurrency issue such as a sequencer check failure or transaction abort. See the guidelines above for deciding between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE.
		StatusOutOfRange         = codes.OutOfRange         // The operation was attempted past the valid range. E.g., seeking or reading past end-of-file. Unlike INVALID_ARGUMENT, this error indicates a problem that may be fixed if the system state changes. For example, a 32-bit file system will generate INVALID_ARGUMENT if asked to read at an offset that is not in the range [0,2^32-1], but it will generate OUT_OF_RANGE if asked to read from an offset past the current file size. There is a fair bit of overlap between FAILED_PRECONDITION and OUT_OF_RANGE. We recommend using OUT_OF_RANGE (the more specific error) when it applies so that callers who are iterating through a space can easily look for an OUT_OF_RANGE error to detect when they are done.
		StatusUnimplemented      = codes.Unimplemented      // The operation is not implemented or is not supported/enabled in this service.
		StatusInternal           = codes.Internal           // Internal errors. This means that some invariants expected by the underlying system have been broken. This error code is reserved for serious errors.
	*/
	//	StatusUnavailable        = 14 //codes.Unavailable        // The service is currently unavailable. This is most likely a transient condition, which can be corrected by retrying with a backoff. Note that it is not always safe to retry non-idempotent operations.
	/*
		StatusDataLoss           = codes.DataLoss           // Unrecoverable data loss or corruption.
		StatusUnauthenticated    = codes.Unauthenticated    // The request does not have valid authentication credentials for the operation.
		_maxGRPCCode             = StatusUnauthenticated
	*/
)

// IsErrors - determine if there are errors in an []error
func IsErrors(errs []error) bool {
	return !(len(errs) == 0 || (len(errs) == 1 && errs[0] == nil))
}

type Status interface {
	Code() int
	SetCode(code int) Status

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
	ContentString() string
	SetContent(content any, jsonContent bool) Status

	Header() http.Header

	/*
		SetContentType(str string) Status
		SetContentLocation(location string) Status
		SetContentTypeAndLocation(location string) Status

		SetHeader(header http.Header, keys ...string) Status
		CopyHeader(header http.Header) Status


	*/

	OK() bool
	Http() int
	Description() string
	String() string
}

// Status - struct for status data
type status struct {
	code      int //type codes.Code uint32
	duration  time.Duration
	handled   bool
	requestId string
	location  []string
	errs      []error
	content   any
	header    http.Header
}

// NewStatus - new Status from a code
func NewStatus(code int) Status {
	return newStatus(code)
}

func newStatus(code int) *status {
	s := new(status)
	s.code = code
	s.duration = NilDuration
	return s
}

// NewStatusOK - new OK status
func NewStatusOK() Status {
	return newStatus(http.StatusOK)
}

// NewStatusError - new Status from a code, location, and optional errors
func NewStatusError(code int, location string, errs ...error) Status {
	s := newStatus(code)
	s.location = append(s.location, location)
	if !IsErrors(errs) {
		s.code = http.StatusOK
	} else {
		if code == 0 {
			s.code = http.StatusInternalServerError
		}
		s.addErrors(errs...)
	}
	return s
}

// Code - functions
func (s *status) Code() int { return s.code }
func (s *status) SetCode(code int) Status {
	s.code = code
	return s
}

func (s *status) String() string {
	if s.IsErrors() {
		return fmt.Sprintf("%v %v", s.Description(), s.errs)
	} else {
		return fmt.Sprintf("%v", s.Description())
	}
}

// IsErrors - determine errors status
func (s *status) IsErrors() bool  { return s.errs != nil && len(s.errs) > 0 }
func (s *status) Errors() []error { return s.errs }
func (s *status) FirstError() error {
	if s.IsErrors() {
		return s.errs[0]
	}
	return nil
}

func errorsHandled(s Status) bool {
	if st, ok := any(s).(*status); ok {
		return st.handled
	}
	return false
}

func setErrorsHandled(s Status) {
	if st, ok := any(s).(*status); ok {
		st.handled = true
	}
}

func (s *status) addErrors(errs ...error) {
	for _, e := range errs {
		if e == nil {
			continue
		}
		s.errs = append(s.errs, e)
	}
	//return s
}

// Duration - get duration
func (s *status) Duration() time.Duration { return s.duration }
func (s *status) SetDuration(duration time.Duration) Status {
	s.duration = duration
	return s
}

// RequestId  - request id
func (s *status) RequestId() string { return s.requestId }
func (s *status) SetRequestId(requestId any) Status {
	if len(s.requestId) != 0 {
		return s
	}
	id := RequestId(requestId)
	if len(id) > 0 {
		s.requestId = id
	}
	return s
}

// Location - location
func (s *status) Location() []string { return s.location }
func (s *status) AddLocation(location string) Status {
	if len(location) >= 0 {
		s.location = append(s.location, location)
	}
	return s
}

// IsContent - content
func (s *status) IsContent() bool { return s.content != nil }
func (s *status) Content() any    { return s.content }
func (s *status) ContentString() string {
	switch ptr := s.content.(type) {
	case string:
		return ptr
	case []byte:
		return string(ptr)
	}
	return ""
}

//func (s *status) RemoveContent() {
//	s.content = nil
//}

func (s *status) SetContent(content any, jsonContent bool) Status {
	if content == nil {
		return s
	}
	s.content = content
	if jsonContent {
		s.SetContentType(contentTypeJson)
	}
	return s
}

/*
	func (s *status) SetJsonContent(content any) *Status {
		if content == nil {
			return s
		}
		s.content = content
		s.SetContentType(ContentTypeJson)
		return s
	}
*/
func (s *status) SetContentType(str string) Status {
	if len(str) == 0 {
		return s
	}
	s.Header().Set(contentType, str)
	return s
}
func (s *status) SetContentLocation(location string) Status {
	if len(location) == 0 {
		return s
	}
	s.Header().Set(contentLocation, location)
	return s
}
func (s *status) SetContentTypeAndLocation(location string) Status {
	if len(location) == 0 {
		return s
	}
	s.Header().Set(contentType, location)
	s.Header().Set(contentLocation, location)
	return s
}

// Header - header map
func (s *status) Header() http.Header {
	if s.header == nil {
		s.header = make(http.Header)
	}
	return s.header
}
func (s *status) SetHeader(header http.Header, keys ...string) Status {
	if header == nil {
		return s
	}
	for _, key := range keys {
		s.header.Set(key, header.Get(key))
	}
	return s
}
func (s *status) CopyHeader(header http.Header) Status {
	if header == nil {
		return s
	}
	for k, _ := range s.header {
		header.Set(k, s.header.Get(k))
	}
	return s
}

func (s *status) OK() bool { return s.code == http.StatusOK }

//func (s *status) NotFound() bool { return s.code == http.StatusNotFound }

func (s *status) Http() int {
	// Catch all valid http status codes
	if s.code >= http.StatusContinue {
		return s.code
	}
	// map known
	switch s.code {
	case StatusInvalidArgument:
		return http.StatusInternalServerError
	case StatusDeadlineExceeded:
		return http.StatusGatewayTimeout
	}
	// all others
	return http.StatusInternalServerError
}

func (s *status) Description() string {
	switch s.code {
	// Mapped
	case StatusInvalidContent:
		return "Invalid Content"
	case StatusIOError:
		return "I/O Failure"
	case StatusJsonEncodeError:
		return "Json Encode Failure"
	case StatusJsonDecodeError:
		return "Json Decode Failure"
	case StatusNotProvided:
		return "Not Provided"
	case StatusRateLimited:
		return "Rate Limited"
	case StatusNotStarted:
		return "Not Started"
	case StatusDeadlineExceeded:
		return "Deadline Exceeded"
	case StatusInvalidArgument:
		return "Invalid Argument"
	case StatusHaveContent:
		return "Content Available"

	//case StatusUnavailable:
	//	return "Invalid Argument"

	//Http
	case http.StatusOK:
		return "OK"
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusGatewayTimeout:
		return "Timeout"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusMethodNotAllowed:
		return "Method Not Allowed"
	case http.StatusForbidden:
		return "Permission Denied"
	case http.StatusInternalServerError:
		return "Internal Error"
	case http.StatusServiceUnavailable:
		return "Service Unavailable"
	case http.StatusUnauthorized:
		return "Unauthorized"

		// Unmapped
		/*
			case StatusCancelled:
				return "The operation was cancelled, typically by the caller"
			case StatusUnknown:
				return "Unknown error" // For example, this error may be returned when a Status value received from another address space belongs to an error space that is not known in this address space. Also errors raised by APIs that do not return enough error information may be converted to this error."
			case StatusAlreadyExists:
				return "The entity that a client attempted to create already exists"
			case StatusResourceExhausted:
				return "Some startup has been exhausted" //perhaps a per-user quota, or perhaps the entire file system is out of space."
			case StatusFailedPrecondition:
				return "The operation was rejected because the system is not in a state required for the operation's execution" //For example, the directory to be deleted is non-empty, an rmdir operation is applied to a non-directory, etc. Service implementors can use the following guidelines to decide between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE: (a) Use UNAVAILABLE if the client can retry just the failing call. (b) Use ABORTED if the client should retry at a higher level (e.g., when a client-specified test-and-set fails, indicating the client should restart a read-modify-write sequence). (c) Use FAILED_PRECONDITION if the client should not retry until the system state has been explicitly fixed. E.g., if an "rmdir" fails because the directory is non-empty, FAILED_PRECONDITION should be returned since the client should not retry unless the files are deleted from the directory.
			case StatusAborted:
				return "The operation was aborted" // typically due to a concurrency issue such as a sequencer check failure or transaction abort. See the guidelines above for deciding between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE."
			case StatusOutOfRange:
				return "The operation was attempted past the valid range" // E.g., seeking or reading past end-of-file. Unlike INVALID_ARGUMENT, this error indicates a problem that may be fixed if the system state changes. For example, a 32-bit file system will generate INVALID_ARGUMENT if asked to read at an offset that is not in the range [0,2^32-1], but it will generate OUT_OF_RANGE if asked to read from an offset past the current file size. There is a fair bit of overlap between FAILED_PRECONDITION and OUT_OF_RANGE. We recommend using OUT_OF_RANGE (the more specific error) when it applies so that callers who are iterating through a space can easily look for an OUT_OF_RANGE error to detect when they are done."
			case StatusUnimplemented:
				return "The operation is not implemented or is not supported/enabled in this service"
			case StatusDataLoss:
				return "Unrecoverable data loss or corruption"
		*/
	}
	return fmt.Sprintf("error: code not mapped: %v", s.code)
}
