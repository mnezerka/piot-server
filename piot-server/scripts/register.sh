#!/bin/bash -e

curl -X POST http://localhost:9096/register -v -s -d '{"email":"test@test.com", "password": "test"}'
