#!/bin/bash
# test.sh - Runs test suite for all microservices and shared libraries.

echo "=========================================="
echo "    Running All Backend Test Suites       "
echo "=========================================="

echo "--> Testing shared module..."
(cd shared && go test -v ./...)

for dir in service/*; do
  if [ -d "$dir" ]; then
    service_name=$(basename "$dir")
    echo "--> Running tests for: $service_name..."
    (cd "$dir" && go test -v ./...)
  fi
done

echo "=========================================="
echo "Tests run completed!"
echo "=========================================="
