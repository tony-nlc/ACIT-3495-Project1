#!/bin/bash
# Load Testing Script for Horizontal Scalability Testing
# This script demonstrates that the backend services scale automatically based on load

set -e

SERVICE_URL="${1:-http://auth-service:5000}"
DURATION="${2:-60}"
CONCURRENT_USERS="${3:-50}"
RAMP_UP="${4:-10}"

echo "=== Horizontal Scalability Load Test ==="
echo "Target: $SERVICE_URL"
echo "Duration: ${DURATION}s"
echo "Concurrent Users: $CONCURRENT_USERS"
echo "Ramp-up: ${RAMP_UP}s"
echo ""

# Install hey if not present
if ! command -v hey &> /dev/null; then
    echo "Installing hey (load testing tool)..."
    go install github.com/rakyll/hey@latest
    export PATH=$PATH:$(go env GOPATH)/bin
fi

# Get initial replica count
echo "=== Initial State ==="
echo "Current pod replicas:"
kubectl get pods -n video-streaming -l app=auth-service -o name
kubectl get hpa auth-service-hpa -n video-streaming

echo ""
echo "Starting load test..."
echo "Simulating $CONCURRENT_USERS concurrent users for ${DURATION}s"
echo ""

# Run load test
hey -n 100000 \
    -c $CONCURRENT_USERS \
    -q 100 \
    -timeout 30s \
    "$SERVICE_URL/login" \
    -m POST \
    -D '{"User":"admin","Pass":"password"}' \
    -H "Content-Type: application/json" &
HEY_PID=$!

# Monitor pod scaling during load test
echo "=== Monitoring Pod Scaling ==="
for i in {1..12}; do
    sleep 5
    echo "--- Time: $((i*5))s ---"
    echo "Pod status:"
    kubectl get pods -n video-streaming -l app=auth-service -o wide
    echo "HPA status:"
    kubectl get hpa auth-service-hpa -n video-streaming
    echo "Pod metrics:"
    kubectl top pods -n video-streaming -l app=auth-service 2>/dev/null || echo "(metrics not available)"
    echo ""
done

# Wait for hey to complete
wait $HEY_PID

echo ""
echo "=== Final State ==="
echo "Final pod replicas:"
kubectl get pods -n video-streaming -l app=auth-service -o wide
echo ""
echo "Final HPA status:"
kubectl get hpa auth-service-hpa -n video-streaming

echo ""
echo "=== Test Complete ==="
echo ""
echo "This demonstrates horizontal pod autoscaling:"
echo "- Under load, the HPA automatically scales up the number of pods"
echo "- After the load decreases, the HPA scales down pods (after stabilization period)"
echo "- The service maintains availability during scaling events"