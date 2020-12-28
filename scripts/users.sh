#!/usr/bin/env bash

source $(dirname "$0")/base.sh

send_gql_query '{"query":"{users {email orgs {name}}}"}'
