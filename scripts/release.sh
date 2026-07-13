#!/bin/bash
# release.sh - Helper script for building production images.

echo "=========================================="
echo "      Production Release Builder          "
echo "=========================================="

REGISTRY="myregistry.local"
IMAGE_TAG=$(git rev-parse --short HEAD 2>/dev/null || echo "latest")

echo "--> Tagging images with version: $IMAGE_TAG"

# Build all Docker services defined in compose
if command -v docker &> /dev/null; then
  echo "--> Building production-ready containers..."
  docker compose build --no-cache
  
  # For each service, tag and push (if registry is configured)
  # echo "--> Pushing images to registry: $REGISTRY..."
  # docker tag e-commerce-backend-api-gateway $REGISTRY/api-gateway:$IMAGE_TAG
  # docker push $REGISTRY/api-gateway:$IMAGE_TAG
  
  echo "Production build step finished successfully."
else
  echo "[ERROR]: Docker environment is not active. Unable to run release builds."
  exit 1
fi

echo "=========================================="
echo "Release builder complete!"
echo "=========================================="
