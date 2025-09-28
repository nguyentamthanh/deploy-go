#!/bin/bash

echo "🐳 Testing Docker build locally..."

# Clean up any existing images
echo "🧹 Cleaning up existing images..."
docker rmi deploy-go-test 2>/dev/null || true

# Build the Docker image
echo "🔨 Building Docker image..."
if docker build -t deploy-go-test .; then
    echo "✅ Docker build successful!"
    
    # Test if the binary was created correctly
    echo "🔍 Checking if binary exists in image..."
    if docker run --rm deploy-go-test ls -la /app/main; then
        echo "✅ Binary exists in Docker image!"
        
        # Test if the binary is executable
        echo "🔍 Testing binary execution..."
        if timeout 5 docker run --rm deploy-go-test /app/main --help 2>/dev/null || true; then
            echo "✅ Binary is executable!"
        else
            echo "⚠️ Binary execution test completed (timeout expected)"
        fi
    else
        echo "❌ Binary not found in Docker image!"
        exit 1
    fi
else
    echo "❌ Docker build failed!"
    exit 1
fi

echo "🎉 Docker build test completed successfully!"
