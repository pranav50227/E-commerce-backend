# Build stage
FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy shared module
COPY shared/ ./shared/

# Copy service dependencies
COPY service/ShoppingCartService/go.mod service/ShoppingCartService/go.sum* ./service/ShoppingCartService/
WORKDIR /app/service/ShoppingCartService
RUN go mod download

# Copy service source code and build
COPY service/ShoppingCartService/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/service/ShoppingCartService/main .
EXPOSE 8084
CMD ["./main"]

