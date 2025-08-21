# Multi-stage build for GG Poker application
FROM golang:1.24-alpine AS go-builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the Go application
RUN make build

# Build stage for frontend
FROM node:18-alpine AS frontend-builder

# Set working directory
WORKDIR /app

# Copy package files
COPY web/package*.json ./

# Install dependencies
RUN npm ci --only=production

# Copy source code
COPY web/ ./

# Build the frontend
RUN npm run build

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S poker && \
    adduser -u 1001 -S poker -G poker

# Set working directory
WORKDIR /app

# Copy Go binary from builder
COPY --from=go-builder /app/bin/ggpoker ./ggpoker

# Copy frontend build from frontend-builder
COPY --from=frontend-builder /app/dist ./web/dist

# Copy configuration
COPY config.yaml ./

# Create necessary directories
RUN mkdir -p logs && \
    chown -R poker:poker /app

# Switch to non-root user
USER poker

# Expose ports
EXPOSE 3000 3001 5173

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3001/health || exit 1

# Run the application
CMD ["./ggpoker"]
