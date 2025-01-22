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

func (s *presidioRedaction) callPresidioAnalyzer(ctx context.Context, value string) ([]*PresidioAnalyzerResponse, error) {
	requestPayload := PresidioAnalyzerRequest{
		Text:           value,
		Language:       "en",
		ScoreThreshold: s.config.AnalyzerConfig.ScoreThreshold,
		Entities:       s.config.AnalyzerConfig.Entities,
		Context:        s.config.AnalyzerConfig.Context,
	}

	if (isStringHTTPUrl(s.config.PresidioServiceConfig.AnalyzerEndpoint)) {
		return s.callPresidioAnalyzerHTTP(ctx, requestPayload)
	}
	if (isStringGRPCUrl(s.config.PresidioServiceConfig.AnalyzerEndpoint)) {
		return s.callPresidioAnalyzerGRPC(ctx, requestPayload)
	}

	return []*PresidioAnalyzerResponse{}, fmt.Errorf("invalid analyzer endpoint: %s", s.config.PresidioServiceConfig.AnalyzerEndpoint)
}

func (s *presidioRedaction) callPresidioAnalyzerHTTP(ctx context.Context, requestPayload PresidioAnalyzerRequest) ([]*PresidioAnalyzerResponse, error) {
	opts := protojson.MarshalOptions{
		UseProtoNames: true,
		EmitDefaultValues: true,
	}
	jsonPayload, err := opts.Marshal(&requestPayload)
	if err != nil {
		return []*PresidioAnalyzerResponse{}, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	url := s.config.PresidioServiceConfig.AnalyzerEndpoint
	resp, err := s.sendHTTPRequest(ctx, http.MethodPost, url, jsonPayload, headers)
	if err != nil {
		return []*PresidioAnalyzerResponse{}, err
	}
	defer resp.Body.Close()

	var presidioAnalyzerResponse []*PresidioAnalyzerResponse
	err = json.NewDecoder(resp.Body).Decode(&presidioAnalyzerResponse)
	if err != nil {
		return []*PresidioAnalyzerResponse{}, err
	}

	return presidioAnalyzerResponse, nil
}

func (s *presidioRedaction) callPresidioAnalyzerGRPC(ctx context.Context, requestPayload PresidioAnalyzerRequest) ([]*PresidioAnalyzerResponse, error) {
	connStr := strings.TrimPrefix(s.config.PresidioServiceConfig.AnalyzerEndpoint, "grpc://")
	conn, err := grpc.Dial(connStr, grpc.WithInsecure())
	if err != nil {
		return []*PresidioAnalyzerResponse{}, fmt.Errorf("failed to dial gRPC server: %v", err)
	}

	client := NewPresidioRedactionProcessorClient(conn)
	defer conn.Close()

	response, err := client.Analyze(ctx, &requestPayload)
	if err != nil {
		return []*PresidioAnalyzerResponse{}, fmt.Errorf("failed to call gRPC server: %v", err)
	}

	return response.AnalyzerResults, nil
}

