# 1. Build Stage
FROM golang:1.25.1-alpine AS builder

# Required for fetching Go modules over HTTPS/Git
RUN apk add --no-cache git ca-certificates

ARG SERVICE_NAME
WORKDIR /app

# Copy the service code from the context
COPY . .

# Initialize and build
# We use 'go mod tidy' to fetch dependencies listed in your code
RUN if [ ! -f go.mod ]; then go mod init ${SERVICE_NAME}; fi && \
    go mod tidy && \
    go build -o /app/main .

# 2. Final Stage
# ... (builder stage above stays the same)

# Final Stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app/

# 1. Create the directory
# 2. Set permissions to be world-writable (safe for dev volumes)
RUN mkdir -p /storage && chmod 777 /storage

COPY --from=builder /app/main .

# Ensure we run as root to avoid permission friction with volumes in dev
USER root

EXPOSE 8080
CMD ["./main"]