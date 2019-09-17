#!/bin/bash -e

curl -X POST http://localhost:9097 -d '{"device":"TEST", "readings": [{"address": "ADDR"}]}' -i
