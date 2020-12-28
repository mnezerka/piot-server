#!/usr/bin/env bash

source $(dirname "$0")/base.sh

# get user ID
echo Getting user id
send_gql_query_raw '{"query":"{users {id, email}}"}' | jq -r '.data.users[] | select (.email | contains("test")) | .id' > user_id.tmp
user_id=$(<user_id.tmp)
echo User id: $user_id

Q="{\"query\": \"mutation {updateUser(user: {id: \\\"${user_id}\\\", is_admin: true}) {id, is_admin} }\"}"

echo $Q

send_gql_query "$Q"
