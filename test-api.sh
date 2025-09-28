#!/bin/bash

echo "🧪 Testing API endpoints..."

# Test main endpoint
echo "📡 Testing main endpoint..."
curl -s http://localhost:8080/ | jq .

# Test health endpoint
echo "🏥 Testing health endpoint..."
curl -s http://localhost:8080/health | jq .

# Test hello endpoint
echo "👋 Testing hello endpoint..."
curl -s http://localhost:8080/api/hello | jq .

# Test users endpoint
echo "👥 Testing users endpoint..."
curl -s http://localhost:8080/api/users | jq .

# Test create user
echo "➕ Testing create user..."
curl -s -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Test User", "email": "test@example.com"}' | jq .

# Test posts endpoint
echo "📝 Testing posts endpoint..."
curl -s http://localhost:8080/api/posts | jq .

echo "✅ All tests completed!"
