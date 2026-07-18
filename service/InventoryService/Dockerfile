# Build stage
FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy shared module
COPY shared/ ./shared/

# Copy service dependencies
COPY service/InventoryService/go.mod service/InventoryService/go.sum* ./service/InventoryService/
WORKDIR /app/service/InventoryService
RUN go mod download

# Copy service source code and build
COPY service/InventoryService/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/service/InventoryService/main .
EXPOSE 8082
CMD ["./main"]

