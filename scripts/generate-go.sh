#!/bin/bash
set -e

echo "Generating Go code from OpenAPI specs..."

# Define paths
OPENAPI_DIR="api/openapi"
GATEWAY_API_DIR="services/gateway/internal/api"

# Create generated directory if it doesn't exist
mkdir -p "$GATEWAY_API_DIR"

# Check if oapi-codegen is installed
if ! command -v oapi-codegen &> /dev/null; then
    echo "Error: oapi-codegen is not installed."
    echo "Please run: go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest"
    exit 1
fi

# Generate Go code for gateway service (user API)
echo "Generating Go types and Gin server for gateway service..."
oapi-codegen -config "$OPENAPI_DIR/codegen-config.yaml" \
  -o "$GATEWAY_API_DIR/user.gen.go" \
  "$OPENAPI_DIR/user.yaml"

echo "Go code generation completed successfully!"
echo "Generated files:"
echo "  - $GATEWAY_API_DIR/user.gen.go"
