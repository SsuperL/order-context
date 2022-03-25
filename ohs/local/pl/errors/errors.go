package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// InternalServerErrorCode internal server error code
	InternalServerErrorCode = 500
)

// BadRequest bad request
func BadRequest(message string, details ...interface{}) error {
	return status.Errorf(codes.Code(400001), message)
}

// UnpaidOrderExists ...
func UnpaidOrderExists(message string, details ...interface{}) error {
	// st:=status.Newf(codes.Code(400002), message)
	// ds,err :=st.WithDetails(
	// )
	// if err != nil {
	// 	return ds.Err()
	// }
	return status.Errorf(codes.Code(400002), message)
}

// OrderNotFound ...
func OrderNotFound(message string, details ...interface{}) error {
	return status.Errorf(codes.Code(400003), message)

}

// InvoiceNotFound ...
func InvoiceNotFound(message string, details ...interface{}) error {
	return status.Errorf(codes.Code(400004), message)
}

// InternalServerError ...
func InternalServerError(message string, details ...interface{}) error {
	return status.Errorf(codes.Code(InternalServerErrorCode), message)

}

// MissingMetadataError ...
func MissingMetadataError(message string, details ...interface{}) error {
	return status.Errorf(codes.Code(400005), message)
}

// PreconditionFailed ...
func PreconditionFailed(message string, details ...interface{}) error {
	// return status.Errorf(400006, , message)
	return status.Errorf(codes.Code(400006), message)
}

// Canceled client cancel request
func Canceled(message string) error {
	return status.Errorf(codes.Canceled, message)
}

// DeadlineExceeded timeout
func DeadlineExceeded(message string) error {
	return status.Errorf(codes.DeadlineExceeded, message)
}
