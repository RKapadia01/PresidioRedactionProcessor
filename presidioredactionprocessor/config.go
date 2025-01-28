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
	PresidioServiceConfig PresidioServiceConfig `mapstructure:"presidio_service"`
	AnalyzerConfig        AnalyzerConfig        `mapstructure:"analyzer"`
	AnonymizerConfig      AnonymizerConfig      `mapstructure:"anonymizer,omitempty"`
}

type PresidioServiceConfig struct {
	AnalyzerEndpoint   string   `mapstructure:"analyzer_endpoint"`
	AnonymizerEndpoint string   `mapstructure:"anonymizer_endpoint"`
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

func (c *PresidioServiceConfig) validate() error {
	if c.AnalyzerEndpoint == "" {
		return errors.New("analyzer_endpoint is required")
	}
	if c.AnonymizerEndpoint == "" {
		return errors.New("anonymizer_endpoint is required")
	}

	return nil
}
