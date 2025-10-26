FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Copy proto and pkg directories first (needed by go.mod replace directive)
COPY proto/ ./proto/
COPY pkg/ ./pkg/

# Copy user service
WORKDIR /workspace/services/user
COPY services/user/go.mod services/user/go.sum* ./
RUN go mod download

COPY services/user/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /workspace/services/user/user-service .

EXPOSE 50051

CMD ["./user-service"]
