# Authentication Service

The Authentication Service handles user security, credential verification, and issues stateful/stateless authentication tokens.

## What it does
* **User Registration**: Coordinates with the User Management Service to create new accounts.
* **Authentication**: Verifies user passwords and issues JWT Access and Refresh tokens.
* **Session Refreshing**: Grants new access tokens using a secure refresh token rotation loop.

## Data Flow
* **Registration**: 
  `Client -> API Gateway -> Auth Service (validates inputs) -> User Service (POST /internal/users/ to create profile) -> Returns User profile.`
* **Login**:
  `Client -> API Gateway -> Auth Service -> User Service (GET /internal/users/username/{username} to fetch hash) -> Verifies Password -> Returns JWT tokens.`

## Context Passing
* **User Identity Context**: Initiates user identity creation. During authentication, the Auth Service signs the user's ID into the JWT payload, which the API Gateway subsequently parses and forwards downstream via the `X-User-Id` header.

## Covered Aspects
* Password hashing using secure crypt/bcrypt packages.
* JWT signing with custom payload claims (`userId`, `username`, `exp`).
* Inter-service HTTP calls using Rest clients to query user records from User Management Service.

## Future Aspects
* **OAuth2 / OpenID Connect**: Add Google, Apple, and GitHub SSO provider options.
* **Multi-Factor Authentication (MFA)**: Support OTP generators (Google Authenticator) or SMS verification.
* **Token Blacklisting**: Implement Redis-based token revocation list to allow immediate user logouts.
