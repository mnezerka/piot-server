#!/bin/bash -e

curl -X POST http://localhost:9096/query -d '{"query":"{userProfile {email}}"}' -H "$(cat headers.curl)" -i
