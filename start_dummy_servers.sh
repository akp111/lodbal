#!/bin/bash

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🚀 Starting both backend servers..."

# Start Server1
echo "Starting Server 1 on port 3000..."
(cd "$SCRIPT_DIR/dummy_servers/server1" && go run main.go) &
SERVER1_PID=$!

# Small delay to ensure Server1 starts first
sleep 1

# Start Server2
echo "Starting Server 2 on port 3001..."
(cd "$SCRIPT_DIR/dummy_servers/server2" && go run main.go) &
SERVER2_PID=$!

echo "✅ Both servers started!"
echo "Server 1 PID: $SERVER1_PID - http://localhost:3000"
echo "Server 2 PID: $SERVER2_PID - http://localhost:3001"
echo ""
echo "📝 Test the servers:"
echo "   curl http://localhost:3000/"
echo "   curl http://localhost:3001/"
echo ""
echo "Press Ctrl+C to stop both servers"

# Function to cleanup when script exits
cleanup() {
    echo ""
    echo "🛑 Stopping servers..."
    kill $SERVER1_PID $SERVER2_PID 2>/dev/null
    echo "✅ Both servers stopped"
    exit 0
}

# Set trap to cleanup on Ctrl+C
trap cleanup SIGINT

# Wait for servers to finish (they won't unless killed)
wait $SERVER1_PID $SERVER2_PID
