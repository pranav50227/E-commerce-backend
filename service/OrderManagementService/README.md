# Order Management Service

The Order Management Service tracks, updates, and creates purchase receipts for users.

## What it does
* **Order Creation**: Builds permanent purchase records with unique Order IDs.
* **Status Updates**: Transitions orders through their lifecycle: `Pending` -> `Success` / `Failed` -> `Cancelled`.
* **History Lookup**: Returns orders list specific to a user.

## Data Flow
* **Order Placement**: 
  `Cart Service (Checkout Orchestrator) -> Order Service (POST /api/v1/orders/) -> Saves record -> Returns Receipt.`

## Context Passing
* **User Identity Context**: Fetches `X-User-Id` from request headers. This ensures that users can only view their own transaction history, preventing data leakage.

## Covered Aspects
* Order schemas storing items array with historic prices (recording the exact price paid at the time of purchase rather than current catalog prices).
* Data lifecycle transition logic.

## Future Aspects
* **Payment Integration**: Link third-party payment gateways (Stripe, PayPal, Adyen) to verify payments before transitioning orders to `Success` state.
* **Shipping Integration**: Connect with shipping APIs (EasyPost, FedEx, UPS) to assign tracking numbers and estimate shipping arrivals.
* **Invoice Generation**: Auto-generate invoice PDFs and send them to users via email attachments.
