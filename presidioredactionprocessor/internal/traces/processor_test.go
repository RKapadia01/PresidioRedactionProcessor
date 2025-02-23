package traces

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/RKapadia01/presidioredactionprocessor/internal/presidioclient"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

// mockHTTPClient creates a client that returns predefined responses for analyzer and anonymizer
func mockHTTPClient() *http.Client {
	return &http.Client{
		Transport: &mockTransport{},
	}
}

type mockTransport struct{}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	switch {
	case strings.Contains(req.URL.String(), "analyzer"):
		resp := map[string][]map[string]interface{}{
			"analyzer_results": {
				{
					"start":       0,
					"end":         14,
					"score":       0.8,
					"entity_type": "EMAIL_ADDRESS",
					"text":        "test@example.com",
				},
			},
		}
		jsonResp, _ := json.Marshal(resp)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(string(jsonResp))),
			Header:     make(http.Header),
		}, nil

	case strings.Contains(req.URL.String(), "anonymizer"):
		resp := map[string]string{
			"text": "<REDACTED>",
		}
		jsonResp, _ := json.Marshal(resp)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(string(jsonResp))),
			Header:     make(http.Header),
		}, nil

	default:
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader(`{"error": "invalid endpoint"}`)),
			Header:     make(http.Header),
		}, nil
	}
}

func TestNewPresidioTraceRedaction(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Save the original HTTP client and restore it after the test
	origClient := http.DefaultClient
	http.DefaultClient = mockHTTPClient()
	defer func() { http.DefaultClient = origClient }()

	cfg := &presidioclient.PresidioRedactionProcessorConfig{
		PresidioRunMode: "external",
		ErrorMode:       ottl.PropagateError,
		PresidioServiceConfig: presidioclient.PresidioServiceConfig{
			AnalyzerEndpoint:   "http://analyzer.test",
			AnonymizerEndpoint: "http://anonymizer.test",
		},
		AnalyzerConfig: presidioclient.AnalyzerConfig{
			Language:       "en",
			ScoreThreshold: 0.5,
		},
	}

	err := cfg.Validate()
	require.NoError(t, err, "minimal config should be valid")

	processor := NewPresidioTraceRedaction(ctx, cfg, componenttest.NewNopTelemetrySettings(), logger)
	require.NotNil(t, processor, "processor should not be nil with minimal valid configuration")
	require.NotNil(t, processor.PresidioRedaction, "base redaction should not be nil")
}

func TestProcessTraces(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	// Save the original HTTP client and restore it after the test
	origClient := http.DefaultClient
	http.DefaultClient = mockHTTPClient()
	defer func() { http.DefaultClient = origClient }()

	cfg := &presidioclient.PresidioRedactionProcessorConfig{
		PresidioRunMode: "external",
		ErrorMode:       ottl.PropagateError,
		PresidioServiceConfig: presidioclient.PresidioServiceConfig{
			AnalyzerEndpoint:   "http://analyzer.test",
			AnonymizerEndpoint: "http://anonymizer.test",
		},
		AnalyzerConfig: presidioclient.AnalyzerConfig{
			Language:       "en",
			ScoreThreshold: 0.5,
		},
	}

	processor := NewPresidioTraceRedaction(ctx, cfg, componenttest.NewNopTelemetrySettings(), logger)
	require.NotNil(t, processor)

	t.Run("process empty traces", func(t *testing.T) {
		traces := ptrace.NewTraces()
		result, err := processor.ProcessTraces(ctx, traces)
		assert.NoError(t, err)
		assert.Equal(t, traces, result)
	})

	t.Run("process traces with PII data", func(t *testing.T) {
		traces := ptrace.NewTraces()
		rs := traces.ResourceSpans().AppendEmpty()
		rs.Resource().Attributes().PutStr("email", "test@example.com")

		result, err := processor.ProcessTraces(ctx, traces)
		assert.NoError(t, err)

		// Verify resource attributes were redacted
		resAttrs := result.ResourceSpans().At(0).Resource().Attributes()
		emailVal, exists := resAttrs.Get("email")
		assert.True(t, exists)
		assert.Equal(t, "<REDACTED>", emailVal.Str())
	})

	t.Run("process traces with nil context", func(t *testing.T) {
		traces := ptrace.NewTraces()
		result, err := processor.ProcessTraces(nil, traces)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProcessResourceSpan(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	origClient := http.DefaultClient
	http.DefaultClient = mockHTTPClient()
	defer func() { http.DefaultClient = origClient }()

	cfg := &presidioclient.PresidioRedactionProcessorConfig{
		PresidioRunMode: "external",
		ErrorMode:       ottl.PropagateError,
		PresidioServiceConfig: presidioclient.PresidioServiceConfig{
			AnalyzerEndpoint:   "http://analyzer.test",
			AnonymizerEndpoint: "http://anonymizer.test",
		},
		AnalyzerConfig: presidioclient.AnalyzerConfig{
			Language:       "en",
			ScoreThreshold: 0.5,
		},
	}

	processor := NewPresidioTraceRedaction(ctx, cfg, componenttest.NewNopTelemetrySettings(), logger)
	require.NotNil(t, processor)

	t.Run("process empty resource span", func(t *testing.T) {
		rs := ptrace.NewResourceSpans()
		err := processor.processResourceSpan(ctx, rs)
		assert.NoError(t, err)
	})
}
