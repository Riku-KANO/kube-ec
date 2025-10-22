#!/bin/bash
set -e

echo "Generating Go code from protobuf files..."

# Go to project root
cd "$(dirname "$0")/.."

# Generate common proto
echo "Generating common proto..."
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --proto_path=. \
  proto/common/common.proto

# Generate user proto
echo "Generating user proto..."
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --proto_path=. \
  proto/user/user.proto

# Generate product proto
echo "Generating product proto..."
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --proto_path=. \
  proto/product/product.proto

# Generate order proto
echo "Generating order proto..."
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --proto_path=. \
  proto/order/order.proto

# Generate payment proto
echo "Generating payment proto..."
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --proto_path=. \
  proto/payment/payment.proto

echo "Proto generation completed successfully!"
