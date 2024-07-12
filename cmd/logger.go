package main

import (
	"context"
	"log"
	"os"
	"time"

	"corrigan.io/go_api_seed/internal/config"
	logWriter "github.com/newrelic/go-agent/v3/integrations/logcontext-v2/zerologWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
)

func getCorrelationIdFromContext(ctx context.Context) string {
	correlationId, ok := ctx.Value(correlationIDKey).(string)
	if !ok {
		return ""
	}
	return correlationId
}

func getSessionTokenFromContext(ctx context.Context) string {
	sessionToken, ok := ctx.Value(sessionTokenKey).(string)
	if !ok {
		return ""
	}
	return sessionToken
}

func getIPFromContext(ctx context.Context) string {
	ip, ok := ctx.Value(ipKey).(string)
	if !ok {
		return ""
	}
	return ip
}

type TracingHook struct{}

func (h TracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	correlationId := getCorrelationIdFromContext(ctx) // as per your tracing framework
	if correlationId != "" {
		e.Str("correlation_id", correlationId)
	}
	ip := getIPFromContext(ctx)
	if ip != "" {
		e.Str("ip_address", ip)
	}
	sessionToken := getSessionTokenFromContext(ctx)
	if sessionToken != "" {
		e.Str("session_token", sessionToken)
	}
}

func getLogger(cfg *config.Config) zerolog.Logger {

	if cfg.Env == "local" {
		logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			With().
			Timestamp().
			Caller().
			Stack().
			Logger().
			Hook(TracingHook{})

		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		return logger
	} else {
		newRelicApp, err := newrelic.NewApplication(
			newrelic.ConfigAppName(cfg.Logging.NewRelicAppName),
			newrelic.ConfigLicense(cfg.Logging.NewRelicLicenseKey),
			newrelic.ConfigDistributedTracerEnabled(true),
			newrelic.ConfigAppLogForwardingEnabled(true),
			newrelic.ConfigAppLogDecoratingEnabled(true),
			func(cfg *newrelic.Config) {
				cfg.CustomInsightsEvents.Enabled = true
				cfg.DistributedTracer.Enabled = true
			},
		)
		if err != nil {
			log.Fatalln("unable to create New Relic Application", err)
		}

		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

		writer := logWriter.New(os.Stdout, newRelicApp)

		logger := zerolog.New(writer).
			With().
			Timestamp().
			Caller().
			Stack().
			Logger().
			Hook(TracingHook{})

		return logger
	}
}
