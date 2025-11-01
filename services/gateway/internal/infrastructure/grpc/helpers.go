package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Riku-KANO/kube-ec/pkg/proto/common"
	"github.com/Riku-KANO/kube-ec/services/gateway/internal/domain/errors"
)

// timestampToTime converts common.Timestamp to time.Time
func timestampToTime(ts *common.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return time.Unix(ts.Seconds, int64(ts.Nanos))
}

// mapGRPCError maps gRPC errors to domain errors
func mapGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return errors.ErrInternalError
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return errors.ErrInvalidInput
	case codes.NotFound:
		return errors.ErrNotFound
	case codes.AlreadyExists:
		return errors.ErrInvalidInput
	case codes.Unauthenticated:
		return errors.ErrUnauthorized
	case codes.PermissionDenied:
		return errors.ErrUnauthorized
	case codes.DeadlineExceeded:
		return errors.ErrInternalError
	case codes.Unavailable:
		return errors.ErrInternalError
	default:
		return errors.ErrInternalError
	}
}

// withTimeout creates a context with timeout for gRPC calls
func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		// Context already has deadline, don't override
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, defaultGRPCTimeout)
}
