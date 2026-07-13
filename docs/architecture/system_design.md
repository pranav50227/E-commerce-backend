# System Architecture & Design Documentation

This document describes the architectural patterns, feature sets, scalability strategies, and detailed communication flows within the E-commerce Backend Suite.

---

## Architectural Patterns

The E-commerce Backend suite is built using modern cloud-native architectural patterns designed for microservice autonomy, stateless execution, and high performance.

### 1. API Gateway Pattern
All incoming client traffic enters the system via a single entry point (the **API Gateway** on port `8000`).
* **Routing & Reverse Proxying**: Routes are matched by prefix and forwarded dynamically using a reverse proxy to downstream microservices.
* **Centralized Cross-Cutting Concerns**: Authentication checks, CORS, rate limiting, and request headers enrichment are performed at the gateway level.
* **Security & Auth Offloading**: Prevents individual services from repeating authentication logic.

### 2. Stateless Service Pattern
Every service is entirely stateless. State is managed by in-memory repositories (which map directly to container-attached datastores like Postgres, MySQL, or Redis in production).
* **Benefits**: Enables instant scale-up/scale-down and container restarts without data corruption or configuration state drift.

### 3. Database-per-Service Pattern
To ensure strict decoupling, services do not share a common database.
* **Encapsulation**: Cart, Order, Product, User, and Inventory data are strictly isolated.
* **Cross-Service Verification**: Services query each other's public endpoints (via HTTP) to verify integrity (e.g., the Cart Service requests the Product Service to verify if a product exists before adding it to a cart).

### 4. Layered Clean Architecture (Inside Services)
Each microservice is structured into clean separation of concerns:
```
┌────────────────────────────────────────────────────────┐
│ API / Router (Exposes REST endpoints, parses requests)  │
└───────────────────┬────────────────────────────────────┘
                    ▼
┌────────────────────────────────────────────────────────┐
│ Service Layer (Core business rules, validations)       │
└───────────────────┬────────────────────────────────────┘
                    ▼
┌────────────────────────────────────────────────────────┐
│ Repository / Client (Data access, downstream RPCs)     │
└────────────────────────────────────────────────────────┘
```

---

## Scale & Scalability Strategy

This microservices design is built to scale efficiently under load using the following mechanisms:

### 1. Independent Horizontal Scaling
Because downstream services are decoupled:
* **Product Catalog Service**: Can be scaled to dozens of instances to handle read-heavy catalog searches during shopping sales without wasting resources scaling the write-heavy **Order Management Service**.
* **API Gateway**: Can run multiple replicas behind an external Load Balancer (like AWS ALB or Nginx).

### 2. Shared-Secret JWT Authentication Scalability
* JWT verification uses an HMAC-SHA256 signature checked using a secret key.
* The API Gateway verifies the token cryptographically **without** making a database query or an auth-service HTTP request for every single call. This eliminates the Auth Service database bottleneck.

### 3. Future Scalability Roadmaps (Production Ready)
* **Caching (Redis)**: Introduce Redis caching at the API Gateway or Product Catalog layer to resolve catalog reads instantly.
* **Event-Driven Architecture (Kafka / RabbitMQ)**: Decouple the Checkout flow. Instead of synchronous HTTP chain requests during checkout (Cart -> Order -> Inventory), the Cart Service can emit a `checkout.submitted` event, and the Order and Inventory services can consume it asynchronously to guarantee final consistency.

---

## System Diagrams

A comprehensive backend system's documentation typically contains four key types of diagrams to represent different facets of the system. They are detailed below.

### 1. High-Level Component & Communication Diagram
*Describes the layout of services, their entrypoints, and communication bounds.*

```mermaid
graph TD
    classDef client fill:#3b82f6,stroke:#1d4ed8,color:#fff,stroke-width:2px;
    classDef gateway fill:#10b981,stroke:#047857,color:#fff,stroke-width:2px;
    classDef service fill:#8b5cf6,stroke:#6d28d9,color:#fff,stroke-width:1px;

    Client["📱 Client App"]:::client -->|Port 8000| Gateway["⚡ API Gateway"]:::gateway
    
    subgraph Downstream Microservices
        Auth["🔐 Auth Service :8085"]:::service
        Users["👤 User Service :8080"]:::service
        Products["📦 Product Service :8081"]:::service
        Inventory["🏬 Inventory Service :8082"]:::service
        Orders["🛒 Order Service :8083"]:::service
        Cart["🛍️ Cart Service :8084"]:::service
    end

    Gateway -->|/api/v1/auth| Auth
    Gateway -->|/api/v1/users| Users
    Gateway -->|/api/v1/products| Products
    Gateway -->|/api/v1/inventory| Inventory
    Gateway -->|/api/v1/orders| Orders
    Gateway -->|/api/v1/cart| Cart

    %% Downstream Service Relationships
    Auth -.->|Create/Lookup user| Users
    Cart -.->|Verify product exists| Products
    Cart -.->|Verify stock levels| Inventory
    Cart -.->|Submit checkout| Orders
    Orders -.->|Get prices| Products
    Orders -.->|Deduct stock| Inventory
```

### 2. Authentication & Authorization Flow Sequence
*Shows how the token-based security check is processed statelessly at the Gateway.*

```mermaid
sequenceDiagram
    autonumber
    actor User as Client App
    participant Gateway as API Gateway
    participant Auth as Auth Service
    participant Target as Target Microservice (e.g. Orders)

    %% Registration & Login
    User->>Gateway: POST /api/v1/auth/login
    Gateway->>Auth: Proxy Login credentials
    Auth-->>Gateway: Return JWT Tokens (Access & Refresh)
    Gateway-->>User: Return Tokens

    %% Authorized Request
    User->>Gateway: GET /api/v1/orders (Header Bearer JWT)
    Note over Gateway: Gateway validates JWT cryptographically.<br/>Extracts userId statelessly.
    alt Token Invalid/Expired
        Gateway-->>User: HTTP 401 Unauthorized
    else Token Valid
        Gateway->>Target: Forward Request (Headers: X-User-Id = userId)
        Target-->>Gateway: Return Order Data
        Gateway-->>User: HTTP 200 OK
    end
```

### 3. Shopping Cart Checkout Sequence Diagram
*Shows the chain of interactions between multiple services to successfully place an order.*

```mermaid
sequenceDiagram
    autonumber
    actor Client
    participant Gateway as API Gateway
    participant Cart as Cart Service
    participant Products as Product Service
    participant Inventory as Inventory Service
    participant Orders as Order Service

    Client->>Gateway: POST /api/v1/cart/{userId}/checkout (JWT)
    Gateway->>Cart: Proxy Checkout request
    
    %% Fetch cart details
    Note over Cart: Reads user's cart items
    
    %% Verify Product catalog prices
    Cart->>Products: GET /internal/products/{productId}
    Products-->>Cart: Product Details (Price, Details)
    
    %% Deduct Inventory Stock
    Cart->>Inventory: POST /internal/inventory/deduct (Items, Qty)
    alt Stock Insufficient
        Inventory-->>Cart: HTTP 409 Conflict (Out of Stock)
        Cart-->>Gateway: HTTP 400 Bad Request (Checkout failed)
        Gateway-->>Client: Return error message
    else Stock Deducted Successfully
        Inventory-->>Cart: HTTP 200 OK
        
        %% Create Order
        Cart->>Orders: POST /api/v1/orders (Items, TotalPrice)
        Orders-->>Cart: Order Created (ID, status="Success")
        
        %% Clear user's Cart
        Note over Cart: Empty user's cart items
        Cart-->>Gateway: HTTP 201 Created (Order details)
        Gateway-->>Client: Checkout Complete + Order Receipt
    end
```

### 4. Entity Relation Schema & Context Bounds
*Outlines data models, types, and context borders.*

```mermaid
classDiagram
    class User {
        +String ID
        +String Username
        +String Name
        +String Email
    }
    
    class Product {
        +String ID
        +String Name
        +String Description
        +Double Price
        +String Category
    }
    
    class Inventory {
        +String ProductID
        +Int Quantity
    }

    class CartItem {
        +String ProductID
        +Int Quantity
    }

    class Cart {
        +String UserID
        +CartItem[] Items
    }

    class OrderItem {
        +String ProductID
        +Int Quantity
        +Double Price
    }

    class Order {
        +String ID
        +String UserID
        +OrderItem[] Items
        +Double TotalPrice
        +String Status
        +DateTime CreatedAt
    }

    User "1" --> "1" Cart : Owns
    Cart "1" *-- "many" CartItem : Contains
    Order "1" *-- "many" OrderItem : Contains
    User "1" --> "many" Order : Places
    Product "1" --> "1" Inventory : Has stock tracking
```
