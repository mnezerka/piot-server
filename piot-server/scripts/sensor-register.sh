#!/bin/bash -e

# Register Device
echo Registring device "script"
curl -s -X POST --data '
{
   "device": "script"
   "readings": [
       {
           "address": "sensor1",
           "t": 34
       }
}
' http://localhost:9096/adapter
