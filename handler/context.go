package handler

import (
    "context"
    "net/http"
    "github.com/op/go-logging"
)

type ContextHandler struct {
    log *logging.Logger
    ctx context.Context
    handler http.Handler
}

func NewContextHandler(log *logging.Logger, ctx context.Context, handler http.Handler) *ContextHandler {
    return &ContextHandler{log: log, ctx: ctx, handler: handler}
}

func (h *ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.handler.ServeHTTP(w, r.WithContext(h.ctx))
}
