package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound          = New(codes.NotFound, "resource not found")
	ErrAlreadyExists     = New(codes.AlreadyExists, "resource already exists")
	ErrInvalidArgument   = New(codes.InvalidArgument, "invalid argument")
	ErrUnauthenticated   = New(codes.Unauthenticated, "unauthenticated")
	ErrPermissionDenied  = New(codes.PermissionDenied, "permission denied")
	ErrInternal          = New(codes.Internal, "internal server error")
	ErrUnavailable       = New(codes.Unavailable, "service unavailable")
	ErrInsufficientStock = New(codes.FailedPrecondition, "insufficient stock")
)

type Error struct {
	Code    codes.Code
	Message string
}

func New(code codes.Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

func FromGRPCError(err error) *Error {
	st, ok := status.FromError(err)
	if !ok {
		return &Error{
			Code:    codes.Unknown,
			Message: err.Error(),
		}
	}
	return &Error{
		Code:    st.Code(),
		Message: st.Message(),
	}
}
