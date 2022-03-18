package errors

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	// InternalServerErrorCode internal server error code
	InternalServerErrorCode = 500
)

func (e *Error) Error() string {
	return fmt.Sprintf("error: code: %d; message= %s; metadata = %v", e.Code, e.Message, e.Metadata)
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := proto.Clone(e).(*Error)
	err.Metadata = md
	return err
}

// New returns an error object for the code, message.
func New(code int, message string) *Error {
	return &Error{
		Code:    int32(code),
		Message: message,
	}
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(code int, format string, a ...interface{}) *Error {

	return New(code, fmt.Sprintf(format, a...))
}

// GRPCStatus returns the Status, custom error
func (e *Error) GRPCStatus() *status.Status {
	s, _ := status.New(codes.Code(e.Code), e.Message).WithDetails(&errdetails.ErrorInfo{
		Metadata: e.Metadata,
	})
	return s
}

// BadRequest bad request
func BadRequest(message string) *Error {
	return Newf(400001, message)
}

// UnpaidOrderExists ...
func UnpaidOrderExists(message string) *Error {
	return Newf(400002, message)
}

// OrderNotFound ...
func OrderNotFound(message string) *Error {
	return Newf(400003, message)
}

// InvoiceNotFound ...
func InvoiceNotFound(message string) *Error {
	return Newf(400004, message)
}

// InternalServerError ...
func InternalServerError(message string) *Error {
	return Newf(InternalServerErrorCode, message)
}

// MissingMetadataError ...
func MissingMetadataError(message string) *Error {
	return Newf(400005, message)
}

// PreconditionFailed ...
func PreconditionFailed(message string) *Error {
	// return status.Errorf(400006, , message)
	return Newf(400006, message)
}

// Canceled client cancel request
func Canceled(message string) *Error {
	return Newf(int(codes.Canceled), message)
}

// DeadlineExceeded timeout
func DeadlineExceeded(message string) *Error {
	return Newf(int(codes.DeadlineExceeded), message)
}
