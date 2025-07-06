#
# ────────────────────────────────────────────────────────────────
#  STAGE 1 – Extract argocd bianry
# ────────────────────────────────────────────────────────────────
#
FROM quay.io/argoproj/argocd:v3.0.4 AS argocd-base

# Sanity check: the binary exists where expected
RUN test -f /usr/local/bin/argocd-cmp-server


#
# ────────────────────────────────────────────────────────────────
#  STAGE 2 – Compile the “sinaloa” CLI binary
# ────────────────────────────────────────────────────────────────
#
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


#
# ────────────────────────────────────────────────────────────────
#  STAGE 3 – Final image with Ubuntu base
# ────────────────────────────────────────────────────────────────
#
FROM ubuntu:24.04

# Add CA certificates if needed (HTTPS, etc.)
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates curl jq yq make bash git && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary from the build stage and the deploy scripts
COPY --from=builder /app/build/sinaloa /usr/local/bin/sinaloa
COPY --from=builder /app/scripts/ci-cd /scripts/ci-cd

# Make sure it's executable
RUN chmod +x /usr/local/bin/sinaloa

# Make all .sh scripts inside /scripts/ci-cd (recursively) executable
RUN find /scripts/ci-cd -type f -name "*.sh" -exec chmod +x {} \;

# Copy argocd-cmp-server binary from the base image at the correct path
COPY --from=argocd-base /usr/local/bin/argocd-cmp-server /usr/local/bin/argocd-cmp-server
RUN chmod +x /usr/local/bin/argocd-cmp-server

# Symlink for compatibility
RUN mkdir -p /var/run/argocd && ln -s /usr/local/bin/argocd-cmp-server /var/run/argocd/argocd-cmp-server

CMD ["/usr/local/bin/argocd-cmp-server"]
