# Build stage
FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy shared module
COPY shared/ ./shared/

# Copy service dependencies
COPY service/UserManagementService/go.mod service/UserManagementService/go.sum* ./service/UserManagementService/
WORKDIR /app/service/UserManagementService
RUN go mod download

# Copy service source code and build
COPY service/UserManagementService/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/service/UserManagementService/main .
EXPOSE 8080
CMD ["./main"]

