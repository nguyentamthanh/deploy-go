#!/bin/bash

echo "🧪 Testing different build approaches..."

# Test 1: Current Dockerfile
echo "📦 Testing current Dockerfile..."
if docker build -f Dockerfile -t deploy-go-current . 2>/dev/null; then
    echo "✅ Current Dockerfile: SUCCESS"
else
    echo "❌ Current Dockerfile: FAILED"
fi

# Test 2: Alternative Dockerfile
echo "📦 Testing alternative Dockerfile..."
if docker build -f Dockerfile.alternative -t deploy-go-alternative . 2>/dev/null; then
    echo "✅ Alternative Dockerfile: SUCCESS"
else
    echo "❌ Alternative Dockerfile: FAILED"
fi

# Test 3: Simple Dockerfile
echo "📦 Testing simple Dockerfile..."
if docker build -f Dockerfile.simple -t deploy-go-simple . 2>/dev/null; then
    echo "✅ Simple Dockerfile: SUCCESS"
else
    echo "❌ Simple Dockerfile: FAILED"
fi

echo "🎉 Build tests completed!"
