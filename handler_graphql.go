package main

import (
    "encoding/json"
    "github.com/graph-gophers/graphql-go"
    "net/http"
)

type GraphQLHandler struct {
    schema  *graphql.Schema
}

func NewGraphQLHandler(schema  *graphql.Schema) *GraphQLHandler {
    return &GraphQLHandler{schema: schema}
}

func (h *GraphQLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var params struct {
        Query         string                 `json:"query"`
        OperationName string                 `json:"operationName"`
        Variables     map[string]interface{} `json:"variables"`
    }

    if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // by mn ctx := h.Loaders.Attach(r.Context())
    ctx := r.Context()

    response := h.schema.Exec(ctx, params.Query, params.OperationName, params.Variables)

    responseJSON, err := json.Marshal(response)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(responseJSON)
}


