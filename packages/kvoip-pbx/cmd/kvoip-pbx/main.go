package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kvoip/kvoip-pbx/internal/api"
	"github.com/kvoip/kvoip-pbx/internal/cdr"
	"github.com/kvoip/kvoip-pbx/internal/config"
	"github.com/kvoip/kvoip-pbx/internal/dialog"
	"github.com/kvoip/kvoip-pbx/internal/handlers"
	"github.com/kvoip/kvoip-pbx/internal/media"
	"github.com/kvoip/kvoip-pbx/internal/proxy"
	"github.com/kvoip/kvoip-pbx/internal/server"
	"github.com/kvoip/kvoip-pbx/internal/session"
	"github.com/kvoip/kvoip-pbx/pkg/version"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config inválida", "err", err)
		os.Exit(1)
	}

	logger := newLogger(cfg.LogLevel)
	logger.Info("iniciando PBX",
		"service", cfg.ServiceName,
		"version", version.Version,
		"listen", cfg.ListenAddr(),
		"advertised", cfg.AdvertisedAddr(),
		"http", cfg.HTTPListenAddr(),
		"auth", cfg.AuthEnabled,
		"realm", cfg.AuthRealm,
		"users", len(cfg.SIPUsers),
		"media", cfg.MediaEnabled,
		"rtp", fmt.Sprintf("%d-%d", cfg.RTPPortMin, cfg.RTPPortMax),
		"cdr", cfg.CdrWebhookURL != "",
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	router := proxy.NewRouter()
	sessions := session.NewManager()
	dialogs := dialog.NewManager()
	var mediaEngine *media.Engine
	if cfg.MediaEnabled {
		mediaEngine = media.NewEngine(
			logger,
			cfg.MediaAdvertiseHost,
			cfg.MediaBindHost,
			cfg.RTPPortMin,
			cfg.RTPPortMax,
		)
	}
	cdrNotifier := cdr.NewNotifier(cfg.CdrWebhookURL, cfg.CdrWebhookSecret, logger)
	dispatcher := handlers.NewDispatcher(logger, router, sessions, dialogs, cfg, mediaEngine, cdrNotifier)
	sipServer := server.New(cfg, logger, dispatcher)

	httpAPI := api.New(cfg, logger, router, sessions, dispatcher.Digest())
	httpServer := &http.Server{
		Addr:              cfg.HTTPListenAddr(),
		Handler:           httpAPI.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("API HTTP do PBX iniciada", "addr", cfg.HTTPListenAddr())
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("API HTTP encerrou com erro", "err", err)
			stop()
		}
	}()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(shutdownCtx)
	}()

	if err := sipServer.ListenAndServeUDP(ctx); err != nil {
		logger.Error("servidor SIP encerrou com erro", "err", err)
		os.Exit(1)
	}

	logger.Info("PBX finalizado")
}

func newLogger(level string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	return slog.New(handler)
}
