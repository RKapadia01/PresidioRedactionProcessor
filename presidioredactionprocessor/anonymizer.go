// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package presidioredactionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/redactionprocessor"

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
)

func (s *presidioRedaction) callPresidioAnonymizer(ctx context.Context, value string, analyzerResults []*PresidioAnalyzerResponse) (PresidioAnonymizerResponse, error) {
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

	requestPayload := PresidioAnonymizerRequest{
		Text:            value,
		Anonymizers:     anonymizers,
		AnalyzerResults: analyzerResults,
	}

	if isStringHTTPUrl(s.config.PresidioServiceConfig.AnonymizerEndpoint) {
		return s.callPresidioAnonymizerHTTP(ctx, requestPayload)
	}

	return PresidioAnonymizerResponse{}, fmt.Errorf("invalid anonymizer endpoint: %s", s.config.PresidioServiceConfig.AnonymizerEndpoint)
}

func (s *presidioRedaction) callPresidioAnonymizerHTTP(ctx context.Context, requestPayload PresidioAnonymizerRequest) (PresidioAnonymizerResponse, error) {
	opts := protojson.MarshalOptions{
		UseProtoNames:     true,
		EmitDefaultValues: true,
	}
	jsonPayload, err := opts.Marshal(&requestPayload)
	if err != nil {
		return PresidioAnonymizerResponse{}, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	url := s.config.PresidioServiceConfig.AnonymizerEndpoint
	resp, err := s.sendHTTPRequest(ctx, http.MethodPost, url, jsonPayload, headers)
	if err != nil {
		return PresidioAnonymizerResponse{}, err
	}
	defer resp.Body.Close()

	var presidioAnonymizerResponse PresidioAnonymizerResponse
	err = json.NewDecoder(resp.Body).Decode(&presidioAnonymizerResponse)
	if err != nil {
		return PresidioAnonymizerResponse{}, fmt.Errorf("failed to unmarshal response payload: %v", err)
	}

	return presidioAnonymizerResponse, nil
}
