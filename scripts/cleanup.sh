#!/bin/bash
# cleanup.sh - Deletes local build artifacts, logs, caches, and test binaries.

echo "--> Cleaning up local build artifacts and temporary files..."

# Delete compiled binaries in service directories
for dir in service/*; do
  if [ -d "$dir" ]; then
    service_name=$(basename "$dir")
    echo "Cleaning $service_name binaries..."
    rm -f "$dir/${service_name}_bin" "$dir/server.exe" "$dir/gateway.exe" "$dir/product.exe" "$dir/server"
  fi
done

# Clear logs or temporary directories if they exist
rm -rf tmp/ logs/ cache/

# Clear Go test cache
echo "--> Cleaning Go test cache..."
go clean -testcache

echo "Cleanup completed successfully!"
