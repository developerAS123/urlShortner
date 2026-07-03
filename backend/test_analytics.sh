#!/bin/bash
sleep 2

echo "Registering user (might fail if already exists, that's fine)..."
response=$(curl -s -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"email":"analytics_test@example.com","password":"password123"}')

if echo "$response" | grep -q "Email already in use"; then
  echo "User exists, logging in..."
  response=$(curl -s -X POST http://localhost:8080/api/login \
    -H "Content-Type: application/json" \
    -d '{"email":"analytics_test@example.com","password":"password123"}')
fi
token=$(echo $response | grep -o '"token":"[^"]*' | grep -o '[^"]*$')

echo -e "\nShortening URL..."
short_response=$(curl -s -X POST http://localhost:8080/api/shorten \
  -H "Authorization: Bearer $token" \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://github.com"}')
slug=$(echo $short_response | grep -o '"slug":"[^"]*' | grep -o '[^"]*$')
echo "Slug: $slug"

echo -e "\nSimulating 3 clicks on the short URL..."
curl -s -A "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X)" http://localhost:8080/$slug > /dev/null
curl -s -A "Mozilla/5.0 (Windows NT 10.0; Win64; x64)" http://localhost:8080/$slug > /dev/null
curl -s -A "Mozilla/5.0 (Windows NT 10.0; Win64; x64)" http://localhost:8080/$slug > /dev/null

sleep 1 # Wait for async goroutines to insert data

echo -e "\nFetching all links for user..."
curl -s -X GET http://localhost:8080/api/links \
  -H "Authorization: Bearer $token"

echo -e "\n\nFetching analytics for slug: $slug..."
curl -s -X GET http://localhost:8080/api/links/$slug/analytics \
  -H "Authorization: Bearer $token"

echo -e "\n\nTesting Rate Limit (making 11 fast requests to /shorten)..."
for i in {1..11}; do
  curl -s -o /dev/null -w "%{http_code}\n" -X POST http://localhost:8080/api/shorten \
    -H "Authorization: Bearer $token" \
    -H "Content-Type: application/json" \
    -d '{"original_url":"https://example.com"}'
done
