package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/RKapadia01/PresidioRedactionProcessor/internal/presidioclient"
	"go.uber.org/zap"
)

var _ presidioclient.PresidioClient = (*PresidioHttpClient)(nil)

type PresidioHttpClient struct {
	config             *presidioclient.PresidioRedactionProcessorConfig
	logger             *zap.Logger
	client             *http.Client
	concurrencyLimiter chan struct{}
}

func NewPresidioHttpClient(cfg *presidioclient.PresidioRedactionProcessorConfig, logger *zap.Logger) *PresidioHttpClient {
	return &PresidioHttpClient{
		config:             cfg,
		logger:             logger,
		client:             &http.Client{},
		concurrencyLimiter: make(chan struct{}, cfg.PresidioServiceConfig.ConcurrencyLimit),
	}
}

func (s *PresidioHttpClient) ProcessText(ctx context.Context, text string) (string, error) {
	analyzerResults, err := s.CallPresidioAnalyzer(ctx, text)
	if err != nil {
		return "", fmt.Errorf("failed to call presidio analyzer: %v", err)
	}

	if len(analyzerResults) == 0 {
		return text, nil
	}

	anonymizerResponse, err := s.CallPresidioAnonymizer(ctx, text, analyzerResults)
	if err != nil {
		return "", fmt.Errorf("failed to call presidio anonymizer: %v", err)
	}

	return anonymizerResponse.Text, nil
}

func (s *PresidioHttpClient) sendHTTPRequest(ctx context.Context, method, url string, payload []byte, headers map[string]string) (*http.Response, error) {
	// Set a concurrency limiter to avoid overloading the presidio service
	s.concurrencyLimiter <- struct{}{}
	defer func() { <-s.concurrencyLimiter }()

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("service returned status code %d", resp.StatusCode)
	}

	return resp, nil
}
