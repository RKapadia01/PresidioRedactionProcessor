// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package presidioclient

import (
	"fmt"
	"strings"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"go.opentelemetry.io/collector/component"
)

func CreateDefaultConfig() component.Config {
	return &PresidioRedactionProcessorConfig{
		PresidioRunMode: "embedded",
		ErrorMode:       ottl.PropagateError,
		PresidioServiceConfig: PresidioServiceConfig{
			AnalyzerEndpoint:   "grpc://localhost:50051",
			AnonymizerEndpoint: "grpc://localhost:50052",
		},
		AnalyzerConfig: AnalyzerConfig{
			Language:       "en",
			ScoreThreshold: 0.5,
		},
	}
}

type PresidioRedactionProcessorConfig struct {
	PresidioRunMode       string                `mapstructure:"mode"`
	ErrorMode             ottl.ErrorMode        `mapstructure:"error_mode"`
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

func (c *PresidioRedactionProcessorConfig) Validate() error {
	if c.PresidioRunMode != "embedded" && c.PresidioRunMode != "external" {
		return fmt.Errorf("invalid presidio_run_mode: %s, must be either 'embedded' or 'external'", c.PresidioRunMode)
	}

	if c.PresidioRunMode == "external" {
		if err := c.validateExternalServiceConfig(); err != nil {
			return err
		}
	}

	for i, anonymizer := range c.AnonymizerConfig.Anonymizers {
		if err := validateEntityAnonymizer(anonymizer, i); err != nil {
			return err
		}
	}

	return nil
}

func (c *PresidioRedactionProcessorConfig) validateExternalServiceConfig() error {
	if c.PresidioServiceConfig.AnalyzerEndpoint == "" {
		return fmt.Errorf("presidio_service.analyzer_endpoint is required when mode is external")
	}

	if !strings.HasPrefix(c.PresidioServiceConfig.AnalyzerEndpoint, "http://") &&
		!strings.HasPrefix(c.PresidioServiceConfig.AnalyzerEndpoint, "https://") {
		return fmt.Errorf("invalid analyzer_endpoint protocol: %s, must start with http:// or https://",
			c.PresidioServiceConfig.AnalyzerEndpoint)
	}

	if c.PresidioServiceConfig.AnonymizerEndpoint == "" {
		return fmt.Errorf("presidio_service.anonymizer_endpoint is required when mode is service")
	}
	if !strings.HasPrefix(c.PresidioServiceConfig.AnonymizerEndpoint, "http://") &&
		!strings.HasPrefix(c.PresidioServiceConfig.AnonymizerEndpoint, "https://") {
		return fmt.Errorf("invalid anonymizer_endpoint protocol: %s, must start with http:// or https://",
			c.PresidioServiceConfig.AnonymizerEndpoint)
	}
	return nil
}

func validateEntityAnonymizer(anonymizer EntityAnonymizer, index int) error {
	if anonymizer.Entity == "" {
		return fmt.Errorf("anonymizer[%d].entity is required", index)
	}

	validTypes := map[string]bool{
		"REPLACE": true,
		"HASH":    true,
		"REDACT":  true,
		"MASK":    true,
		"ENCRYPT": true,
	}

	if !validTypes[anonymizer.Type] {
		return fmt.Errorf("anonymizer[%d].type has invalid value: %s", index, anonymizer.Type)
	}

	switch anonymizer.Type {
	case "REPLACE":
		if anonymizer.NewValue == "" {
			return fmt.Errorf("anonymizer[%d].new_value is required for REPLACE type", index)
		}
	case "HASH":
		if anonymizer.HashType == "" {
			return fmt.Errorf("anonymizer[%d].hash_type is required for HASH type", index)
		}
		validHashTypes := map[string]bool{"sha256": true, "sha512": true, "md5": true}
		if !validHashTypes[anonymizer.HashType] {
			return fmt.Errorf("anonymizer[%d].hash_type has invalid value: %s", index, anonymizer.HashType)
		}
	case "MASK":
		if anonymizer.CharsToMask <= 0 {
			return fmt.Errorf("anonymizer[%d].chars_to_mask must be greater than 0 for MASK type", index)
		}
		if anonymizer.MaskingChar == "" {
			anonymizer.MaskingChar = "*"
		}
	case "ENCRYPT":
		if anonymizer.Key == "" {
			return fmt.Errorf("anonymizer[%d].key is required for ENCRYPT type", index)
		}
	}

	return nil
}
