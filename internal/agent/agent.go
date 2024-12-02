package agent

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/alpineworks/katalog-agent/internal/kubernetes"
	agentclient "github.com/alpineworks/katalog/backend/pkg/agent"
	"github.com/alpineworks/katalog/backend/pkg/agentservice"
)

type Agent struct {
	kubernetesClient   *kubernetes.KubernetesClient
	agentServiceClient *agentclient.AgentServiceClient
}

func NewAgent(kc *kubernetes.KubernetesClient, asc *agentclient.AgentServiceClient) *Agent {
	return &Agent{
		kubernetesClient:   kc,
		agentServiceClient: asc,
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

	slog.Debug("namespaces", slog.String("namespaces", strings.Join(namespaces, ",")))

	var totalDeployments []kubernetes.Deployment
	for _, namespace := range namespaces {
		deployments, err := a.kubernetesClient.GetDeploymentsFromNamespace(collectionCtx, namespace)
		if err != nil {
			slog.Error("failed to get deployments", slog.String("error", err.Error()))
			return
		}

		totalDeployments = append(totalDeployments, deployments...)

		slog.Debug("deployments", slog.String("namespace", namespace), slog.Int("deployments", len(deployments)))
	}

	slog.Info("successfully collected data", slog.Int("namespaces", len(namespaces)), slog.Int("deployments", len(totalDeployments)))

	response, err := a.agentServiceClient.PublishDeployments(ctx, translateDeployments(totalDeployments))
	if err != nil {
		slog.Error("failed to publish deployments", slog.String("error", err.Error()))
		return
	}

	if !response.Success {
		if response.Error == nil {
			slog.Error("failed to publish deployments", slog.String("error", "unknown error"))
			return
		}

		slog.Error("failed to publish deployments", slog.String("error", *response.Error))
		return
	}

	slog.Info("successfully published deployments", slog.Int("deployments", len(totalDeployments)))
}

func translateDeployments(deployments []kubernetes.Deployment) *agentservice.PublishDeploymentsRequest {
	var translatedDeployments []*agentservice.Deployment
	for _, d := range deployments {
		var containers []*agentservice.Container
		for _, c := range d.Containers {
			containers = append(containers, &agentservice.Container{
				Name:  c.Name,
				Image: c.Image,
				Tag:   c.Tag,
			})

			slog.Debug("translated container", slog.Any("containers", containers))
		}

		translatedDeployments = append(translatedDeployments, &agentservice.Deployment{
			Name:         d.Name,
			Replicas:     d.Replicas,
			TrueReplicas: d.TrueReplicas,
			Labels:       d.Labels,
			Containers:   containers,
		})

		slog.Debug("translated deployments", slog.Any("deployments", translatedDeployments))
	}

	slog.Debug("final translated deployments", slog.Any("deployments", translatedDeployments))

	return &agentservice.PublishDeploymentsRequest{
		Deployments: translatedDeployments,
	}
}
