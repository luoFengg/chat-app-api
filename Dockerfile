# ================================
# Stage 1: Build
# ================================
FROM golang:1.21-alpine AS builder
# Set working directory
WORKDIR /app
# Install dependencies
RUN apk add --no-cache git
# Copy go.mod dan go.sum
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download
# Copy semua source code
COPY . .
# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
# ================================
# Stage 2: Run
# ================================
FROM alpine:latest
# Install ca-certificates untuk HTTPS
RUN apk --no-cache add ca-certificates tzdata
# Set timezone
ENV TZ=Asia/Jakarta
WORKDIR /root/
# Copy binary dari stage builder
COPY --from=builder /app/main .
# Copy .env (opsional, bisa pakai env vars)
# COPY --from=builder /app/.env .
# Expose port
EXPOSE 8080
# Run binary
CMD ["./main"]