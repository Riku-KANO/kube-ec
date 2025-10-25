FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Copy proto directory first (needed by go.mod replace directive)
COPY proto/ ./proto/

# Copy gateway service
WORKDIR /workspace/services/gateway
COPY services/gateway/go.mod services/gateway/go.sum* ./
RUN go mod download

COPY services/gateway/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway-service ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /workspace/services/gateway/gateway-service .

EXPOSE 8080

CMD ["./gateway-service"]
