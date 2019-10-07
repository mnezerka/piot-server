#!/bin/bash -e

curl -X POST http://$PIOT_SERVER_HOSTNAME:9096/login -s -d '{"email":"test@test.com", "password": "test"}' | jq -r ".token" > token

echo "Authorization: Bearer `cat token`" > headers.curl
