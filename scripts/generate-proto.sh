#!/bin/bash
set -e

echo "Generating Go code from protobuf files..."

# Go to project root
cd "$(dirname "$0")/.."

# Generate common proto
echo "Generating common proto..."
protoc --go_out=. --go_opt=paths=source_relative \
  --proto_path=. \
  pkg/proto/common/common.proto

# Generate auth proto
echo "Generating auth proto..."
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --proto_path=. \
  services/auth/proto/auth.proto

# Generate user proto
echo "Generating user proto..."
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --proto_path=. \
  services/user/proto/user.proto

# Generate product proto
echo "Generating product proto..."
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --proto_path=. \
  services/product/proto/product.proto

# Generate order proto
echo "Generating order proto..."
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --proto_path=. \
  services/order/proto/order.proto

# Generate payment proto
echo "Generating payment proto..."
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  --proto_path=. \
  services/payment/proto/payment.proto

echo "Proto generation completed successfully!"
