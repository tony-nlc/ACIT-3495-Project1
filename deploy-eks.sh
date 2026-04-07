#!/bin/bash
set -e

# AWS EKS Deployment Script
# This script creates an EKS cluster and deploys the microservices

CLUSTER_NAME="video-streaming-cluster"
REGION="us-west-2"
NODE_GROUP_NAME="video-streaming-nodes"

echo "=== AWS EKS Deployment Script ==="
echo "Cluster: $CLUSTER_NAME"
echo "Region: $REGION"
echo ""

# Check if eksctl is installed
if ! command -v eksctl &> /dev/null; then
    echo "eksctl not found. Installing..."
    brew install eksctl
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "kubectl not found. Installing..."
    brew install kubectl
fi

# Check if AWS CLI is installed
if ! command -v aws &> /dev/null; then
    echo "aws CLI not found. Please install AWS CLI first."
    exit 1
fi

# Create EKS cluster
echo "Creating EKS cluster..."
eksctl create cluster \
    --name $CLUSTER_NAME \
    --region $REGION \
    --nodegroup-name $NODE_GROUP_NAME \
    --node-type t3.medium \
    --nodes-min 2 \
    --nodes-max 10 \
    --nodes 3 \
    --managed

echo "Cluster created successfully!"
echo ""

# Update kubeconfig
echo "Updating kubeconfig..."
aws eks update-kubeconfig --name $CLUSTER_NAME --region $REGION

# Install metrics-server for HPA
echo "Installing metrics-server..."
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Install ingress-nginx
echo "Installing ingress-nginx..."
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.4/deploy/static/provider/cloud/deploy.yaml

# Build and push Docker images to ECR
echo "Building and pushing images to ECR..."

# Get AWS account ID and region
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query 'Account' --output text)
ECR_REGISTRY="${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com"

# Create ECR repositories
for repo in auth-service file-service upload-service streaming-service upload-app streaming-app; do
    aws ecr create-repository --repository-name $repo --region $REGION || true
done

# Login to ECR
aws ecr get-login-password --region $REGION | docker login --username AWS --password-stin $ECR_REGISTRY

# Build and push images
for service in auth-service file-service upload-service streaming-service upload-app streaming-app; do
    echo "Building and pushing $service..."
    if [ "$service" = "upload-app" ] || [ "$service" = "streaming-app" ]; then
        docker build -t $service:latest $service/
    else
        docker build -t $service:latest -f dockerfile $service/
    fi

    docker tag $service:latest ${ECR_REGISTRY}/${service}:latest
    docker push ${ECR_REGISTRY}/${service}:latest
done

# Update image names in Kubernetes manifests
echo "Updating Kubernetes manifests with ECR image paths..."

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
echo "Check services:"
kubectl get svc -n video-streaming
echo ""
echo "Get ingress IP:"
kubectl get ingress -n video-streaming
echo ""
echo "Check HPA status:"
kubectl get hpa -n video-streaming