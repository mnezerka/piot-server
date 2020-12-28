function send_gql_query() {
    send_gql_query_raw "$1" | jq
}

function send_gql_query_raw() {
    curl -s -X POST http://localhost:9096/query -d "$1" -H "$(cat headers.curl)"
}
