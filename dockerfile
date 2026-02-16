# Use the version requested by your go.mod
FROM golang:1.25-alpine AS builder

# This ARG allows us to customize the build for each service
ARG SERVICE_NAME

WORKDIR /app

# Copy the specific service folder content into the build context
COPY . .

# Ensure we are in the correct directory if the context is the service folder
# then initialize, tidy, and build
RUN go mod init ${SERVICE_NAME} || true && \
    go mod tidy && \
    go build -o /app/main .

# Final Stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

# Expose the standard port (you'll map these differently in compose)
EXPOSE 8080

CMD ["./main"]