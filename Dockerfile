# Stage 1: Build Go binary
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY main.go ./

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server .

# Stage 2: Scratch image (smallest possible)
FROM scratch

WORKDIR /app

# Copy CA certificates for HTTPS (needed for Resend API)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary
COPY --from=builder /app/server .

# Copy static files
COPY index.html .
COPY style.css .
COPY favicon.svg .

# Expose port
EXPOSE 4001

# Run the binary
ENTRYPOINT ["/app/server"]