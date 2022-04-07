package errors

import (
	"log"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// InternalServerErrorCode internal server error code
	InternalServerErrorCode = 500
)

// BadRequest bad request
func BadRequest(message string, details ...interface{}) error {
	log.Printf("bad request: %v", message)
	st := status.New(codes.InvalidArgument, message)
	errInfo := &errdetails.ErrorInfo{
		Reason: "400001",
	}
	// return status.Errorf(codes.InvalidArgument, message)
	ds, err := st.WithDetails(errInfo)
	if err == nil {
		return ds.Err()
	}

	return st.Err()
}

// UnpaidOrderExists ...
func UnpaidOrderExists(message string, details ...interface{}) error {
	// st := status.Newf(codes.Code(400002), message)
	// ds, err := st.WithDetails()
	// if err != nil {
	// 	return ds.Err()
	// }
	log.Printf("unpaid order exists: %v", message)

	st := status.New(codes.AlreadyExists, message)
	errInfo := &errdetails.ErrorInfo{
		Reason: "400002",
	}

	ds, err := st.WithDetails(errInfo)
	if err == nil {
		return ds.Err()
	}

	return st.Err()
	// return status.Errorf(codes.AlreadyExists, message)
}

// OrderNotFound ...
func OrderNotFound(message string, details ...interface{}) error {
	log.Printf("order not found: %v", message)
	// return status.Errorf(codes.NotFound, message)
	st := status.New(codes.NotFound, message)
	// 传递自定义错误码
	errInfo1 := &errdetails.ErrorInfo{
		Reason: "400003",
	}
	// bd := &rpc.BadRequest{
	// 	FieldViolations: []*rpc.BadRequest_FieldViolation{{
	// 		Field:       "err_code",
	// 		Description: "400003",
	// 	}},
	// }
	ds, err := st.WithDetails(errInfo1)
	if err == nil {
		return ds.Err()
	}

	return st.Err()
}

// InvoiceNotFound ...
func InvoiceNotFound(message string, details ...interface{}) error {
	log.Printf("invoice not found: %v", message)
	st := status.New(codes.NotFound, message)
	errInfo := &errdetails.ErrorInfo{
		Reason: "400004",
	}
	ds, err := st.WithDetails(errInfo)
	if err == nil {
		return ds.Err()
	}

	return st.Err()
	// return status.Errorf(codes.NotFound, message)
}

// InternalServerError ...
func InternalServerError(message string, details ...interface{}) error {
	log.Printf("Internal Server Error: %v", message)
	return status.Errorf(codes.Internal, message)

}

// MissingMetadataError ...
func MissingMetadataError(message string, details ...interface{}) error {
	return status.Errorf(codes.InvalidArgument, message)
}

// PreconditionFailed ...
func PreconditionFailed(message string, details ...interface{}) error {
	// return status.Errorf(400006, , message)
	st := status.New(codes.FailedPrecondition, message)
	errInfo := &errdetails.ErrorInfo{
		Reason: "400005",
	}
	ds, err := st.WithDetails(errInfo)
	if err == nil {
		return ds.Err()
	}

	return st.Err()
	// return status.Errorf(codes.FailedPrecondition, message)
}

// Canceled client cancel request
func Canceled(message string) error {
	return status.Errorf(codes.Canceled, message)
}

// DeadlineExceeded timeout
func DeadlineExceeded(message string) error {
	return status.Errorf(codes.DeadlineExceeded, message)
}
