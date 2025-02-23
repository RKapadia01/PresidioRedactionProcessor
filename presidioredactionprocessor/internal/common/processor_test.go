package common

import (
	"context"
	"fmt"
	"testing"

	"github.com/RKapadia01/presidioredactionprocessor/internal/presidioclient"
	"github.com/RKapadia01/presidioredactionprocessor/internal/presidioclient/grpcclient"
	"github.com/RKapadia01/presidioredactionprocessor/internal/presidioclient/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.uber.org/zap"
)

func TestCreateBaseRedaction(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name          string
		config        *presidioclient.PresidioRedactionProcessorConfig
		expectedError bool
		expectedGRPC  bool
		expectedHTTP  bool
	}{
		{
			name: "Valid GRPC Configuration",
			config: &presidioclient.PresidioRedactionProcessorConfig{
				PresidioServiceConfig: presidioclient.PresidioServiceConfig{
					AnalyzerEndpoint:   "grpc://analyzer:3000",
					AnonymizerEndpoint: "grpc://anonymizer:3000",
					ConcurrencyLimit:   5,
				},
			},
			expectedError: false,
			expectedGRPC:  true,
			expectedHTTP:  false,
		},
		{
			name: "Valid HTTP Configuration",
			config: &presidioclient.PresidioRedactionProcessorConfig{
				PresidioServiceConfig: presidioclient.PresidioServiceConfig{
					AnalyzerEndpoint:   "http://analyzer:3000",
					AnonymizerEndpoint: "http://anonymizer:3000",
					ConcurrencyLimit:   10,
				},
			},
			expectedError: false,
			expectedGRPC:  false,
			expectedHTTP:  true,
		},
		{
			name: "Invalid Mixed Protocol Configuration",
			config: &presidioclient.PresidioRedactionProcessorConfig{
				PresidioServiceConfig: presidioclient.PresidioServiceConfig{
					AnalyzerEndpoint:   "grpc://analyzer:3000",
					AnonymizerEndpoint: "http://anonymizer:3000",
					ConcurrencyLimit:   5,
				},
			},
			expectedError: true,
			expectedGRPC:  false,
			expectedHTTP:  false,
		},
		{
			name: "Invalid Endpoint Configuration",
			config: &presidioclient.PresidioRedactionProcessorConfig{
				PresidioServiceConfig: presidioclient.PresidioServiceConfig{
					AnalyzerEndpoint:   "invalid://analyzer:3000",
					AnonymizerEndpoint: "invalid://anonymizer:3000",
					ConcurrencyLimit:   5,
				},
			},
			expectedError: true,
			expectedGRPC:  false,
			expectedHTTP:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redaction, err := CreateBaseRedaction(tt.config, logger)

			if tt.expectedError {
				require.Error(t, err)
				assert.Nil(t, redaction)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, redaction)

			assert.Equal(t, tt.config, redaction.Config)
			assert.Equal(t, logger, redaction.Logger)
			assert.NotNil(t, redaction.Client)
			assert.NotNil(t, redaction.ConcurrencyLimiter)
			assert.Equal(t, tt.config.PresidioServiceConfig.ConcurrencyLimit, cap(redaction.ConcurrencyLimiter))

			if tt.expectedGRPC {
				assert.IsType(t, &grpcclient.PresidioGrpcClient{}, redaction.Client)
			}
			if tt.expectedHTTP {
				assert.IsType(t, &httpclient.PresidioHttpClient{}, redaction.Client)
			}
		})
	}
}

func TestProcessAttribute(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	tests := []struct {
		name          string
		attributes    map[string]interface{}
		expectedError bool
		expectedValue string
		mockProcessor func(context.Context, string) (string, error)
	}{
		{
			name: "Successfully process attribute",
			attributes: map[string]interface{}{
				"key": "sensitive value",
			},
			expectedError: false,
			expectedValue: "redacted value",
			mockProcessor: func(ctx context.Context, text string) (string, error) {
				return "redacted value", nil
			},
		},
		{
			name: "Empty string attribute",
			attributes: map[string]interface{}{
				"key": "",
			},
			expectedError: false,
			expectedValue: "",
			mockProcessor: func(ctx context.Context, text string) (string, error) {
				return "", nil
			},
		},
		{
			name: "Non-string attribute",
			attributes: map[string]interface{}{
				"key": 123,
			},
			expectedError: false,
			expectedValue: "",
			mockProcessor: func(ctx context.Context, text string) (string, error) {
				return "", nil
			},
		},
		{
			name: "Failed redaction",
			attributes: map[string]interface{}{
				"key": "sensitive value that fails",
			},
			expectedError: true,
			expectedValue: "sensitive value that fails", // original value should remain
			mockProcessor: func(ctx context.Context, text string) (string, error) {
				return "", fmt.Errorf("mock redaction error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			mockClient := &MockPresidioClient{
				processTextFunc: tt.mockProcessor,
			}

			redaction := &PresidioRedaction{
				Logger: logger,
				Client: mockClient,
			}

			// Create pcommon.Map from test attributes
			attrs := pcommon.NewMap()
			for k, v := range tt.attributes {
				switch val := v.(type) {
				case string:
					attrs.PutStr(k, val)
				case int:
					attrs.PutInt(k, int64(val))
				}
			}

			err := redaction.ProcessAttribute(ctx, attrs)

			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.expectedValue != "" {
				val, _ := attrs.Get("key")
				assert.Equal(t, tt.expectedValue, val.Str())
			}
		})
	}
}

type MockPresidioClient struct {
	processTextFunc func(context.Context, string) (string, error)
}

func (m *MockPresidioClient) ProcessText(ctx context.Context, text string) (string, error) {
	return m.processTextFunc(ctx, text)
}

func TestGetErrorMode(t *testing.T) {
	cfg := &presidioclient.PresidioRedactionProcessorConfig{
		ErrorMode: "test-mode",
	}
	redaction := &PresidioRedaction{Config: cfg}

	assert.Equal(t, cfg.ErrorMode, redaction.GetErrorMode())
}

func TestGetLogger(t *testing.T) {
	logger := zap.NewNop()
	redaction := &PresidioRedaction{Logger: logger}

	assert.Equal(t, logger, redaction.GetLogger())
}
