# 1. Build Stage
FROM golang:1.23-alpine AS builder

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
FROM alpine:latest
# Re-add certificates so the app can make HTTPS calls to other services
RUN apk add --no-cache ca-certificates
WORKDIR /app/

# Create storage for the volume
RUN mkdir -p /storage && chmod 777 /storage

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]