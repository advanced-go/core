package runtime

import (
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"net/http"
	"time"
)

const (
	ContentTypeText = "text/plain" // charset=utf-8
	ContentTypeJson = "application/json"
	ContentType     = "Content-Type"
)

// https://grpc.github.io/grpc/core/md_doc_statuscodes.html

const (
	NilDuration = time.Duration(-1)
)

const (
	StatusInvalidContent     = codes.Code(90)           // Content is not available, is nil, or is of the wrong type, usually found via unmarshalling
	StatusIOError            = codes.Code(91)           // I/O operation failed
	StatusJsonDecodeError    = codes.Code(92)           // Json decoding failed
	StatusJsonEncodeError    = codes.Code(93)           // Json decoding failed
	StatusNotProvided        = codes.Code(94)           // No status is available
	StatusRateLimited        = codes.Code(95)           // Rate limited
	StatusOK                 = codes.OK                 // Not an error; returned on success.
	StatusCancelled          = codes.Canceled           // The operation was cancelled, typically by the caller.
	StatusUnknown            = codes.Unknown            // Unknown error. For example, this error may be returned when a Status value received from another address space belongs to an error space that is not known in this address space. Also errors raised by APIs that do not return enough error information may be converted to this error.
	StatusInvalidArgument    = codes.InvalidArgument    // The client specified an invalid argument. Note that this differs from FAILED_PRECONDITION. INVALID_ARGUMENT indicates arguments that are problematic regardless of the state of the system (e.g., a malformed file name).
	StatusDeadlineExceeded   = codes.DeadlineExceeded   // The deadline expired before the operation could complete. For operations that change the state of the system, this error may be returned even if the operation has completed successfully. For example, a successful response from a server could have been delayed long
	StatusNotFound           = codes.NotFound           // Some requested entity (e.g., file or directory) was not found. Note to server developers: if a request is denied for an entire class of users, such as gradual feature rollout or undocumented allowlist, NOT_FOUND may be used. If a request is denied for some users within a class of users, such as user-based access control, PERMISSION_DENIED must be used.
	StatusAlreadyExists      = codes.AlreadyExists      // The entity that a client attempted to create (e.g., file or directory) already exists.
	StatusPermissionDenied   = codes.PermissionDenied   // The caller does not have permission to execute the specified operation. PERMISSION_DENIED must not be used for rejections caused by exhausting some host (use RESOURCE_EXHAUSTED instead for those errors). PERMISSION_DENIED must not be used if the caller can not be identified (use UNAUTHENTICATED instead for those errors). This error code does not imply the request is valid or the requested entity exists or satisfies other pre-conditions.
	StatusResourceExhausted  = codes.ResourceExhausted  // Some host has been exhausted, perhaps a per-user quota, or perhaps the entire file system is out of space.
	StatusFailedPrecondition = codes.FailedPrecondition // The operation was rejected because the system is not in a state required for the operation's execution. For example, the directory to be deleted is non-empty, an rmdir operation is applied to a non-directory, etc. Service implementors can use the following guidelines to decide between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE: (a) Use UNAVAILABLE if the client can retry just the failing call. (b) Use ABORTED if the client should retry at a higher level (e.g., when a client-specified test-and-set fails, indicating the client should restart a read-modify-write sequence). (c) Use FAILED_PRECONDITION if the client should not retry until the system state has been explicitly fixed. E.g., if an "rmdir" fails because the directory is non-empty, FAILED_PRECONDITION should be returned since the client should not retry unless the files are deleted from the directory.
	StatusAborted            = codes.Aborted            // The operation was aborted, typically due to a concurrency issue such as a sequencer check failure or transaction abort. See the guidelines above for deciding between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE.
	StatusOutOfRange         = codes.OutOfRange         // The operation was attempted past the valid range. E.g., seeking or reading past end-of-file. Unlike INVALID_ARGUMENT, this error indicates a problem that may be fixed if the system state changes. For example, a 32-bit file system will generate INVALID_ARGUMENT if asked to read at an offset that is not in the range [0,2^32-1], but it will generate OUT_OF_RANGE if asked to read from an offset past the current file size. There is a fair bit of overlap between FAILED_PRECONDITION and OUT_OF_RANGE. We recommend using OUT_OF_RANGE (the more specific error) when it applies so that callers who are iterating through a space can easily look for an OUT_OF_RANGE error to detect when they are done.
	StatusUnimplemented      = codes.Unimplemented      // The operation is not implemented or is not supported/enabled in this service.
	StatusInternal           = codes.Internal           // Internal errors. This means that some invariants expected by the underlying system have been broken. This error code is reserved for serious errors.
	StatusUnavailable        = codes.Unavailable        // The service is currently unavailable. This is most likely a transient condition, which can be corrected by retrying with a backoff. Note that it is not always safe to retry non-idempotent operations.
	StatusDataLoss           = codes.DataLoss           // Unrecoverable data loss or corruption.
	StatusUnauthenticated    = codes.Unauthenticated    // The request does not have valid authentication credentials for the operation.
	_maxGRPCCode             = StatusUnauthenticated
)

// IsErrors - determine if there are errors in an []error
func IsErrors(errs []error) bool {
	return !(len(errs) == 0 || (len(errs) == 1 && errs[0] == nil))
}

// Status - struct for status data
type Status struct {
	code      codes.Code
	duration  time.Duration
	location  string
	requestId string
	errs      []error
	handled   bool
	content   []byte
	//md        metadata.MD
}

// NewStatus - new Status from a code
func NewStatus(code codes.Code) *Status {
	s := new(Status)
	s.code = code
	s.duration = NilDuration
	return s
}

// NewHttpStatus - new Status from a http status code
func NewHttpStatus(code int) *Status {
	return NewStatus(codes.Code(code))
}

// NewStatusOK - new OK status
func NewStatusOK() *Status {
	return NewStatus(StatusOK)
}

// NewStatusError - new Status from a code, location, and optional errors
func NewStatusError(code codes.Code, location string, errs ...error) *Status {
	s := NewStatus(code)
	s.location = location
	if !IsErrors(errs) {
		s.code = StatusOK
	} else {
		if code == 0 {
			s.code = StatusInternal
		}
		s.addErrors(errs...)
	}
	return s
}

/*


func NewStatusWithContext(code codes.Code, location string, ctx context.Context, errs ...error) *Status {
	s := NewStatus(code, location, errs...)
	//s.SetContext(ctx)
	return s
}

// NewHttpStatus - new Status from a http.Response, location, and optional errors
// func NewHttpStatus(resp *http.Response, location string, errs ...error) *Status {
func NewHttpStatus(resp *http.Response, errs ...error) *Status {
	var code codes.Code
	if resp == nil {
		code = StatusInvalidContent
	} else {
		code = codes.Code(resp.StatusCode)
	}
	s := NewStatusCode(code)
	if IsErrors(errs) {
		s.addErrors(errs...)
		s.code = http.StatusInternalServerError
	}
	return s
}

func NewStatusInvalidArgument(location string, err error) *Status {
	return NewStatus(StatusInvalidArgument, err)
}

// NewStatusError - new Internal status with location and errors
// func NewStatusError(location string, errs ...error) *Status {
func NewStatusError(errs ...error) *Status {
	return NewStatus(StatusInternal, errs...)
}

*/

// IsGRPCCode - gRPC code functions
func (s *Status) IsGRPCCode() bool { return s.code >= codes.OK && s.code <= _maxGRPCCode }
func (s *Status) Code() codes.Code { return s.code }
func (s *Status) SetCode(code codes.Code) *Status {
	s.code = code
	return s
}

func (s *Status) String() string {
	if s.IsGRPCCode() {
		if s.IsErrors() {
			if s.location == "" {
				return fmt.Sprintf("%v %v", s.code, s.errs)
			} else {
				return fmt.Sprintf("%v %v %v", s.code, s.location, s.errs)
			}
		} else {
			return fmt.Sprintf("%v", s.code)
		}
	} else {
		if s.IsErrors() {
			if s.Location() == "" {
				return fmt.Sprintf("%v %v", s.Description(), s.errs)
			} else {
				return fmt.Sprintf("%v %v %v", s.Description(), s.location, s.errs)
			}
		} else {
			return fmt.Sprintf("%v", s.Description())
		}
	}
}

// ErrorsHandled - determine errors status
func (s *Status) ErrorsHandled() bool { return s.handled }
func (s *Status) SetErrorsHandled()   { s.handled = true }
func (s *Status) IsErrors() bool      { return s.errs != nil && len(s.errs) > 0 }
func (s *Status) Errors() []error     { return s.errs }
func (s *Status) FirstError() error {
	if s.IsErrors() {
		return s.errs[0]
	}
	return nil
}
func (s *Status) addErrors(errs ...error) *Status {
	for _, e := range errs {
		if e == nil {
			continue
		}
		s.errs = append(s.errs, e)
	}
	return s
}

//func (s *Status) RemoveErrors() { s.errs = nil }

// Duration - get duration
func (s *Status) Duration() time.Duration { return s.duration }
func (s *Status) SetDuration(duration time.Duration) *Status {
	s.duration = duration
	return s
}

// Location - location
func (s *Status) Location() string { return s.location }
func (s *Status) SetLocation(location string) *Status {
	s.location = location
	return s
}

// RequestId  - request id
func (s *Status) RequestId() string { return s.requestId }
func (s *Status) SetRequestId(requestId any) *Status {
	if str, ok := requestId.(string); ok {
		s.requestId = str
	} else {
		if st, ok1 := requestId.(*Status); ok1 && st != nil {
			s.requestId = st.RequestId()
		}
	}
	return s
}

/*
func (s *Status) SetLocationAndId(location string, requestId any) *Status {
	s.SetLocation(location)
	s.SetRequestId(requestId)
	return s
}


*/

// IsContent - content
func (s *Status) IsContent() bool { return s.content != nil }
func (s *Status) Content() []byte { return s.content }
func (s *Status) RemoveContent() {
	s.content = nil
}
func (s *Status) SetContent(content any) *Status {
	if content == nil {
		return s
	}

	switch data := content.(type) {
	case string:
		buf := []byte(data)
		s.content = buf
	case []byte:
		s.content = data
	case error:
		s.content = []byte(data.Error())
	default:
		buf, err := json.Marshal(data)
		if err != nil {
			s.content = []byte("invalid non Json serializable type")
		} else {
			s.content = buf
		}
	}
	return s
}

func (s *Status) OK() bool              { return s.code == StatusOK || s.code == http.StatusOK }
func (s *Status) InvalidArgument() bool { return s.code == StatusInvalidArgument }
func (s *Status) Unauthenticated() bool {
	return s.code == StatusUnauthenticated || s.code == http.StatusUnauthorized
}
func (s *Status) PermissionDenied() bool {
	return s.code == StatusPermissionDenied || s.code == http.StatusForbidden
}
func (s *Status) NotFound() bool { return s.code == StatusNotFound || s.code == http.StatusNotFound }
func (s *Status) Internal() bool {
	return s.code == StatusInternal || s.code == http.StatusInternalServerError
}
func (s *Status) Timeout() bool {
	return s.code == StatusDeadlineExceeded || s.code == http.StatusGatewayTimeout
}
func (s *Status) ServiceUnavailable() bool {
	return s.code == StatusUnavailable || s.code == http.StatusServiceUnavailable
}

func (s *Status) Http() int {
	if s.code >= http.StatusContinue {
		return int(s.code)
	}
	//return http.StatusInternalServerError

	var code = http.StatusInternalServerError
	switch s.code {
	case StatusOK:
		code = http.StatusOK
	case StatusInvalidArgument:
		code = http.StatusBadRequest
	case StatusUnauthenticated:
		code = http.StatusUnauthorized
	case StatusPermissionDenied:
		code = http.StatusForbidden
	case StatusNotFound:
		code = http.StatusNotFound
	case StatusRateLimited:
		code = http.StatusTooManyRequests
	case StatusInternal:
		code = http.StatusInternalServerError
	case StatusUnavailable:
		code = http.StatusServiceUnavailable
	case StatusDeadlineExceeded:
		code = http.StatusGatewayTimeout
	case StatusInvalidContent:
		code = http.StatusNoContent
	case StatusCancelled,
		StatusUnknown,
		StatusAlreadyExists,
		StatusResourceExhausted,
		StatusFailedPrecondition,
		StatusAborted,
		StatusOutOfRange,
		StatusUnimplemented,
		StatusDataLoss:
	}
	return code
}

func (s *Status) Description() string {
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
		return "Not provided"
	case StatusRateLimited:
		return "Rate limited"

	case StatusOK, http.StatusOK:
		return "OK"
	case StatusInvalidArgument, http.StatusBadRequest:
		return "Bad Request"
	case StatusDeadlineExceeded, http.StatusGatewayTimeout:
		return "Timeout"
	case StatusNotFound, http.StatusNotFound:
		return "Not Found"
	case StatusPermissionDenied, http.StatusForbidden:
		return "Permission Denied"
	case StatusInternal, http.StatusInternalServerError:
		return "Internal Error"
	case StatusUnavailable, http.StatusServiceUnavailable:
		return "Service Unavailable"
	case StatusUnauthenticated, http.StatusUnauthorized:
		return "Unauthorized"

	// Unmapped
	case StatusCancelled:
		return "The operation was cancelled, typically by the caller"
	case StatusUnknown:
		return "Unknown error" // For example, this error may be returned when a Status value received from another address space belongs to an error space that is not known in this address space. Also errors raised by APIs that do not return enough error information may be converted to this error."
	case StatusAlreadyExists:
		return "The entity that a client attempted to create already exists"
	case StatusResourceExhausted:
		return "Some host has been exhausted" //perhaps a per-user quota, or perhaps the entire file system is out of space."
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
	}
	return fmt.Sprintf("code not mapped: %v", s.code)
}
