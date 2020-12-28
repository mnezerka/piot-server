#!/usr/bin/env bash

source $(dirname "$0")/base.sh

# get user ID
echo Getting user id
send_gql_query_raw '{"query":"{users {id, email}}"}' | jq -r '.data.users[] | select (.email | contains("test")) | .id' > user_id.tmp
user_id=$(<user_id.tmp)
echo User id: $user_id

echo Getting org id
send_gql_query_raw '{"query":"{orgs {id, name}}"}' | jq -r '.data.orgs[] | select (.name | contains("test")) | .id' > org_id.tmp
org_id=$(<org_id.tmp)
echo User id: $org_id

Q="{\"query\": \"mutation { addOrgUser(orgId: \\\"${org_id}\\\", userId: \\\"${user_id}\\\" ) {} }\"} "

send_gql_query "$Q"
