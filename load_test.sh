#!/bin/bash

# Analytics Load Test Runner
# This script runs the load test with different configurations

echo "🚀 Analytics Load Test Runner"
echo "============================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Default values
TOTAL_EVENTS=${1:-1000000}
BATCH_SIZE=${2:-1000}
NUM_USERS=${3:-10000}

echo "📊 Configuration:"
echo "   Total Events: $TOTAL_EVENTS"
echo "   Batch Size: $BATCH_SIZE"
echo "   Number of Users: $NUM_USERS"
echo ""

# Check if the server is running
if ! curl -s http://localhost:8080/api/health > /dev/null; then
    echo "⚠️  Warning: Analytics server is not running on localhost:8080"
    echo "   Start the server with: go run main.go"
    echo ""
fi

echo "🔧 Running load test from loadtest directory..."
cd loadtest

# Run the load test
echo "🏃 Starting load test..."
go run main.go $TOTAL_EVENTS $BATCH_SIZE $NUM_USERS

if [ $? -ne 0 ]; then
    echo "❌ Load test failed"
    exit 1
fi

cd ..

echo ""
echo "✅ Load test complete!"
echo ""
echo "📊 Quick stats commands:"
echo "   curl http://localhost:8080/api/stats"
echo "   curl http://localhost:8080/api/debug/events"
echo ""
echo "💾 Database info:"
echo "   Database size: $(du -sh analytics.db* 2>/dev/null | cut -f1 || echo 'No database files found')"
echo ""
echo "💡 Pro tips:"
echo "   - For 5M events: ./load_test.sh 5000000 2000 20000"
echo "   - For 10M events: ./load_test.sh 10000000 5000 50000"
echo "   - Monitor with: watch -n 1 'curl -s http://localhost:8080/api/stats | jq .total_events'"