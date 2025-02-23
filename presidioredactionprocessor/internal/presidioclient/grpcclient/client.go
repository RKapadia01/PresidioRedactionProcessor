package grpcclient

import (
	"context"
	"fmt"

	"github.com/RKapadia01/presidioredactionprocessor/internal/presidioclient"
	"go.uber.org/zap"
)

type PresidioGrpcClient struct {
	config *presidioclient.PresidioRedactionProcessorConfig
	logger *zap.Logger
}

func NewPresidioGrpcClient(cfg *presidioclient.PresidioRedactionProcessorConfig) *PresidioGrpcClient {
	return &PresidioGrpcClient{
		config: cfg,
	}
}

func (s *PresidioGrpcClient) ProcessText(ctx context.Context, value string) (string, error) {
	response, err := s.CallPresidioGRPC(ctx, value)

	if err != nil {
		return "", fmt.Errorf("failed to call Presidio gRPC server: %v", err)
	}
	return response.Text, nil
}
