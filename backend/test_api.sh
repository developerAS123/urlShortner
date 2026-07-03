#!/bin/bash
sleep 2

echo "Registering user..."
response=$(curl -s -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}')
echo $response
token=$(echo $response | grep -o '"token":"[^"]*' | grep -o '[^"]*$')

echo -e "\nShortening URL..."
short_response=$(curl -s -X POST http://localhost:8080/api/shorten \
  -H "Authorization: Bearer $token" \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://google.com"}')
echo $short_response
slug=$(echo $short_response | grep -o '"slug":"[^"]*' | grep -o '[^"]*$')

echo -e "\nTesting Redirect for slug: $slug"
curl -s -I http://localhost:8080/$slug | grep -i "location"
