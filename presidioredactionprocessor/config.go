// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package presidioredactionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/redactionprocessor"

type PresidioServiceConfig struct {
	UseDocker 	              bool   `mapstructure:"use_docker"`
	DockerAnalyzerEndpoint    string `mapstructure:"docker_analyzer_endpoint"`
	DockerAnonymizerEndpoint  string `mapstructure:"docker_anonymizer_endpoint"`
	ConcurrencyLimit          int    `mapstructure:"concurrency_limit,omitempty"`
	PythonPath 	              string `mapstructure:"python_path"`
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
	CharsToMask int    `mapstructure:"chars_to_mask,omitempty"`
	FromEnd     bool   `mapstructure:"from_end,omitempty"`
	HashType    string `mapstructure:"hash_type,omitempty"`
	Key         string `mapstructure:"key,omitempty"`
}
