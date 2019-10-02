#!/bin/bash -e

curl -s -X POST http://localhost:9096/query -d '{"query":"{things {id, name, org {id, name}}}"}' -H "$(cat headers.curl)" | jq
