#!/bin/bash
# setup.sh - Runs once after cloning the repository.

echo "=========================================="
echo "      E-commerce Backend Setup Tool       "
echo "=========================================="

# 1. Download Go dependencies for shared library
echo "--> Tidying and downloading shared dependencies..."
(cd shared && go mod tidy && go mod download)

# 2. Download Go dependencies for each microservice
for dir in service/*; do
  if [ -d "$dir" ]; then
    service_name=$(basename "$dir")
    echo "--> Tidying and downloading dependencies for service: $service_name..."
    (cd "$dir" && go mod tidy && go mod download)
  fi
done

# 3. Spin up infrastructure container dependencies
echo "--> Launching infrastructure containers in background..."
if command -v docker-compose &> /dev/null; then
  docker-compose up -d
elif command -v docker &> /dev/null; then
  docker compose up -d
else
  echo "[WARNING]: Docker Compose not detected. Skipping infrastructure setup."
fi

# 4. Install dev and code generation utilities
echo "--> Installing code generation utilities..."
go install github.com/swaggo/swag/cmd/swag@latest

echo "=========================================="
echo "Setup complete! Ready for development."
echo "=========================================="
