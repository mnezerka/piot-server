#!/bin/bash -e


# Send sensort data
echo Sending sensor data
curl -s -X POST --data '
{
   "device": "script",
   "readings": [
       {
           "address": "sensor1",
           "t": 34
       }
   ]
}
' http://$PIOT_SERVER_HOSTNAME:9096/adapter


