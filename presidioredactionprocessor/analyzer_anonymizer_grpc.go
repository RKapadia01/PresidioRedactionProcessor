// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package presidioredactionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/redactionprocessor"

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	
	"google.golang.org/grpc"
    "google.golang.org/protobuf/encoding/protojson"
)

func (s *presidioRedaction) callPresidioGRPC(ctx context.Context, value string) (PresidioAnonymizerResponse, error) {
	requestPayload := PresidioAnalyzerRequest{
		Text:           value,
		Language:       "en",
		ScoreThreshold: s.config.AnalyzerConfig.ScoreThreshold,
		Entities:       s.config.AnalyzerConfig.Entities,
		Context:        s.config.AnalyzerConfig.Context,
	}

	return s.callPresidioAnalyzerAndAnonymizer(ctx, requestPayload)
}

func (s *presidioRedaction) callPresidioAnalyzerGRPC(ctx context.Context, requestPayload PresidioAnalyzerRequest) (PresidioAnonymizerResponse, error) {
	connStr := strings.TrimPrefix(s.config.PresidioServiceConfig.AnalyzerEndpoint, "grpc://")
	conn, err := grpc.Dial(connStr, grpc.WithInsecure())
	if err != nil {
		return []*PresidioAnalyzerResponse{}, fmt.Errorf("failed to dial gRPC server: %v", err)
	}

	client := NewPresidioRedactionProcessorClient(conn)
	defer conn.Close()

	response, err := client.AnalyzeAndAnonymize(ctx, &requestPayload)
	if err != nil {
		return PresidioAnonymizerResponse{}, fmt.Errorf("failed to call gRPC server: %v", err)
	}

	return *response, nil
}

