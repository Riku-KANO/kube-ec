package grpc

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey string

const (
	requestIDKey contextKey = "request_id"
)

// UnaryRequestIDInterceptor adds request ID to context for tracing
func UnaryRequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		requestID := extractOrGenerateRequestID(ctx)
		ctx = context.WithValue(ctx, requestIDKey, requestID)

		log.Printf("[%s] %s called", requestID, info.FullMethod)

		resp, err := handler(ctx, req)

		if err != nil {
			log.Printf("[%s] %s failed: %v", requestID, info.FullMethod, err)
		} else {
			log.Printf("[%s] %s completed", requestID, info.FullMethod)
		}

		return resp, err
	}
}

// extractOrGenerateRequestID extracts request ID from metadata or generates a new one
func extractOrGenerateRequestID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if ids := md.Get("x-request-id"); len(ids) > 0 {
			return ids[0]
		}
	}

	// Generate new request ID
	return uuid.New().String()
}

// GetRequestID retrieves request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return "unknown"
}
