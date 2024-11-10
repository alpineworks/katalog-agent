package agent

import (
	"context"
	"log/slog"
	"time"

	"github.com/alpineworks/katalog-agent/internal/kubernetes"
)

type Agent struct {
	kubernetesClient *kubernetes.KubernetesClient
}

func NewAgent(kc *kubernetes.KubernetesClient) *Agent {
	return &Agent{
		kubernetesClient: kc,
	}
}

func (a *Agent) Collect() {
	ctx := context.Background()

	// limit the time we spend collecting data to 1 minute
	collectionCtx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	namespaces, err := a.kubernetesClient.GetNamespaces(collectionCtx)
	if err != nil {
		slog.Error("failed to get namespaces", slog.String("error", err.Error()))
		return
	}

	var deployments []kubernetes.Deployment
	for _, namespace := range namespaces {
		deployments, err = a.kubernetesClient.GetDeploymentsFromNamespace(collectionCtx, namespace)
		if err != nil {
			slog.Error("failed to get deployments", slog.String("error", err.Error()))
			return
		}
	}

	slog.Info("successfully collected data", slog.Int("namespaces", len(namespaces)), slog.Int("deployments", len(deployments)))

	// TODO: send data to the collector
}
