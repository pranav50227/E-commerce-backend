# Build stage
FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy shared module
COPY shared/ ./shared/

# Copy service dependencies
COPY service/api-gateway/go.mod service/api-gateway/go.sum* ./service/api-gateway/
WORKDIR /app/service/api-gateway
RUN go mod download

# Copy service source code and build
COPY service/api-gateway/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/gateway/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/service/api-gateway/main .
EXPOSE 8000
CMD ["./main"]
