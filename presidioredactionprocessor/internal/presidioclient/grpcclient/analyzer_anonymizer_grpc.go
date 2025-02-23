// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package grpcclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/RKapadia01/PresidioRedactionProcessor/internal/presidioclient"
	"google.golang.org/grpc"
)

func (s *PresidioGrpcClient) CallPresidioGRPC(ctx context.Context, value string) (presidioclient.PresidioAnonymizerResponse, error) {
	anonymizers := make(map[string]*presidioclient.PresidioAnonymizer)

	for _, entityAnonymizer := range s.config.AnonymizerConfig.Anonymizers {
		anonymizers[entityAnonymizer.Entity] = &presidioclient.PresidioAnonymizer{
			Type:        strings.ToLower(entityAnonymizer.Type),
			NewValue:    entityAnonymizer.NewValue,
			MaskingChar: entityAnonymizer.MaskingChar,
			CharsToMask: entityAnonymizer.CharsToMask,
			FromEnd:     entityAnonymizer.FromEnd,
			HashType:    entityAnonymizer.HashType,
			Key:         entityAnonymizer.Key,
		}
	}

	requestPayload := presidioclient.PresidioAnalyzerAnomymizerRequest{
		Text:           value,
		Language:       "en",
		ScoreThreshold: s.config.AnalyzerConfig.ScoreThreshold,
		Entities:       s.config.AnalyzerConfig.Entities,
		Context:        s.config.AnalyzerConfig.Context,
		Anonymizers:    anonymizers,
	}

	return s.callPresidioAnalyzerAndAnonymizer(ctx, requestPayload)
}

func (s *PresidioGrpcClient) callPresidioAnalyzerAndAnonymizer(ctx context.Context, requestPayload presidioclient.PresidioAnalyzerAnomymizerRequest) (presidioclient.PresidioAnonymizerResponse, error) {
	connStr := strings.TrimPrefix(s.config.PresidioServiceConfig.AnalyzerEndpoint, "grpc://")
	conn, err := grpc.Dial(connStr, grpc.WithInsecure())
	if err != nil {
		return presidioclient.PresidioAnonymizerResponse{}, fmt.Errorf("failed to dial gRPC server: %v", err)
	}

	client := presidioclient.NewPresidioRedactionProcessorClient(conn)
	defer conn.Close()

	response, err := client.AnalyzeAndAnonymize(ctx, &requestPayload)
	if err != nil {
		return presidioclient.PresidioAnonymizerResponse{}, fmt.Errorf("failed to call gRPC server: %v", err)
	}

	return *response, nil
}
