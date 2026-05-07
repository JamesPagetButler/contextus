// Command contextus-synthesis runs the Synthesis subscriber against a NATS
// broker. It connects, subscribes to the three Spec v1.3 §5.3 subjects
// (ctx.edge.boundary.>, ctx.corpus.diversity, ctx.bridge.intervention.>),
// and dispatches incoming payloads into a synthesis.Agent that mints
// InsightSignals when findings cross the persistence-boundary threshold.
//
// The binary uses a NoopPersister until the Wyrd `contextus/` subpackage
// (Wyrd issue #6) ships a Wyrd-backed Persister. Until then, minted signals
// are observed only through structured logs.
//
// Configuration:
//
//	-nats-url   string   NATS broker URL (default "nats://127.0.0.1:4222")
//	-log-level  string   "debug" | "info" | "warn" | "error" (default "info")
//
// Environment overrides (read first; flags override):
//
//	CONTEXTUS_NATS_URL
//	CONTEXTUS_LOG_LEVEL
//
// Graceful shutdown: SIGINT or SIGTERM cancels the context, the service stops
// subscriptions, and the Subscriber drains in-flight messages before close.
//
// References:
//   - Spec v1.3 §4.4 Synthesis as Persistence Boundary
//   - Spec v1.3 §5.3 NATS Subject Hierarchy
//   - go-coding-guide.md (slog logging; context propagation; goroutine
//     termination plans)
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/JamesPagetButler/contextus/internal/contextus/natsadapter"
	"github.com/JamesPagetButler/contextus/internal/contextus/synthesis"
)

// Default values for configuration; can be overridden by env or flags.
const (
	defaultNATSURL  = "nats://127.0.0.1:4222"
	defaultLogLevel = "info"
	drainTimeout    = 30 * time.Second
)

func main() {
	if err := run(); err != nil {
		// Final error log at the top of the call stack per go-coding-guide.
		slog.Error("contextus-synthesis: fatal", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

func run() error {
	cfg := loadConfig()
	logger := newLogger(cfg.logLevel)
	slog.SetDefault(logger)
	logger.Info("contextus-synthesis: starting",
		slog.String("nats_url", cfg.natsURL),
		slog.String("log_level", cfg.logLevel),
	)

	// Top-level context: cancelled on SIGINT/SIGTERM. Goroutine termination
	// for both the NATS subscription handlers and the main loop.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	sub, err := natsadapter.Connect(ctx, cfg.natsURL)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	persister := &synthesis.NoopPersister{}
	logger.Warn("contextus-synthesis: using NoopPersister (Wyrd contextus/ subpackage not yet available)")

	agent := synthesis.New(persister)

	svc := natsadapter.New(sub, agent, logger)
	if err := svc.Start(ctx); err != nil {
		// Stop nothing — Start rolls back partial subscriptions on error.
		drainCtx, drainCancel := context.WithTimeout(context.Background(), drainTimeout)
		defer drainCancel()
		_ = sub.Drain(drainCtx)
		return fmt.Errorf("start: %w", err)
	}
	logger.Info("contextus-synthesis: subscribed to all three subjects; waiting for signals")

	<-ctx.Done()

	// Shutdown sequence:
	//   1. Stop subscriptions (no new messages dispatched into the agent).
	//   2. Drain the connection (in-flight messages finish; connection closes).
	logger.Info("contextus-synthesis: shutting down", slog.String("reason", ctxReason(ctx)))
	svc.Stop()

	// Use a fresh context for drain; the main ctx is already cancelled.
	drainCtx, drainCancel := context.WithTimeout(context.Background(), drainTimeout)
	defer drainCancel()
	if err := sub.Drain(drainCtx); err != nil {
		// Drain timeout or error — return so main can exit non-zero.
		return fmt.Errorf("drain: %w", err)
	}
	logger.Info("contextus-synthesis: shut down cleanly")
	return nil
}

type config struct {
	natsURL  string
	logLevel string
}

func loadConfig() config {
	cfg := config{
		natsURL:  envOrDefault("CONTEXTUS_NATS_URL", defaultNATSURL),
		logLevel: envOrDefault("CONTEXTUS_LOG_LEVEL", defaultLogLevel),
	}
	flag.StringVar(&cfg.natsURL, "nats-url", cfg.natsURL, "NATS broker URL")
	flag.StringVar(&cfg.logLevel, "log-level", cfg.logLevel, `log level: "debug" | "info" | "warn" | "error"`)
	flag.Parse()
	return cfg
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func newLogger(level string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		// Unknown level: default to info and warn about it.
		lvl = slog.LevelInfo
	}
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl})
	return slog.New(h)
}

// ctxReason returns a short reason string for ctx.Done(). Used in shutdown
// logging so operators can tell SIGINT from SIGTERM (or unexpected cancel).
func ctxReason(ctx context.Context) string {
	if err := ctx.Err(); errors.Is(err, context.Canceled) {
		return "signal"
	}
	return "timeout"
}
