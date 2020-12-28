#!/usr/bin/env bash

source $(dirname "$0")/base.sh

send_gql_query '{"query":"{things {id, name, org {id, name} availability_topic}}"}'
