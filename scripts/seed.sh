#!/bin/bash
# seed.sh - Feeds initial seed data (users, products, and categories) into the running microservices.

echo "=========================================="
echo "         Seeding Database Content         "
echo "=========================================="

GATEWAY_URL="http://localhost:8000"

# 1. Register a test admin user
echo "--> Registering admin user..."
REG_RESP=$(curl -s -w "\n%{http_code}" -X POST "$GATEWAY_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"adminpassword123","name":"System Administrator","email":"admin@example.com"}')

REG_CODE=$(echo "$REG_RESP" | tail -n1)
if [ "$REG_CODE" -ne 201 ] && [ "$REG_CODE" -ne 400 ]; then
  echo "[ERROR]: Failed to register admin (HTTP Code: $REG_CODE). Ensure the stack is running on port 8000."
  exit 1
fi
echo "Admin registered (or already exists)."

# 2. Authenticate admin user to obtain token
echo "--> Logging in to retrieve JWT Access Token..."
LOGIN_RESP=$(curl -s -X POST "$GATEWAY_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"adminpassword123"}')

TOKEN=$(echo "$LOGIN_RESP" | grep -o '"accessToken":"[^"]*' | grep -o '[^"]*$')

if [ -z "$TOKEN" ]; then
  echo "[ERROR]: Could not retrieve JWT Token. Ensure credentials are correct."
  exit 1
fi
echo "Access token successfully retrieved."

# 3. Create Sample Products
echo "--> Seeding sample catalog items..."
declare -a PRODUCTS=(
  '{"name":"Studio Headphones","description":"Wireless over-ear headphones with ANC","price":199.99,"category":"Electronics"}'
  '{"name":"Mechanical Keyboard","description":"Hot-swappable tactile RGB keyboard","price":89.50,"category":"Electronics"}'
  '{"name":"Premium Leather Wallet","description":"Genuine bi-fold leather cardholder","price":45.00,"category":"Apparel"}'
  '{"name":"Ergonomic Office Chair","description":"Breathable mesh high back chair","price":249.99,"category":"Furniture"}'
)

for product in "${PRODUCTS[@]}"; do
  echo "Adding product: $(echo "$product" | grep -o '"name":"[^"]*' | cut -d'"' -f4)"
  curl -s -o /dev/null -X POST "$GATEWAY_URL/api/v1/products/" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "$product"
done

echo "=========================================="
echo "Seeding complete!"
echo "=========================================="
