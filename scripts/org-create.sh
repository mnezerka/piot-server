#!/bin/bash -e

curl -X POST -H "$(cat headers.curl)" --data '
{ "query": "mutation {createOrg(name: \"test\", description: \"test organization\") {id, name}}" }
' http://localhost:9096/query | jq
