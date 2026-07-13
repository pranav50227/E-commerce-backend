# User Management Service

The User Management Service acts as the single source of truth for user profile details and credentials.

## What it does
* **Profile CRUD**: Manages user profiles (Username, Email, Name, Password Hash).
* **Internal APIs**: Exposes high-performance internal endpoints for username and ID queries consumed by the Auth Service.

## Data Flow
* **Profile Queries**: 
  `Client -> API Gateway (Verifies JWT) -> User Service (GET /api/v1/users/{userId})`
* **Internal Creation**: 
  `Auth Service (during registration) -> User Service (POST /internal/users/)`

## Context Passing
* **User Identity Context**: Reads `X-User-Id` header injected by the API Gateway to ensure users can only view or edit their own profile, protecting details from unauthorized cross-tenant requests.

## Covered Aspects
* Clean architecture (Repository layer isolating data storage, Service layer validating email/username uniqueness, Handler layer parsing REST variables).
* In-memory thread-safe thread-lock maps simulating database queries.

## Future Aspects
* **Persistent Database**: Transition from in-memory maps to PostgreSQL or MySQL database with GORM schemas.
* **Avatar Upload**: Integrate AWS S3 or Cloudinary storage for profile pictures.
* **Audit Logging**: Maintain history of profile updates, password resets, and changes for security compliance.
