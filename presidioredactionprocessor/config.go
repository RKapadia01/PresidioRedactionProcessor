// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package presidioredactionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/redactionprocessor"

import "errors"

type Config struct {
	AnalyzerEndpoint   string           `mapstructure:"analyzer_endpoint"`
	AnonymizerEndpoint string           `mapstructure:"anonymizer_endpoint"`
	AnalyzerConfig     AnalyzerConfig   `mapstructure:"analyzer"`
	AnonymizerConfig   AnonymizerConfig `mapstructure:"anonymizer,omitempty"`
	ConcurrencyLimit   int              `mapstructure:"concurrency_limit,omitempty"`

	TraceConditions []string `mapstructure:"trace_conditions,omitempty"`
	LogConditions   []string `mapstructure:"log_conditions,omitempty"`
}

func (c *Config) validate() error {
	if c.AnalyzerEndpoint == "" {
		return errors.New("analyzer_endpoint is required")
	}
	if c.AnonymizerEndpoint == "" {
		return errors.New("anonymizer_endpoint is required")
	}

	return nil
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
	CharsToMask int    `mapstructure:"chars_to_mask,omitempty"`
	FromEnd     bool   `mapstructure:"from_end,omitempty"`
	HashType    string `mapstructure:"hash_type,omitempty"`
	Key         string `mapstructure:"key,omitempty"`
}
