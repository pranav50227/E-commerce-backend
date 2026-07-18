# Build stage
FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy shared module
COPY shared/ ./shared/

# Copy service dependencies
COPY service/OrderManagementService/go.mod service/OrderManagementService/go.sum* ./service/OrderManagementService/
WORKDIR /app/service/OrderManagementService
RUN go mod download

# Copy service source code and build
COPY service/OrderManagementService/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/service/OrderManagementService/main .
EXPOSE 8083
CMD ["./main"]

