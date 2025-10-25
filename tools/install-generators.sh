#!/bin/bash
set -e

echo "Installing OpenAPI code generators..."

# Install oapi-codegen for Go
echo "Installing oapi-codegen..."
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest

echo "All generators installed successfully!"
