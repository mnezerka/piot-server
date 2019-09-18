#!/bin/bash -e

curl -X POST -i -H "$(cat headers.curl)" --data '
{ "query": "mutation {createCustomer(name: \"cust2\", description: \"cust2 desc\") {name}}" }
' http://localhost:9096/query
