#!/bin/bash

url="http://localhost:3000/api/v1/create"

for i in {1..100}
do
  curl -X POST $url \
  -H "Content-Type: application/json" \
  -d "{\"id\": $i, \"title\": \"Demo Todo $i\", \"completed\": false}"
done