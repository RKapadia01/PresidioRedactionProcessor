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

func TestCallPresidioAnalyzerSuccess(t *testing.T) {
	// Mock server to simulate Presidio Analyzer API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var request PresidioAnalyzerRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		response := []PresidioAnalyzerResponse{
			{EntityType: "PERSON", Start: 0, End: 5, Score: 0.85},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Initialize the processor
	logger, _ := zap.NewProduction()
	config := &PresidioRedactionProcessorConfig{
		AnalyzerConfig: AnalyzerConfig{
			ScoreThreshold: 0.5,
			Entities:       []string{"PERSON"},
			Context:        []string{"context"},
		},
		PresidioServiceConfig: PresidioServiceConfig{
			AnalyzerEndpoint: mockServer.URL,
		},
	}
	processor := newPresidioRedaction(context.Background(), config, logger)

	// Call the function
	ctx := context.Background()
	value := "John Doe"
	response, err := processor.callPresidioAnalyzer(ctx, value)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "PERSON", response[0].EntityType)
}

func TestCallPresidioAnalyzerWithHTTPRequestError(t *testing.T) {
	logger, _ := zap.NewProduction()
	config := &PresidioRedactionProcessorConfig{
		PresidioServiceConfig: PresidioServiceConfig{
			AnalyzerEndpoint: "http://invalid-url",
		},
	}
	processor := newPresidioRedaction(context.Background(), config, logger)

	// Call the function
	ctx := context.Background()
	value := "John Doe"
	_, err := processor.callPresidioAnalyzer(ctx, value)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute HTTP request")
}

func TestCallPresidioAnalyzerNonOKStatusCode(t *testing.T) {
	// Mock server to simulate Presidio Analyzer API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	// Initialize the processor
	logger, _ := zap.NewProduction()
	config := &PresidioRedactionProcessorConfig{
		PresidioServiceConfig: PresidioServiceConfig{
			AnalyzerEndpoint: mockServer.URL,
		},
	}
	processor := newPresidioRedaction(context.Background(), config, logger)

	// Call the function
	ctx := context.Background()
	value := "John Doe"
	_, err := processor.callPresidioAnalyzer(ctx, value)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "service returned status code 500")
}

func TestCallPresidioAnalyzerJSONDecodeError(t *testing.T) {
	// Mock server to simulate Presidio Analyzer API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer mockServer.Close()

	// Initialize the processor
	logger, _ := zap.NewProduction()
	config := &PresidioRedactionProcessorConfig{
		PresidioServiceConfig: PresidioServiceConfig{
			AnalyzerEndpoint: mockServer.URL,
		},
	}
	processor := newPresidioRedaction(context.Background(), config, logger)

	// Call the function
	ctx := context.Background()
	value := "John Doe"
	_, err := processor.callPresidioAnalyzer(ctx, value)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}
