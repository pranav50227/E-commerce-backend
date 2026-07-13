# E-commerce Backend API Reference

This document serves as a reference for all API endpoints exposed by the API Gateway on port `8000`.

---

## Authentication Requirements

Endpoints marked with 🔐 require a JWT Bearer token in the `Authorization` header:
```http
Authorization: Bearer <your_jwt_access_token>
```
Endpoints marked with 🔓 are public and bypass authentication check.

---

## 1. Authentication Service (`/api/v1/auth`)

### Register User
* **Method**: `POST`
* **Route**: `/api/v1/auth/register`
* **Auth**: 🔓 Public
* **Request Body**:
```json
{
  "username": "johndoe",
  "password": "securepassword123",
  "name": "John Doe",
  "email": "johndoe@example.com"
}
```
* **Success Response (`201 Created`)**:
```json
{
  "id": "usr-12345",
  "username": "johndoe",
  "name": "John Doe",
  "email": "johndoe@example.com"
}
```

### Login
* **Method**: `POST`
* **Route**: `/api/v1/auth/login`
* **Auth**: 🔓 Public
* **Request Body**:
```json
{
  "username": "johndoe",
  "password": "securepassword123"
}
```
* **Success Response (`200 OK`)**:
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Refresh Token
* **Method**: `POST`
* **Route**: `/api/v1/auth/refresh-token`
* **Auth**: 🔓 Public
* **Request Body**:
```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
* **Success Response (`200 OK`)**:
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## 2. User Service (`/api/v1/users`)

### Get User Profile
* **Method**: `GET`
* **Route**: `/api/v1/users/{userId}`
* **Auth**: 🔐 Authorized
* **Success Response (`200 OK`)**:
```json
{
  "id": "usr-12345",
  "username": "johndoe",
  "name": "John Doe",
  "email": "johndoe@example.com"
}
```

### Update User Profile
* **Method**: `PUT`
* **Route**: `/api/v1/users/{userId}`
* **Auth**: 🔐 Authorized
* **Request Body**:
```json
{
  "name": "Johnathan Doe",
  "email": "johndoe.new@example.com"
}
```
* **Success Response (`200 OK`)**:
```json
{
  "id": "usr-12345",
  "username": "johndoe",
  "name": "Johnathan Doe",
  "email": "johndoe.new@example.com"
}
```

### Delete User Profile
* **Method**: `DELETE`
* **Route**: `/api/v1/users/{userId}`
* **Auth**: 🔐 Authorized
* **Success Response (`200 OK`)**:
```json
{
  "message": "User deleted successfully"
}
```

---

## 3. Product Catalog Service (`/api/v1/products`)

### List All Products
* **Method**: `GET`
* **Route**: `/api/v1/products/`
* **Auth**: 🔓 Public
* **Success Response (`200 OK`)**:
```json
[
  {
    "id": "prod1",
    "name": "Wireless Headphones",
    "description": "Noise-cancelling over-ear headphones",
    "price": 199.99,
    "category": "Electronics"
  }
]
```

### Get Product by ID
* **Method**: `GET`
* **Route**: `/api/v1/products/{productId}`
* **Auth**: 🔓 Public
* **Success Response (`200 OK`)**:
```json
{
  "id": "prod1",
  "name": "Wireless Headphones",
  "description": "Noise-cancelling over-ear headphones",
  "price": 199.99,
  "category": "Electronics"
}
```

### Search Products
* **Method**: `GET`
* **Route**: `/api/v1/products/search?q={query}`
* **Auth**: 🔓 Public
* **Success Response (`200 OK`)**:
```json
[
  {
    "id": "prod1",
    "name": "Wireless Headphones",
    "description": "Noise-cancelling over-ear headphones",
    "price": 199.99,
    "category": "Electronics"
  }
]
```

### Get Products by Category
* **Method**: `GET`
* **Route**: `/api/v1/products/category/{categoryName}`
* **Auth**: 🔓 Public
* **Success Response (`200 OK`)**:
```json
[
  {
    "id": "prod1",
    "name": "Wireless Headphones",
    "price": 199.99,
    "category": "Electronics"
  }
]
```

### Create Product
* **Method**: `POST`
* **Route**: `/api/v1/products/`
* **Auth**: 🔐 Authorized
* **Request Body**:
```json
{
  "name": "Mechanical Keyboard",
  "description": "RGB Backlit keyboard with Blue switches",
  "price": 89.99,
  "category": "Electronics"
}
```
* **Success Response (`201 Created`)**:
```json
{
  "id": "prod-56789",
  "name": "Mechanical Keyboard",
  "description": "RGB Backlit keyboard with Blue switches",
  "price": 89.99,
  "category": "Electronics"
}
```

### Update Product
* **Method**: `PUT`
* **Route**: `/api/v1/products/{productId}`
* **Auth**: 🔐 Authorized
* **Request Body**:
```json
{
  "name": "Mechanical Keyboard Pro",
  "description": "Silent mechanical keyboard",
  "price": 109.99,
  "category": "Electronics"
}
```
* **Success Response (`200 OK`)**:
```json
{
  "id": "prod-56789",
  "name": "Mechanical Keyboard Pro",
  "description": "Silent mechanical keyboard",
  "price": 109.99,
  "category": "Electronics"
}
```

### Delete Product
* **Method**: `DELETE`
* **Route**: `/api/v1/products/{productId}`
* **Auth**: 🔐 Authorized
* **Success Response (`200 OK`)**:
```json
{
  "message": "Product deleted successfully"
}
```

---

## 4. Inventory Service (`/api/v1/inventory`)

### Get Entire Inventory
* **Method**: `GET`
* **Route**: `/api/v1/inventory/`
* **Auth**: 🔐 Authorized
* **Success Response (`200 OK`)**:
```json
{
  "prod1": 100,
  "prod2": 50,
  "prod3": 200
}
```

### Get Stock of Product
* **Method**: `GET`
* **Route**: `/api/v1/inventory/{productId}`
* **Auth**: 🔐 Authorized
* **Success Response (`200 OK`)**:
```json
{
  "productId": "prod1",
  "quantity": 100
}
```

### Set Stock of Product
* **Method**: `PUT`
* **Route**: `/api/v1/inventory/{productId}`
* **Auth**: 🔐 Authorized
* **Request Body**:
```json
{
  "quantity": 120
}
```
* **Success Response (`200 OK`)**:
```json
{
  "productId": "prod1",
  "quantity": 120,
  "message": "Inventory updated"
}
```

### Restock Inventory
* **Method**: `POST`
* **Route**: `/api/v1/inventory/restock`
* **Auth**: 🔐 Authorized
* **Request Body**:
```json
{
  "productId": "prod1",
  "quantity": 30
}
```
* **Success Response (`200 OK`)**:
```json
{
  "productId": "prod1",
  "quantity": 150,
  "message": "Product restocked successfully"
}
```

---

## 5. Shopping Cart Service (`/api/v1/cart`)

### Get Cart
* **Method**: `GET`
* **Route**: `/api/v1/cart/{userId}`
* **Auth**: 🔐 Authorized
* **Success Response (`200 OK`)**:
```json
{
  "userId": "usr-12345",
  "items": [
    {
      "productId": "prod1",
      "quantity": 2
    }
  ]
}
```

### Add Item to Cart
* **Method**: `POST`
* **Route**: `/api/v1/cart/{userId}/items`
* **Auth**: 🔐 Authorized
* **Request Body**:
```json
{
  "productId": "prod1",
  "quantity": 2
}
```
* **Success Response (`200 OK`)**:
```json
{
  "userId": "usr-12345",
  "items": [
    {
      "productId": "prod1",
      "quantity": 2
    }
  ]
}
```

### Update Item in Cart
* **Method**: `PUT`
* **Route**: `/api/v1/cart/{userId}/items/{productId}`
* **Auth**: 🔐 Authorized
* **Request Body**:
```json
{
  "quantity": 5
}
```
* **Success Response (`200 OK`)**:
```json
{
  "userId": "usr-12345",
  "items": [
    {
      "productId": "prod1",
      "quantity": 5
    }
  ]
}
```

### Remove Item from Cart
* **Method**: `DELETE`
* **Route**: `/api/v1/cart/{userId}/items/{productId}`
* **Auth**: 🔐 Authorized
* **Success Response (`200 OK`)**:
```json
{
  "userId": "usr-12345",
  "items": []
}
```

### Clear Cart
* **Method**: `DELETE`
* **Route**: `/api/v1/cart/{userId}`
* **Auth**: 🔐 Authorized
* **Success Response (`200 OK`)**:
```json
{
  "message": "Cart cleared successfully"
}
```

### Checkout Cart
* **Method**: `POST`
* **Route**: `/api/v1/cart/{userId}/checkout`
* **Auth**: 🔐 Authorized
* **Success Response (`201 Created`)**:
```json
{
  "id": "ord-88392",
  "userId": "usr-12345",
  "items": [
    {
      "productId": "prod1",
      "quantity": 5,
      "price": 199.99
    }
  ],
  "totalPrice": 999.95,
  "status": "Success",
  "createdAt": "2026-07-13T10:19:08Z"
}
```

---

## 6. Order Management Service (`/api/v1/orders`)

### Get Orders List
* **Method**: `GET`
* **Route**: `/api/v1/orders/`
* **Auth**: 🔐 Authorized
* **Success Response (`200 OK`)**:
```json
[
  {
    "id": "ord-88392",
    "userId": "usr-12345",
    "items": [
      {
        "productId": "prod1",
        "quantity": 5,
        "price": 199.99
      }
    ],
    "totalPrice": 999.95,
    "status": "Success",
    "createdAt": "2026-07-13T10:19:08Z"
  }
]
```

### Get Order details
* **Method**: `GET`
* **Route**: `/api/v1/orders/{orderId}`
* **Auth**: 🔐 Authorized
* **Success Response (`200 OK`)**:
```json
{
  "id": "ord-88392",
  "userId": "usr-12345",
  "items": [
    {
      "productId": "prod1",
      "quantity": 5,
      "price": 199.99
    }
  ],
  "totalPrice": 999.95,
  "status": "Success",
  "createdAt": "2026-07-13T10:19:08Z"
}
```

### Create Order Manually
* **Method**: `POST`
* **Route**: `/api/v1/orders/`
* **Auth**: 🔐 Authorized
* **Request Body**:
```json
{
  "items": [
    {
      "productId": "prod1",
      "quantity": 1
    }
  ]
}
```
* **Success Response (`201 Created`)**:
```json
{
  "id": "ord-99482",
  "userId": "usr-12345",
  "items": [
    {
      "productId": "prod1",
      "quantity": 1,
      "price": 199.99
    }
  ],
  "totalPrice": 199.99,
  "status": "Success",
  "createdAt": "2026-07-13T10:19:08Z"
}
```

### Update Order Status
* **Method**: `PUT`
* **Route**: `/api/v1/orders/{orderId}/status`
* **Auth**: 🔐 Authorized
* **Request Body**:
```json
{
  "status": "Cancelled"
}
```
* **Success Response (`200 OK`)**:
```json
{
  "id": "ord-99482",
  "userId": "usr-12345",
  "items": [
    {
      "productId": "prod1",
      "quantity": 1,
      "price": 199.99
    }
  ],
  "totalPrice": 199.99,
  "status": "Cancelled",
  "createdAt": "2026-07-13T10:19:08Z"
}
```

### Cancel Order
* **Method**: `DELETE`
* **Route**: `/api/v1/orders/{orderId}`
* **Auth**: 🔐 Authorized
* **Success Response (`200 OK`)**:
```json
{
  "message": "Order cancelled successfully"
}
```
