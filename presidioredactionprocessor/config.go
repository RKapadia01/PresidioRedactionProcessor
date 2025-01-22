// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package presidioredactionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/redactionprocessor"

import (
	"strings"
)

func isStringHTTPUrl(endpoint string) bool {
	return strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://")
}

func isStringGRPCUrl(endpoint string) bool {
	return strings.HasPrefix(endpoint, "grpc://") || strings.HasPrefix(endpoint, "grpcs://")
}

type PresidioServiceConfig struct {
	AnalyzerEndpoint    string `mapstructure:"analyzer_endpoint"`
	AnonymizerEndpoint  string `mapstructure:"anonymizer_endpoint"`
	ConcurrencyLimit    int    `mapstructure:"concurrency_limit,omitempty"`
}

type Config struct {
	PresidioServiceConfig PresidioServiceConfig `mapstructure:"presidio_service"`
	AnalyzerConfig        AnalyzerConfig   	    `mapstructure:"analyzer"`
	AnonymizerConfig      AnonymizerConfig 	    `mapstructure:"anonymizer,omitempty"`
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
