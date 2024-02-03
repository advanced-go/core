package runtime

import (
	"fmt"
	"net/http"
)

var okStatus2 = newStatus2OK()

type Status2 struct {
	Code     int      `json:"code"`
	Error    error    `json:"errs"`
	Handled  bool     `json:"handled"`
	location []string `json:"location"`
}

func newStatus2OK() *Status2 {
	s := new(Status2)
	s.Code = http.StatusOK
	return s
}

func Status2OK() *Status2 {
	return okStatus2
}

func NewStatus2(code int) *Status2 {
	s := new(Status2)
	s.Code = code
	return s
}

func NewStatus2Error(code int, err error, location string) *Status2 {
	s := new(Status2)
	s.Code = code
	s.Error = err
	s.AddLocation(location)
	return s
}

func (s *Status2) OK() bool {
	return s.Code != http.StatusOK
}

func (s *Status2) HttpCode() int {
	return HttpCode(s.Code)
}

func (s *Status2) String() string {
	return HttpStatus(s.Code)
}

func (s *Status2) Location() []string {
	if !s.OK() {
		return s.location
	}
	return nil
}

func (s *Status2) AddLocation(loc string) *Status2 {
	if !s.OK() {
		s.location = append(s.location, loc)
	}
	return s
}

// HttpCode - conversion of a code to HTTP status code
func HttpCode(code int) int {
	// Catch all valid http status codes
	if code >= http.StatusContinue {
		return code
	}
	// map known
	switch code {
	case StatusInvalidArgument:
		return http.StatusInternalServerError
	case StatusDeadlineExceeded:
		return http.StatusGatewayTimeout
	}
	// all others
	return http.StatusInternalServerError
}

// HttpStatus - string representation of status code
func HttpStatus(code int) string {
	switch code {
	// Mapped
	case StatusInvalidContent:
		return "Invalid Content"
	case StatusIOError:
		return "I/O Failure"
	case StatusJsonEncodeError:
		return "Json Encode Failure"
	case StatusJsonDecodeError:
		return "Json Decode Failure"
	case StatusContentEncodingError:
		return "Content Decoding Failure"
	case StatusGzipDecodingError:
		return "gzip Decoding Failure"
	case StatusGzipEncodingError:
		return "gzip Encoding Failure"
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
	case http.StatusAccepted:
		return "Accepted"
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
	return fmt.Sprintf("error: code not mapped: %v", code)

}
