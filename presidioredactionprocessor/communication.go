// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package presidioredactionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/redactionprocessor"

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

func (s *presidioRedaction) sendHTTPRequest(ctx context.Context, method, url string, payload []byte, headers map[string]string) (*http.Response, error) {
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
