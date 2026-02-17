# 1. Build Stage
FROM golang:1.23-alpine AS builder

# This ARG comes from your docker-compose 'args' section
ARG SERVICE_NAME

WORKDIR /app

# Copy everything from the service's context
COPY . .

# If the service folder doesn't have a go.mod, this initializes it 
# using the SERVICE_NAME passed from Compose.
RUN if [ ! -f go.mod ]; then go mod init ${SERVICE_NAME}; fi && \
    go mod tidy && \
    go build -o /app/main .

# 2. Final Stage
FROM alpine:latest
WORKDIR /app/

# IMPORTANT: Create the storage directory that your volume mounts to
# and ensure it is writable by the application.
RUN mkdir -p /storage && chmod 777 /storage

# Copy the binary from the builder
COPY --from=builder /app/main .

# Standard Gin Port (Mapped in Compose)
EXPOSE 8080

CMD ["./main"]