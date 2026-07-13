#!/bin/bash
# build.sh - Compiles all modules and service binaries.

echo "=========================================="
echo "    Building All Backend Services         "
echo "=========================================="

echo "--> Building shared module..."
(cd shared && go build ./...)

for dir in service/*; do
  if [ -d "$dir" ]; then
    service_name=$(basename "$dir")
    echo "--> Compiling service binary: $service_name..."
    (cd "$dir" && go build -o "${service_name}_bin" ./cmd/* || go build -o "${service_name}_bin" ./...)
  fi
done

echo "=========================================="
echo "Build complete!"
echo "=========================================="
