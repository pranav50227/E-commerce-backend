# Shopping Cart Service

The Shopping Cart Service manages active user shopping carts and acts as the orchestrator during the checkout transaction.

## What it does
* **Cart Operations**: Allows users to add, update quantity, remove, and view items in their active cart.
* **Checkout Orchestration**: Coordinates item validation (Product Service), stock deductions (Inventory Service), and order placements (Order Service) into a single cohesive checkout transaction.

## Data Flow (Checkout Flow)
1. **Trigger**: 
   `Client -> Gateway -> Cart Service (POST /api/v1/cart/{userId}/checkout)`
2. **Product Price Check**: 
   `Cart Service -> Product Service (GET /internal/products/{productId})`
3. **Inventory Reservation**: 
   `Cart Service -> Inventory Service (POST /internal/inventory/deduct)`
4. **Order Placement**: 
   `Cart Service -> Order Service (POST /api/v1/orders/)`
5. **Clear Cart**: Cart resets item list for the user.

## Context Passing
* **User Identity Context**: Reads `X-User-Id` header from the API Gateway. This header determines whose cart is being modified or checked out, securing user cart details from unauthorized cross-reads.

## Covered Aspects
* Shopping cart data structure (Items array containing Product ID and quantity mapped to User IDs).
* Orchestration client code mapping REST connections to other downstream services.

## Future Aspects
* **Persistent Cache**: Store cart items in a persistent Redis cache with Time-To-Live (TTL) configuration.
* **Abandoned Cart Workflows**: Trigger cron events reminding users to complete checkout for cart items left idle for over 24 hours.
* **Promo Code Validations**: Integrate promotional discount coupon validations into the pricing calculations.
