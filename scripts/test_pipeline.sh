#!/bin/bash

# Mock upload
echo "Testing File Upload..."
curl -X POST http://localhost:8080/api/v1/upload \
  -F "file=@README.md" \
  -F "type=document"

echo -e "\n\nTesting Search..."
curl -X GET "http://localhost:8080/api/v1/search?q=architecture"
