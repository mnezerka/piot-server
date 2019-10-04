#!/bin/bash -e

curl -X POST http://localhost:9096/query -d '{"query":"{users {email orgs {name}}}"}' -H "$(cat headers.curl)" | jq
