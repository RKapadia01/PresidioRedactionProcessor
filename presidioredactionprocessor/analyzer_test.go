package presidioredactionprocessor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

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

func TestCallPresidioAnalyzer_ValidateConfig(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request PresidioAnalyzerRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		assert.NoError(t, err)

		// Verify request contains config values
		assert.Equal(t, 0.7, request.ScoreThreshold)
		assert.Equal(t, []string{"PERSON"}, request.Entities)
		assert.Equal(t, []string{"context"}, request.Context)

		response := []*PresidioAnalyzerResponse{}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	logger, _ := zap.NewProduction()
	config := &PresidioRedactionProcessorConfig{
		PresidioServiceConfig: PresidioServiceConfig{
			AnalyzerEndpoint: mockServer.URL,
		},
		AnalyzerConfig: AnalyzerConfig{
			ScoreThreshold: 0.7,
			Entities:       []string{"PERSON"},
			Context:        []string{"context"},
		},
	}
	processor := newPresidioRedaction(context.Background(), config, logger)

	ctx := context.Background()
	_, err := processor.callPresidioAnalyzer(ctx, "John Doe")
	assert.NoError(t, err) // Throw no error since the config matches
}

func TestProcessor_RespectsMaxConcurrentRequests(t *testing.T) {
	// Mock server with deliberate delay
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []*PresidioAnalyzerResponse{}
		time.Sleep(100 * time.Millisecond) // Ensure requests overlap
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	logger, _ := zap.NewProduction()
	config := &PresidioRedactionProcessorConfig{
		PresidioServiceConfig: PresidioServiceConfig{
			AnalyzerEndpoint: mockServer.URL,
			ConcurrencyLimit: 2,
		},
	}
	processor := newPresidioRedaction(context.Background(), config, logger)

	var wg sync.WaitGroup
	var concurrent int32
	var maxConcurrent int32

	// Launch requests
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			processor.concurrencyLimiter <- struct{}{}        // Acquire semaphore and block if full
			defer func() { <-processor.concurrencyLimiter }() // Ensures release

			// Track concurrent executions
			currCount := atomic.AddInt32(&concurrent, 1)
			if currCount > atomic.LoadInt32(&maxConcurrent) {
				atomic.StoreInt32(&maxConcurrent, currCount)
			}

			// Simulate work
			time.Sleep(50 * time.Millisecond)

			// Decrement count
			atomic.AddInt32(&concurrent, -1)
		}()
	}

	wg.Wait()
	assert.LessOrEqual(t, atomic.LoadInt32(&maxConcurrent), int32(2),
		"Concurrent requests exceeded limit")
}
