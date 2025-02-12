package presidioredactionprocessor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetRedactedValue(t *testing.T) {
	// Mock analyzer server
	mockAnalyzer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []*PresidioAnalyzerResponse{
			{EntityType: "PERSON", Start: 0, End: 8, Score: 0.85},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockAnalyzer.Close()

	// Mock anonymizer server
	mockAnonymizer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := &PresidioAnonymizerResponse{
			Text: "<REDACTED>",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockAnonymizer.Close()

	logger, _ := zap.NewProduction()
	config := &PresidioRedactionProcessorConfig{
		PresidioServiceConfig: PresidioServiceConfig{
			AnalyzerEndpoint:   mockAnalyzer.URL,
			AnonymizerEndpoint: mockAnonymizer.URL,
		},
		AnalyzerConfig: AnalyzerConfig{
			ScoreThreshold: 0.5,
			Entities:       []string{"PERSON"},
		},
		AnonymizerConfig: AnonymizerConfig{
			Anonymizers: []EntityAnonymizer{
				{
					Entity:   "PERSON",
					Type:     "replace",
					NewValue: "<REDACTED>",
				},
			},
		},
	}
	processor := newPresidioRedaction(context.Background(), config, logger)

	ctx := context.Background()
	result, err := processor.getRedactedValue(ctx, "John Doe")
	assert.NoError(t, err)
	assert.Equal(t, "<REDACTED>", result)
}

func TestGetRedactedValue_EmptyString(t *testing.T) {
	// Mock analyzer server
	mockAnalyzer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request PresidioAnalyzerRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		response := []*PresidioAnalyzerResponse{}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockAnalyzer.Close()

	// Mock anonymizer server
	mockAnonymizer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := &PresidioAnonymizerResponse{
			Text: "",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockAnonymizer.Close()

	logger, _ := zap.NewProduction()
	config := &PresidioRedactionProcessorConfig{
		PresidioServiceConfig: PresidioServiceConfig{
			AnalyzerEndpoint:   mockAnalyzer.URL,
			AnonymizerEndpoint: mockAnonymizer.URL,
		},
	}
	processor := newPresidioRedaction(context.Background(), config, logger)

	ctx := context.Background()
	result, err := processor.getRedactedValue(ctx, "")
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetRedactedValue_NoDetectedEntities(t *testing.T) {
	// Mock analyzer server
	mockAnalyzer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []*PresidioAnalyzerResponse{}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockAnalyzer.Close()

	// Mock anonymizer server
	mockAnonymizer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := &PresidioAnonymizerResponse{
			Text: "",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockAnonymizer.Close()

	logger, _ := zap.NewProduction()
	config := &PresidioRedactionProcessorConfig{
		PresidioServiceConfig: PresidioServiceConfig{
			AnalyzerEndpoint:   mockAnalyzer.URL,
			AnonymizerEndpoint: mockAnonymizer.URL,
		},
	}
	processor := newPresidioRedaction(context.Background(), config, logger)

	ctx := context.Background()
	value := "No sensitive data here"
	result, err := processor.getRedactedValue(ctx, value)
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}
