#!/bin/bash -e

curl -X POST http://localhost:9096/query -d '{"query":"{user (email: \"test@test.com\" ) {email, password}}"}' -H "$(cat headers.curl)" -i
