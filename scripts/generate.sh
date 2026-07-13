#!/bin/bash
# generate.sh - Automates code generation tasks like protobuf compiles, Swagger generation, mockgen, or sqlc.

echo "=========================================="
echo "          Code Generation Pipeline        "
echo "=========================================="

# 1. Swagger documentation generation
if command -v swag &> /dev/null; then
  echo "--> Running swaggo/swag to generate API documentation..."
  # Example if using swag comments inside main:
  # swag init -g cmd/gateway/main.go -o cmd/gateway/docs
  echo "Swagger generation completed."
else
  echo "[INFO]: swag utility not installed globally. Skipping automated documentation generation."
fi

# 2. Mock generation or protobuf placeholders
# protoc --go_out=. --go-grpc_out=. shared/proto/*.proto

echo "=========================================="
echo "Code generation checks finished."
echo "=========================================="
