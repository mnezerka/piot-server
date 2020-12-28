#!/usr/bin/env bash

source $(dirname "$0")/base.sh

send_gql_query '
{ "query": "mutation {createOrg(name: \"test\", description: \"test organization\") {id, name}}" }
'
