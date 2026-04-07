#!/bin/bash
# Local Kubernetes Deployment Script using Kind
# This script allows testing Kubernetes deployment locally

set -e

echo "=== Local Kubernetes Deployment (Kind) ==="

# Check if kind is installed
if ! command -v kind &> /dev/null; then
    echo "kind not found. Installing..."
    brew install kind
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "kubectl not found. Installing..."
    brew install kubectl
fi

# Create kind cluster
echo "Creating kind cluster..."
kind get clusters | grep video-streaming && echo "Cluster already exists" || kind create cluster --name video-streaming

# Build Docker images
echo "Building Docker images..."
docker build -t auth-service:latest -f dockerfile auth-service/
docker build -t file-service:latest -f dockerfile file-service/
docker build -t upload-service:latest -f dockerfile upload-service/
docker build -t streaming-service:latest -f dockerfile streaming-service/
docker build -t upload-app:latest upload-app/
docker build -t streaming-app:latest streaming-app/

# Load images into kind cluster
echo "Loading images into kind cluster..."
kind load docker-image auth-service:latest --name video-streaming
kind load docker-image file-service:latest --name video-streaming
kind load docker-image upload-service:latest --name video-streaming
kind load docker-image streaming-service:latest --name video-streaming
kind load docker-image upload-app:latest --name video-streaming
kind load docker-image streaming-app:latest --name video-streaming

# Install metrics-server for HPA
echo "Installing metrics-server..."
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Install ingress-nginx
echo "Installing ingress-nginx..."
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.4/deploy/static/provider/cloud/deploy.yaml

# Apply Kubernetes manifests
echo "Applying Kubernetes manifests..."
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/mysql-deployment.yaml
kubectl apply -f k8s/auth-service-deployment.yaml
kubectl apply -f k8s/file-service-deployment.yaml
kubectl apply -f k8s/upload-service-deployment.yaml
kubectl apply -f k8s/streaming-service-deployment.yaml
kubectl apply -f k8s/upload-app-deployment.yaml
kubectl apply -f k8s/streaming-app-deployment.yaml
kubectl apply -f k8s/ingress.yaml
kubectl apply -f k8s/hpa.yaml

echo ""
echo "=== Deployment Complete! ==="
echo ""
echo "Check pod status:"
kubectl get pods -n video-streaming
echo ""
echo "Services:"
kubectl get svc -n video-streaming
echo ""
echo "Ingress:"
kubectl get ingress -n video-streaming
echo ""
echo "HPA:"
kubectl get hpa -n video-streaming
echo ""
echo "To test the application, add the following to your /etc/hosts:"
echo "127.0.0.1 video-streaming.local"
echo ""
echo "Then access the application at:"
echo "  - Upload app: http://video-streaming.local/upload"
echo "  - Streaming app: http://video-streaming.local/streaming"