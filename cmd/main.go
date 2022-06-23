package main

import (
	"context"
	"fmt"
	"github.com/d7561985/tel"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
	"net/http"
	"os/signal"
	"projects/practic_5course_cesar/cmd/crypto"
	"projects/practic_5course_cesar/cmd/crypto/handler"
	"projects/practic_5course_cesar/internal/cryptosvc"
	"projects/practic_5course_cesar/internal/storage"
	"projects/practic_5course_cesar/pkg/custom"
	"syscall"
	"time"
)

func main() {
	cfg, err := crypto.ReadConfig()
	if err != nil {
		panic(fmt.Errorf("read config: %w", err))
	}

	telemetry := tel.New(cfg.Tel)
	defer telemetry.Close()

	ctx, stop := signal.NotifyContext(telemetry.Ctx(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var (
		enCrypt = new(custom.English).New()
		ukCrypt = new(custom.Ukrainian).New()
		cache   = storage.New(&cfg.Storage)
	)

	service := cryptosvc.New(enCrypt, ukCrypt, cache)

	var errCh = make(chan error)

	restHandlers := handler.New(
		service,
		schema.NewDecoder(),
		telemetry.Named("rest_handler"),
	)

	restRouter := handler.MakeRouter(ctx, restHandlers)

	httpServer := createHTTPServer(cfg.Listen, restRouter)

	go func() {
		telemetry.Info("listen and serve", zap.String("address", cfg.Listen))

		if err = httpServer.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	shutdown := func() {
		stop()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //nolint:gomnd,govet
		defer cancel()

		if err = httpServer.Shutdown(ctx); err != nil {
			telemetry.Error("http server: shutdown", zap.Error(err))

			return
		}

		telemetry.Info("service shutdown: graceful!", zap.Error(err))
	}

	select {
	case err := <-errCh:
		telemetry.Error("shutdown catch error", zap.Error(err))
		shutdown()
	case <-ctx.Done():
		telemetry.Info("shutdown context done")
		shutdown()
	}
}

const seconds10 = 10 * time.Second

func createHTTPServer(port string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:           port,
		Handler:        handler,
		ReadTimeout:    seconds10,
		WriteTimeout:   seconds10,
		IdleTimeout:    seconds10,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}
}
