#!/bin/bash -e

curl -X POST http://$PIOT_SERVER_HOSTNAME:9096/register -v -s -d '{"email":"test@test.com", "password": "test"}'
