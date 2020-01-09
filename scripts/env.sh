#!/bin/bash -e

MONGO_ADDR=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' piot-server_mongodb_1)
echo Setting mongodb address to $MONGO_ADDR
export MONGODB_URI=mongodb://$MONGO_ADDR:27017

PIOT_SERVER_ADDR=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' piot_piot-server_1)
echo Setting piot server hostname to $PIOT_SERVER_ADDR
export PIOT_SERVER_HOSTNAME=$PIOT_SERVER_ADDR



