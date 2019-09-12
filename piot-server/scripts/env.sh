#!/bin/bash -e

MONGO_ADDR=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mongodb)

echo Setting mongodb address to $MONGO_ADDR

export MONGODB_URI=mongodb://$MONGO_ADDR:27017
