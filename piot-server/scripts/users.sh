#!/bin/bash -e

curl -X POST http://localhost:9096/query -d '{"query":"{users {email}}"}' -H "$(cat headers.curl)" | jq
