# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies including swag
RUN apk add --no-cache git make swag

# Set working directory
WORKDIR /app

# Copy go mod files for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Generate swagger docs and build the application
RUN make build

# Runtime stage
FROM alpine:latest

# Install runtime dependencies (tzdata as used in original build)
RUN apk add --no-cache tzdata

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/out ./out

# Set the entrypoint
ENTRYPOINT ["./out"]