# Stage 1: Build the Go binary
FROM golang:1.24 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source code (including Makefile and .go files)
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o sinaloa ./src/main.go
RUN mkdir -p build && mv sinaloa build/

# Stage 2: Final lightweight image
FROM alpine:3.19

# Add CA certificates if needed (HTTPS, etc.)
RUN apk add --no-cache ca-certificates

# Copy the binary from the build stage
COPY --from=builder /app/build/sinaloa /usr/local/bin/sinaloa

# Make sure it's executable
RUN chmod +x /usr/local/bin/sinaloa

# Execute cmd or execute a cmd passed by arg
CMD ["sinaloa", "--help"]
