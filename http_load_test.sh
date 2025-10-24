#!/bin/bash

# HTTP Load Test Runner for Analytics API
# This script tests the /api/track endpoint with concurrent HTTP requests

echo "🔥 Analytics HTTP Load Test"
echo "============================"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Default values
TOTAL_REQUESTS=${1:-10000}
CONCURRENCY=${2:-10}
NUM_USERS=${3:-1000}
DURATION=${4:-300}
SERVER_URL=${5:-"http://localhost:8080"}

echo "📊 HTTP Load Test Configuration:"
echo "   Server URL: $SERVER_URL"
echo "   Total Requests: $TOTAL_REQUESTS"
echo "   Concurrency: $CONCURRENCY"
echo "   Users: $NUM_USERS"
echo "   Max Duration: ${DURATION}s"
echo ""

# Check if the server is running
echo "🔍 Checking server connectivity..."
if curl -s --max-time 5 "$SERVER_URL/api/health" > /dev/null; then
    echo "✅ Server is responding at $SERVER_URL"
else
    echo "❌ Cannot connect to server at $SERVER_URL"
    echo "   Make sure the analytics server is running: go run main.go"
    exit 1
fi

echo ""
echo "🏃 Starting HTTP load test..."
echo "   Press Ctrl+C to stop early"
echo ""

cd loadtest

# Run the HTTP load test
go run http_test.go $TOTAL_REQUESTS $CONCURRENCY $NUM_USERS $DURATION "$SERVER_URL"

if [ $? -ne 0 ]; then
    echo "❌ HTTP load test failed"
    exit 1
fi

cd ..

echo ""
echo "🎯 Load test complete! Check these endpoints:"
echo "   📊 Stats: curl $SERVER_URL/api/stats"
echo "   🔍 Recent events: curl $SERVER_URL/api/debug/events"
echo ""
echo "💡 Performance tips:"
echo "   - For high throughput: ./http_load_test.sh 100000 50 10000 120"
echo "   - For sustained load: ./http_load_test.sh 50000 20 5000 600"
echo "   - Monitor server: watch -n 1 'curl -s $SERVER_URL/api/stats | jq .total_events'"