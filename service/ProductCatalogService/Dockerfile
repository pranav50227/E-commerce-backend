# Build stage
FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy shared module
COPY shared/ ./shared/

# Copy service dependencies
COPY service/ProductCatalogService/go.mod service/ProductCatalogService/go.sum* ./service/ProductCatalogService/
WORKDIR /app/service/ProductCatalogService
RUN go mod download

# Copy service source code and build
COPY service/ProductCatalogService/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/service/ProductCatalogService/main .
EXPOSE 8081
CMD ["./main"]

