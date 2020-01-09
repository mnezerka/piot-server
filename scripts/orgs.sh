#!/bin/bash -e

curl -X POST http://localhost:9096/query -d '{"query":"{orgs {id, name, created}}"}' -H "$(cat headers.curl)" | jq
