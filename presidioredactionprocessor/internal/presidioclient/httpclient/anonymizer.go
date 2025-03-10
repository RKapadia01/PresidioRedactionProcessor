// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/RKapadia01/PresidioRedactionProcessor/presidioredactionprocessor/internal/presidioclient"
	"google.golang.org/protobuf/encoding/protojson"
)

func (s *PresidioHttpClient) CallPresidioAnonymizer(ctx context.Context, value string, analyzerResults []*presidioclient.PresidioAnalyzerResponse) (presidioclient.PresidioAnonymizerResponse, error) {
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

	requestPayload := presidioclient.PresidioAnonymizerRequest{
		Text:            value,
		Anonymizers:     anonymizers,
		AnalyzerResults: analyzerResults,
	}

	return s.callPresidioAnonymizerHTTP(ctx, requestPayload)
}

func (s *PresidioHttpClient) callPresidioAnonymizerHTTP(ctx context.Context, requestPayload presidioclient.PresidioAnonymizerRequest) (presidioclient.PresidioAnonymizerResponse, error) {
	opts := protojson.MarshalOptions{
		UseProtoNames:     true,
		EmitDefaultValues: true,
	}
	jsonPayload, err := opts.Marshal(&requestPayload)
	if err != nil {
		return presidioclient.PresidioAnonymizerResponse{}, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	url := s.config.PresidioServiceConfig.AnonymizerEndpoint
	resp, err := s.sendHTTPRequest(ctx, http.MethodPost, url, jsonPayload, headers)
	if err != nil {
		return presidioclient.PresidioAnonymizerResponse{}, err
	}
	defer resp.Body.Close()

	var presidioAnonymizerResponse presidioclient.PresidioAnonymizerResponse
	err = json.NewDecoder(resp.Body).Decode(&presidioAnonymizerResponse)
	if err != nil {
		return presidioclient.PresidioAnonymizerResponse{}, fmt.Errorf("failed to unmarshal response payload: %v", err)
	}

	return presidioAnonymizerResponse, nil
}
