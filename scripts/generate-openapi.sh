#!/bin/bash
set -e

echo "Generating code from OpenAPI specifications..."

# Define paths
OPENAPI_DIR="api/openapi"
GENERATED_DIR="$OPENAPI_DIR/generated"

# Create generated directory if it doesn't exist
mkdir -p "$GENERATED_DIR/go"
mkdir -p "$GENERATED_DIR/typescript"

# Generate Go code for user service
echo "Generating Go types and Gin server for user service..."
oapi-codegen -config "$OPENAPI_DIR/codegen-config.yaml" \
  -o "$GENERATED_DIR/go/user.gen.go" \
  "$OPENAPI_DIR/user.yaml"

echo "Go code generation completed!"

# Generate TypeScript types for frontend
echo "Generating TypeScript types for Next.js frontend..."
if command -v npx &> /dev/null; then
  npx openapi-typescript "$OPENAPI_DIR/user.yaml" -o "$GENERATED_DIR/typescript/user.ts"
  echo "TypeScript code generation completed!"
else
  echo "Warning: npm/npx not found. Skipping TypeScript generation."
  echo "Run 'npm install -g openapi-typescript' to enable TypeScript generation."
fi

echo "All code generation completed successfully!"
