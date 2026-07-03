#!/bin/bash
sleep 2

echo "Logging in as existing analytics_test user..."
response=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"analytics_test@example.com","password":"password123"}')
token=$(echo $response | grep -o '"token":"[^"]*' | grep -o '[^"]*$')

echo "Getting links for user..."
links_response=$(curl -s -X GET http://localhost:8080/api/links \
  -H "Authorization: Bearer $token")
slug=$(echo $links_response | grep -o '"slug":"[^"]*' | head -1 | grep -o '[^"]*$')

echo "Waiting for AI background worker to finish (10 seconds)..."
sleep 10

echo -e "\nFetching AI summary for slug: $slug..."
curl -s -X GET http://localhost:8080/api/links/$slug/summary \
  -H "Authorization: Bearer $token"
echo -e "\nDone."
