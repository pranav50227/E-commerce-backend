#!/bin/bash
# dev.sh - Starts the project in development mode.

echo "--> Launching E-commerce Backend Services..."

if command -v docker-compose &> /dev/null; then
  docker-compose up
elif command -v docker &> /dev/null; then
  docker compose up
else
  echo "[ERROR]: Docker not found. Starting microservices manually is recommended if Docker is unavailable."
  exit 1
fi
