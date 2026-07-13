#!/bin/bash
# lint.sh - Static code analysis and code verification.

echo "=========================================="
echo "          Running Go Code Linters         "
echo "=========================================="

echo "--> Checking shared module..."
(cd shared && go vet ./...)

for dir in service/*; do
  if [ -d "$dir" ]; then
    service_name=$(basename "$dir")
    echo "--> Running vetting/lint check on: $service_name..."
    (cd "$dir" && go vet ./...)
  fi
done

# If golangci-lint is installed, run it
if command -v golangci-lint &> /dev/null; then
  echo "--> Running golangci-lint on workspace..."
  golangci-lint run
else
  echo "[INFO]: golangci-lint is not installed. Standard 'go vet' was used."
fi

echo "=========================================="
echo "Linting complete!"
echo "=========================================="
