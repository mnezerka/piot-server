#!/bin/bash -e

# Register Device
echo Senting data packet to register device "script" ans sensor "sensor1"
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

# Get device ID
echo Getting ID of the script thing
curl -s -X POST http://$PIOT_SERVER_HOSTNAME:9096/query -d '{"query":"{things {id, name}}"}' -H "$(cat headers.curl)"
curl -s -X POST http://$PIOT_SERVER_HOSTNAME:9096/query -d '{"query":"{things {id, name}}"}' -H "$(cat headers.curl)" | jq -r '.data.things[] | select (.name | contains("script")) | .id'  > thing_id.tmp
thingId=$(<thing_id.tmp)
echo First thing id: $thingId


# Get sensor ID
echo Getting ID of the sensor thing
curl -s -X POST http://$PIOT_SERVER_HOSTNAME:9096/query -d '{"query":"{things {id, name}}"}' -H "$(cat headers.curl)" | jq -r '.data.things[] | select (.name | contains("sensor1")) | .id'  > sensor_id.tmp
sensorId=$(<sensor_id.tmp)
echo First sensor id: $sensorId


# Get org ID
echo Getting ID of the test org
curl -s -X POST http://$PIOT_SERVER_HOSTNAME:9096/query -d '{"query":"{orgs {id, name}}"}' -H "$(cat headers.curl)" | jq -r '.data.orgs[] | select (.name | contains("test")) | .id'  > org_id.tmp
#curl -s -X POST http://localhost:9096/query -d '{"query":"{orgs {id, name}}"}' -H "$(cat headers.curl)" | jq -r ".data.orgs[0].id" > org_id.tmp
orgId=$(<org_id.tmp)
echo First org id is  $orgId, assigning device $thingId to it


# Assign Device to Org
gql="{ \"query\": \"mutation {updateThing(thing: {id: \\\"$thingId\\\" orgId: \\\"$orgId\\\"}) {id}}\" }"
curl -s -X POST -H "$(cat headers.curl)" --data "$gql" http://$PIOT_SERVER_HOSTNAME:9096/query | jq

# Assign Sensor to Org
gql="{ \"query\": \"mutation {updateThing(thing: {id: \\\"$sensorId\\\" orgId: \\\"$orgId\\\"}) {id}}\" }"
curl -s -X POST -H "$(cat headers.curl)" --data "$gql" http://$PIOT_SERVER_HOSTNAME:9096/query | jq


