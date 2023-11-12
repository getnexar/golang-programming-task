package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/getnexar/golang-programming-task/doc-index/pkg/config"
	"github.com/getnexar/golang-programming-task/doc-index/pkg/index"
	"github.com/getnexar/golang-programming-task/doc-index/pkg/server"
	"github.com/getnexar/golang-programming-task/doc-index/pkg/server/handlers"
)

var (
	Module fx.Option
)

func init() {
	Module = fx.Provide(
		provideConfig,
		provideLogger,
		provideIndex,
		provideHandlers,
		provideHttpRouter,
		provideHttpServer,
	)
}

func provideConfig() (*config.Config, error) {
	return config.Load()
}

func provideLogger(config *config.Config) (*zap.SugaredLogger, error) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	switch loglevel := strings.ToLower(config.LogLevel); loglevel {
	case "debug":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "error":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}

func provideIndex(logger *zap.SugaredLogger, config *config.Config) (*index.Index, error) {
	return index.NewIndex(config, logger)
}

func provideHandlers(logger *zap.SugaredLogger, index *index.Index) handlers.HandlersInterface {
	return handlers.NewHandlers(logger, index)
}

func provideHttpRouter(handlers handlers.HandlersInterface) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/healthcheck", handlers.Healthcheck)
	r.Get("/search", handlers.Search)
	r.Post("/search", handlers.Search)
	r.Delete("/document", handlers.DeleteDocument)
	return r
}

func provideHttpServer(
	handler http.Handler,
	config *config.Config,
) (*server.HttpServer, error) {
	return server.NewHttpServer(
		config.ServerAddress,
		handler,
	), nil
}
