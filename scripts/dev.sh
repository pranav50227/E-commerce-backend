#!/bin/bash
# dev.sh - Starts the project in development mode.

echo "=========================================="
echo "      E-commerce Backend Service Manager  "
echo "=========================================="

# Check if docker daemon is actually running
docker_running=false
if command -v docker &> /dev/null; then
  if docker info &> /dev/null; then
    docker_running=true
  fi
fi

if [ "$docker_running" = true ]; then
  echo "--> Docker is available. Launching via Docker Compose..."
  if command -v docker-compose &> /dev/null; then
    docker-compose up
  else
    docker compose up
  fi
else
  echo "[WARNING]: Docker is not running or not installed."
  echo "--> Falling back to running services locally on the local server..."

  # Run local services in parallel
  
  echo "Starting User Management Service..."
  (cd service/UserManagementService && go run cmd/server/main.go) &
  PID_USER=$!
  
  echo "Starting Product Catalog Service..."
  (cd service/ProductCatalogService && go run cmd/server/main.go) &
  PID_PRODUCT=$!
  
  echo "Starting Inventory Service..."
  (cd service/InventoryService && go run cmd/server/main.go) &
  PID_INVENTORY=$!
  
  echo "Starting Order Management Service..."
  (cd service/OrderManagementService && go run cmd/server/main.go) &
  PID_ORDER=$!
  
  echo "Starting Shopping Cart Service..."
  (cd service/ShoppingCartService && go run cmd/server/main.go) &
  PID_CART=$!
  
  echo "Starting Auth Service..."
  (cd service/auth-service && go run cmd/server/main.go) &
  PID_AUTH=$!
  
  echo "Starting API Gateway..."
  (cd service/api-gateway && go run cmd/gateway/main.go) &
  PID_GATEWAY=$!
  
  # Setup trap to terminate all background processes on exit (Ctrl+C)
  trap "echo 'Stopping all local services...'; kill $PID_USER $PID_PRODUCT $PID_INVENTORY $PID_ORDER $PID_CART $PID_AUTH $PID_GATEWAY 2>/dev/null; exit 0" INT TERM EXIT
  
  echo "------------------------------------------"
  echo "All services successfully launched locally!"
  echo "Gateway is running on port 8000."
  echo "Press Ctrl+C to stop all services."
  echo "------------------------------------------"
  
  # Wait for all background jobs to finish
  wait
fi

