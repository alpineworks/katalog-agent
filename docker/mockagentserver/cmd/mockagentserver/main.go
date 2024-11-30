package main

import (
	"context"
	"log"
	"net"

	"github.com/alpineworks/katalog/backend/pkg/agentservice"
	"google.golang.org/grpc"
)

type AgentServiceServer struct {
	agentservice.UnimplementedAgentServiceServer
}

func newAgentServiceServer() *AgentServiceServer {
	return &AgentServiceServer{}
}

func (s *AgentServiceServer) PublishDeployments(ctx context.Context, pdr *agentservice.PublishDeploymentsRequest) (*agentservice.PublishDeploymentsResponse, error) {
	return &agentservice.PublishDeploymentsResponse{
		Success: true,
	}, nil
}

func main() {

	grpcServer := grpc.NewServer()
	agentservice.RegisterAgentServiceServer(grpcServer, newAgentServiceServer())

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Starting agent server on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
