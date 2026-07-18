# Build stage
FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy shared module
COPY shared/ ./shared/

# Copy service dependencies
COPY service/auth-service/go.mod service/auth-service/go.sum* ./service/auth-service/
WORKDIR /app/service/auth-service
RUN go mod download

# Copy service source code and build
COPY service/auth-service/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/service/auth-service/main .
EXPOSE 8085
CMD ["./main"]

