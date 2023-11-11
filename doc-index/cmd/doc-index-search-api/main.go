package main

import (
	"context"

	"github.com/getnexar/golang-programming-task/doc-index/pkg/server"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(Module, fx.Invoke(
		registerLifecycleHooks,
	)).Run()
}

func registerLifecycleHooks(lifecycle fx.Lifecycle, logger *zap.SugaredLogger, server *server.HttpServer) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return server.Start()
		},
		OnStop: func(ctx context.Context) error {
			server.Stop()
			logger.Sync()
			return nil
		},
	})
}
