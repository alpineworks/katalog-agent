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
	agentclient "github.com/alpineworks/katalog/backend/pkg/agent"
	"github.com/alpineworks/ootel"
	"github.com/michaelpeterswa/go-mtls"
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
		slog.Error("could not initialize ootel client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer func() {
		_ = shutdown(ctx)
	}()

	var x509Files *mtls.X509Files
	if c.KatalogCertificateFile != "" && c.KatalogKeyFile != "" && c.KatalogCAFile != "" {
		x509Files = mtls.NewX509Files(c.KatalogCertificateFile, c.KatalogKeyFile, c.KatalogCAFile)
	}

	var agentServiceClient *agentclient.AgentServiceClient
	if x509Files != nil {
		agentServiceClient, err = agentclient.NewAgentServiceClient(c.KatalogHost, agentclient.WithMutualTLS(x509Files))
		if err != nil {
			slog.Error("could not create agent service client", slog.String("error", err.Error()))
			os.Exit(1)
		}
	} else {
		agentServiceClient, err = agentclient.NewAgentServiceClient(c.KatalogHost)
		if err != nil {
			slog.Error("could not create agent service client", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}

	kubernetesClient, err := kubernetes.NewKubernetesClient()
	if err != nil {
		slog.Error("could not create kubernetes client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	katalogAgent := agent.NewAgent(kubernetesClient, agentServiceClient)

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
