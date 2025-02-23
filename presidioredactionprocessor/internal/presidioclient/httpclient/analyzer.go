// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RKapadia01/presidioredactionprocessor/internal/presidioclient"
	"google.golang.org/protobuf/encoding/protojson"
)

func (s *PresidioHttpClient) CallPresidioAnalyzer(ctx context.Context, value string) ([]*presidioclient.PresidioAnalyzerResponse, error) {
	requestPayload := presidioclient.PresidioAnalyzerRequest{
		Text:           value,
		Language:       "en",
		ScoreThreshold: s.config.AnalyzerConfig.ScoreThreshold,
		Entities:       s.config.AnalyzerConfig.Entities,
		Context:        s.config.AnalyzerConfig.Context,
	}

	return s.callPresidioAnalyzerHTTP(ctx, requestPayload)
}

func (s *PresidioHttpClient) callPresidioAnalyzerHTTP(ctx context.Context, requestPayload presidioclient.PresidioAnalyzerRequest) ([]*presidioclient.PresidioAnalyzerResponse, error) {
	opts := protojson.MarshalOptions{
		UseProtoNames:     true,
		EmitDefaultValues: true,
	}
	jsonPayload, err := opts.Marshal(&requestPayload)
	if err != nil {
		return []*presidioclient.PresidioAnalyzerResponse{}, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	url := s.config.PresidioServiceConfig.AnalyzerEndpoint
	resp, err := s.sendHTTPRequest(ctx, http.MethodPost, url, jsonPayload, headers)
	if err != nil {
		return []*presidioclient.PresidioAnalyzerResponse{}, err
	}
	defer resp.Body.Close()

	var presidioAnalyzerResponse []*presidioclient.PresidioAnalyzerResponse
	err = json.NewDecoder(resp.Body).Decode(&presidioAnalyzerResponse)
	if err != nil {
		return []*presidioclient.PresidioAnalyzerResponse{}, fmt.Errorf("failed to unmarshal response payload: %v", err)
	}

	return presidioAnalyzerResponse, nil
}
