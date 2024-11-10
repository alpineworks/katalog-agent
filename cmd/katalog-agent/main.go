package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"

	"github.com/alpineworks/katalog-agent/internal/agent"
	"github.com/alpineworks/katalog-agent/internal/config"
	"github.com/alpineworks/katalog-agent/internal/kubernetes"
	"github.com/alpineworks/katalog-agent/internal/logging"
	"github.com/alpineworks/ootel"
)

func main() {
	slogHandler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(slogHandler))

	slog.Info("welcome to katalog-agent!")

	c, err := config.NewConfig()
	if err != nil {
		slog.Error("could not create config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slogLevel, err := logging.LogLevelToSlogLevel(c.LogLevel)
	if err != nil {
		slog.Error("could not parse log level", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.SetLogLoggerLevel(slogLevel)

	ctx := context.Background()

	ootelClient := ootel.NewOotelClient(
		ootel.WithMetricConfig(
			ootel.NewMetricConfig(
				c.MetricsEnabled,
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
		panic(err)
	}

	defer func() {
		_ = shutdown(ctx)
	}()

	kubernetesClient, err := kubernetes.NewKubernetesClient()
	if err != nil {
		slog.Error("could not create kubernetes client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	katalogAgent := agent.NewAgent(kubernetesClient)

	cronClient := cron.New()
	_, err = cronClient.AddFunc(c.CronSchedule, katalogAgent.Collect)
	if err != nil {
		slog.Error("could not add cron function", slog.String("error", err.Error()))
		os.Exit(1)
	}

	cronClient.Start()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		sig := <-sigs

		slog.Info("received signal", slog.String("signal", sig.String()))
		cronClient.Stop()

		done <- true
	}()

	slog.Info("katalog-agent is running")
	<-done
	slog.Info("katalog-agent is shutting down")

}
