#!/bin/bash

echo "======================================"
echo " Mock SSO Flow - Starting Server"
echo "======================================"
echo ""

# Kill any existing server on port 8080
lsof -ti:8080 | xargs kill -9 2>/dev/null || true
sleep 1

# Change to project directory
cd "/media/OliveGreen/Authentication/Assignment four"

# Build the server
echo "Building server..."
go build -o sso-server server.go

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build successful!"
echo ""
echo "Starting server on http://localhost:8080"
echo "Press Ctrl+C to stop"
echo ""

# Run the server
./sso-server
