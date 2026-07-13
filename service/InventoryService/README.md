# Inventory Service

The Inventory Service monitors warehouse stock levels and performs stock deductions during transactions.

## What it does
* **Stock Tracking**: Maintains record counts of available items.
* **Atomic Deductions**: Provides thread-safe, all-or-nothing stock deductions for cart checkouts.
* **Restock Operations**: Handles replenishment updates.

## Data Flow
* **Stock Deduction**: 
  `Cart Service (checkout trigger) -> Inventory Service (POST /internal/inventory/deduct) -> Checks if all products are in stock -> Subtracts quantities -> Returns OK.`

## Context Passing
* **User Identity Context**: Actions like viewing/updating stock require the Gateway authenticated credentials. Internal deductions are carried out by other microservice calls passing correlation/transaction details.

## Covered Aspects
* Thread-safe atomic updates using mutex locks (`sync.Mutex`) preventing overselling and race conditions on hot items.
* Mock database queries returning real-time stocks.

## Future Aspects
* **Distributed Locking (Redis Lock)**: Replace simple Go Mutexes with distributed locks to allow scaling the service to multiple servers safely.
* **Low Stock Alerts**: Publish messages to Kafka/RabbitMQ when stock drops below threshold (e.g. notify admins).
* **Multi-Warehouse Support**: Route stock deductions to the warehouse closest to the client's location.
