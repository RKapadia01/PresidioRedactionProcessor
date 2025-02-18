// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package presidioredactionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/redactionprocessor"

import (
	"errors"
	"strings"
)

func isStringHTTPUrl(endpoint string) bool {
	return strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://")
}

func isStringGRPCUrl(endpoint string) bool {
	return strings.HasPrefix(endpoint, "grpc://") || strings.HasPrefix(endpoint, "grpcs://")
}

type PresidioRedactionProcessorConfig struct {
	PresidioRunMode       string                `mapstructure:"mode,omitempty" default:"embedded"`
	PresidioServiceConfig PresidioServiceConfig `mapstructure:"presidio_service"`
	AnalyzerConfig        AnalyzerConfig        `mapstructure:"analyzer"`
	AnonymizerConfig      AnonymizerConfig      `mapstructure:"anonymizer,omitempty"`
}

type PresidioServiceConfig struct {
	AnalyzerEndpoint   string   `mapstructure:"analyzer_endpoint" default:"grpc://localhost:50051"`
	AnonymizerEndpoint string   `mapstructure:"anonymizer_endpoint" default:"grpc://localhost:50052"`
	ConcurrencyLimit   int      `mapstructure:"concurrency_limit,omitempty"`
	TraceConditions    []string `mapstructure:"process_trace_if,omitempty"`
	LogConditions      []string `mapstructure:"process_log_if,omitempty"`
}

type AnalyzerConfig struct {
	Language       string   `mapstructure:"language"`
	ScoreThreshold float64  `mapstructure:"score_threshold,omitempty"`
	Entities       []string `mapstructure:"entities,omitempty"`
	Context        []string `mapstructure:"context,omitempty"`
}

type AnonymizerConfig struct {
	Anonymizers []EntityAnonymizer `mapstructure:"anonymizers"`
}

type EntityAnonymizer struct {
	Entity      string `mapstructure:"entity"`
	Type        string `mapstructure:"type"`
	NewValue    string `mapstructure:"new_value,omitempty"`
	MaskingChar string `mapstructure:"masking_char,omitempty"`
	CharsToMask int32  `mapstructure:"chars_to_mask,omitempty"`
	FromEnd     bool   `mapstructure:"from_end,omitempty"`
	HashType    string `mapstructure:"hash_type,omitempty"`
	Key         string `mapstructure:"key,omitempty"`
}

func (c *PresidioRedactionProcessorConfig) validate() error {
	if c.PresidioRunMode == "service" {
		if c.PresidioServiceConfig.AnalyzerEndpoint == "" {
			return errors.New("presidio_service.analyzer_endpoint is required when mode is external")
		}
		if c.PresidioServiceConfig.AnonymizerEndpoint == "" {
			return errors.New("presidio_service.anonymizer_endpoint is required when mode is external")
		}
	}

	return nil
}
