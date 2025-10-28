FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Copy proto and pkg directories first (needed by go.mod replace directive)
COPY proto/ ./proto/
COPY pkg/ ./pkg/

# Copy auth service
WORKDIR /workspace/services/auth
COPY services/auth/go.mod services/auth/go.sum* ./
RUN go mod download

COPY services/auth/ ./

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-service ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /workspace/services/auth/auth-service .

EXPOSE 50052

CMD ["./auth-service"]
