#!/bin/bash
set -e

echo "Building Docker images for Kubernetes deployment..."

# Build auth-service
echo "Building auth-service..."
docker build -t auth-service:latest -f dockerfile auth-service/

# Build file-service
echo "Building file-service..."
docker build -t file-service:latest -f dockerfile file-service/

# Build upload-service
echo "Building upload-service..."
docker build -t upload-service:latest -f dockerfile upload-service/

# Build streaming-service
echo "Building streaming-service..."
docker build -t streaming-service:latest -f dockerfile streaming-service/

# Build upload-app
echo "Building upload-app..."
docker build -t upload-app:latest upload-app/

# Build streaming-app
echo "Building streaming-app..."
docker build -t streaming-app:latest streaming-app/

echo "All images built successfully!"
echo "Tags:"
echo "  - auth-service:latest"
echo "  - file-service:latest"
echo "  - upload-service:latest"
echo "  - streaming-service:latest"
echo "  - upload-app:latest"
echo "  - streaming-app:latest"