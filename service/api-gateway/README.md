# API Gateway Service

The API Gateway is the central entry point and reverse proxy for all client traffic. It manages routing, security, and propagates authenticated user context downstream.

## What it does
* **Dynamic Reverse Proxy**: Intercepts HTTP requests and forwards them to downstream microservices based on URL prefixes.
* **Authentication Offloading**: Parses and cryptographically validates JWT access tokens in incoming `Authorization` headers.
* **Context Propagation**: Extracts the `userId` claim from valid tokens and injects it as an `X-User-Id` header before routing.
* **Unified API Documentation**: Hosts the integrated Swagger UI panel on port `8000`.

## Context Passing & Propagation
1. **Client Request**: Client sends a request with `Authorization: Bearer <token>`.
2. **Gateway Validation**: The gateway checks the signature using a shared JWT secret.
3. **Context Injection**: The gateway injects the header `X-User-Id: <userId>` into the request.
4. **Downstream Consumption**: The downstream microservice retrieves the context using:
   ```go
   userId := c.GetHeader("X-User-Id")
   ```

## Covered Aspects
* Route prefix configuration mapping paths to internal microservice URLs.
* HMAC-SHA256 JWT signature verification and token expiration checking.
* Unified Swagger document embedding and UI serving.
* Reverse proxy redirection preserving host headers and query parameters.

## Future Aspects
* **Rate Limiting**: Implement a Redis-backed token bucket algorithm to protect downstream services from DDoS.
* **SSL/TLS Termination**: Secure endpoint communication directly at the gateway layer.
* **Correlation IDs**: Generate and append `X-Correlation-Id` headers to requests to enable trace logs aggregation across services.
* **Service Discovery**: Integrate Consul or Eureka to remove hardcoded downstream URLs.
