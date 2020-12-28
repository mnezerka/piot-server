#!/usr/bin/env bash

THING_NAME=${1:-thing1}

source $(dirname "$0")/base.sh

echo Creating new thing with name=${THING_NAME}
send_gql_query "{ \"query\": \"mutation {createThing(name: \\\"${THING_NAME}\\\", type: \\\"device\\\") {id, name}}\" }"

# Get device ID
echo Getting ID of the thing named as "${THING_NAME}"
send_gql_query_raw '{"query":"{things (all: true) {id, name}}"}' | \
    jq -r ".data.things[] | select (.name | contains(\"${THING_NAME}\")) | .id"  > thing_id.tmp
thing_id=$(<thing_id.tmp)
echo thing id: $thing_id

echo "Setting thing attributes (enabled)"
Q="{\"query\": \"mutation {updateThing(thing: {id: \\\"${thing_id}\\\", enabled: true, last_seen_interval: 10}) {id, enabled} }\"}"
echo $Q
send_gql_query "$Q"
