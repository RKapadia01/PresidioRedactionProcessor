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

func TestCallPresidioAnonymizerSuccess(t *testing.T) {
	// Mock server to simulate Presidio Anonymizer API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var request PresidioAnonymizerRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		response := PresidioAnonymizerResponse{
			Text: "*****",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Initialize the processor
	logger, _ := zap.NewProduction()
	config := &PresidioRedactionProcessorConfig{
		AnonymizerConfig: AnonymizerConfig{
			Anonymizers: []EntityAnonymizer{
				{
					Entity:      "PERSON",
					Type:        "mask",
					NewValue:    "",
					MaskingChar: "*",
					CharsToMask: 5,
					FromEnd:     false,
					HashType:    "",
					Key:         "",
				},
			},
		},
		PresidioServiceConfig: PresidioServiceConfig{
			AnonymizerEndpoint: mockServer.URL,
		},
	}
	processor := newPresidioRedaction(context.Background(), config, logger)

	// Call the function
	ctx := context.Background()
	value := "John Doe"
	analyzerResults := []*PresidioAnalyzerResponse{
		{EntityType: "PERSON", Start: 0, End: 5, Score: 0.85},
	}
	response, err := processor.callPresidioAnonymizer(ctx, value, analyzerResults)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "*****", response.Text)
}

func TestCallPresidioAnonymizerHTTPRequestError(t *testing.T) {
	logger, _ := zap.NewProduction()
	config := &PresidioRedactionProcessorConfig{
		PresidioServiceConfig: PresidioServiceConfig{
			AnonymizerEndpoint: "http://invalid-url",
		},
	}
	processor := newPresidioRedaction(context.Background(), config, logger)

	// Call the function
	ctx := context.Background()
	value := "John Doe"
	analyzerResults := []*PresidioAnalyzerResponse{}
	_, err := processor.callPresidioAnonymizer(ctx, value, analyzerResults)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute HTTP request")
}

func TestCallPresidioAnonymizerNonOKStatusCode(t *testing.T) {
	// Mock server to simulate Presidio Anonymizer API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	// Initialize the processor
	logger, _ := zap.NewProduction()
	config := &PresidioRedactionProcessorConfig{
		PresidioServiceConfig: PresidioServiceConfig{
			AnonymizerEndpoint: mockServer.URL,
		},
	}
	processor := newPresidioRedaction(context.Background(), config, logger)

	// Call the function
	ctx := context.Background()
	value := "John Doe"
	analyzerResults := []*PresidioAnalyzerResponse{}
	_, err := processor.callPresidioAnonymizer(ctx, value, analyzerResults)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "service returned status code 500")
}

func TestCallPresidioAnonymizerJSONDecodeError(t *testing.T) {
	// Mock server to simulate Presidio Anonymizer API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer mockServer.Close()

	// Initialize the processor
	logger, _ := zap.NewProduction()
	config := &PresidioRedactionProcessorConfig{
		PresidioServiceConfig: PresidioServiceConfig{
			AnonymizerEndpoint: mockServer.URL,
		},
	}
	processor := newPresidioRedaction(context.Background(), config, logger)

	// Call the function
	ctx := context.Background()
	value := "John Doe"
	analyzerResults := []*PresidioAnalyzerResponse{}
	_, err := processor.callPresidioAnonymizer(ctx, value, analyzerResults)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestCallPresidioAnonymizer_DifferentAnonymizerTypes(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request PresidioAnonymizerRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify different anonymizer types
		assert.Equal(t, "mask", request.Anonymizers["PERSON"].Type)
		assert.Equal(t, "hash", request.Anonymizers["EMAIL"].Type)

		response := &PresidioAnonymizerResponse{
			Text: "***** abc123",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	logger, _ := zap.NewProduction()
	config := &PresidioRedactionProcessorConfig{
		PresidioServiceConfig: PresidioServiceConfig{
			AnonymizerEndpoint: mockServer.URL,
		},
		AnonymizerConfig: AnonymizerConfig{
			Anonymizers: []EntityAnonymizer{
				{
					Entity:      "PERSON",
					Type:        "mask",
					MaskingChar: "*",
				},
				{
					Entity:   "EMAIL",
					Type:     "hash",
					HashType: "md5",
				},
			},
		},
	}
	processor := newPresidioRedaction(context.Background(), config, logger)

	ctx := context.Background()
	analyzerResults := []*PresidioAnalyzerResponse{
		{EntityType: "PERSON", Start: 0, End: 8, Score: 0.85},
		{EntityType: "EMAIL", Start: 9, End: 25, Score: 0.9},
	}

	response, err := processor.callPresidioAnonymizer(ctx, "John Doe john@example.com", analyzerResults)
	assert.NoError(t, err)
	assert.Equal(t, "***** abc123", response.Text)
}
