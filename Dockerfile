# Stage 1: Build the Go binary
FROM golang:tip-20250620-alpine3.21 AS builder

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

##################################################

# Stage 2: Final ubuntu image
FROM ubuntu:24.04

# Add CA certificates if needed (HTTPS, etc.)
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary from the build stage and the deploy scripts
COPY --from=builder /app/build/sinaloa /usr/local/bin/sinaloa
COPY --from=builder /app/scripts/ci-cd /scripts/ci-cd

# Make sure it's executable
RUN chmod +x /usr/local/bin/sinaloa

# Make all .sh scripts inside /scripts/ci-cd (recursively) executable
RUN find /scripts/ci-cd -type f -name "*.sh" -exec chmod +x {} \;

# Execute cmd or execute a cmd passed by arg
CMD ["sinaloa", "--help"]
