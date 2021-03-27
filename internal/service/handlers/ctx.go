package handlers

import (
	"context"
	"github.com/olegfomenko/game-service/internal/horizon"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	conCtxKey ctxKey = iota
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxConnector(con *horizon.Connector) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, conCtxKey, con)
	}
}

func Connector(r *http.Request) *horizon.Connector {
	return r.Context().Value(conCtxKey).(*horizon.Connector)
}
