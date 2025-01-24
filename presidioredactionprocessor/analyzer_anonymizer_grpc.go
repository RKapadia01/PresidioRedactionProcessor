// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package presidioredactionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/redactionprocessor"

import (
	"context"
	"fmt"
	"strings"
	
	"google.golang.org/grpc"
)

func (s *presidioRedaction) callPresidioGRPC(ctx context.Context, value string) (PresidioAnonymizerResponse, error) {
	anonymizers := make(map[string]*PresidioAnonymizer)

	for _, entityAnonymizer := range s.config.AnonymizerConfig.Anonymizers {
		anonymizers[entityAnonymizer.Entity] = &PresidioAnonymizer{
			Type:        strings.ToLower(entityAnonymizer.Type),
			NewValue:    entityAnonymizer.NewValue,
			MaskingChar: entityAnonymizer.MaskingChar,
			CharsToMask: entityAnonymizer.CharsToMask,
			FromEnd:     entityAnonymizer.FromEnd,
			HashType:    entityAnonymizer.HashType,
			Key:         entityAnonymizer.Key,
		}
	}

	requestPayload := PresidioAnalyzerAnomymizerRequest{
		Text:           value,
		Language:       "en",
		ScoreThreshold: s.config.AnalyzerConfig.ScoreThreshold,
		Entities:       s.config.AnalyzerConfig.Entities,
		Context:        s.config.AnalyzerConfig.Context,
		Anonymizers:    anonymizers,
	}

	return s.callPresidioAnalyzerAndAnonymizer(ctx, requestPayload)
}

func (s *presidioRedaction) callPresidioAnalyzerAndAnonymizer(ctx context.Context, requestPayload PresidioAnalyzerAnomymizerRequest) (PresidioAnonymizerResponse, error) {
	connStr := strings.TrimPrefix(s.config.PresidioServiceConfig.AnalyzerEndpoint, "grpc://")
	conn, err := grpc.Dial(connStr, grpc.WithInsecure())
	if err != nil {
		return PresidioAnonymizerResponse{}, fmt.Errorf("failed to dial gRPC server: %v", err)
	}

	client := NewPresidioRedactionProcessorClient(conn)
	defer conn.Close()

	response, err := client.AnalyzeAndAnonymize(ctx, &requestPayload)
	if err != nil {
		return PresidioAnonymizerResponse{}, fmt.Errorf("failed to call gRPC server: %v", err)
	}

	return *response, nil
}

