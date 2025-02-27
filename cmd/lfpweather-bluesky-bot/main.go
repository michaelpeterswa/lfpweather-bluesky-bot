package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"alpineworks.io/ootel"
	"github.com/michaelpeterswa/lfpweather-bluesky-bot/internal/bluesky"
	"github.com/michaelpeterswa/lfpweather-bluesky-bot/internal/config"
	"github.com/michaelpeterswa/lfpweather-bluesky-bot/internal/forecast"
	"github.com/michaelpeterswa/lfpweather-bluesky-bot/internal/logging"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "error"
	}

	slogLevel, err := logging.LogLevelToSlogLevel(logLevel)
	if err != nil {
		log.Fatalf("could not convert log level: %s", err)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
	})))
	c, err := config.NewConfig()
	if err != nil {
		slog.Error("could not create config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx := context.Background()

	exporterType := ootel.ExporterTypePrometheus
	if c.Local {
		exporterType = ootel.ExporterTypeOTLPGRPC
	}

	ootelClient := ootel.NewOotelClient(
		ootel.WithMetricConfig(
			ootel.NewMetricConfig(
				c.MetricsEnabled,
				exporterType,
				c.MetricsPort,
			),
		),
		ootel.WithTraceConfig(
			ootel.NewTraceConfig(
				c.TracingEnabled,
				c.TracingSampleRate,
				c.TracingService,
				c.TracingVersion,
			),
		),
	)

	shutdown, err := ootelClient.Init(ctx)
	if err != nil {
		slog.Error("could not create ootel client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(5 * time.Second))
	if err != nil {
		slog.Error("could not create runtime metrics", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = host.Start()
	if err != nil {
		slog.Error("could not create host metrics", slog.String("error", err.Error()))
		os.Exit(1)
	}

	forecastClient := forecast.NewForecastClient(http.DefaultClient, c.LFPWeatherForecastInferenceAPIURL, c.LFPWeatherForecastInferenceAPIAPIKey)

	blueskyClient, err := bluesky.NewBlueskyClient(ctx, c.BlueskyHost, c.BlueskyUsername, c.BlueskyPassword)
	if err != nil {
		slog.Error("could not create bluesky client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer func() {
		_ = shutdown(ctx)
	}()

	cronDaemon := cron.New()
	_, err = cronDaemon.AddFunc(c.DailyForecastCronSchedule, func() {
		summary, err := forecastClient.GetSummary(ctx)
		if err != nil {
			slog.Error("could not get forecast summary", slog.String("error", err.Error()))
			return
		}

		postCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err = blueskyClient.WritePost(postCtx, summary.Summary)
		if err != nil {
			slog.Error("could not write post", slog.String("error", err.Error()))
			return
		}

		slog.Debug("wrote post", slog.String("summary", summary.Summary))
	})
	if err != nil {
		slog.Error("could not add cron job", slog.String("error", err.Error()))
	}

	cronDaemon.Start()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		sig := <-sigs

		slog.Info("received signal", slog.String("signal", sig.String()))
		cronDaemon.Stop()

		done <- true
	}()

	slog.Info("lfpweather-bluesky-bot is running")
	<-done
	slog.Info("lfpweather-bluesky-bot is shutting down")

}
