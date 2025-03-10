package logs

import (
	"context"
	"testing"

	"github.com/RKapadia01/PresidioRedactionProcessor/presidioredactionprocessor/internal/common"
	"github.com/RKapadia01/PresidioRedactionProcessor/presidioredactionprocessor/internal/presidioclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type mockPresidioClient struct {
	processTextFunc func(ctx context.Context, text string) (string, error)
}

func (m *mockPresidioClient) ProcessText(ctx context.Context, text string) (string, error) {
	return m.processTextFunc(ctx, text)
}

func createTestProcessor(client presidioclient.PresidioClient) *LogProcessor {
	return &LogProcessor{
		PresidioRedaction: &common.PresidioRedaction{
			Client: client,
			Logger: zap.NewNop(),
		},
	}
}

func createTestLogs(bodyContent string, attributes map[string]interface{}) plog.Logs {
	logs := plog.NewLogs()
	rl := logs.ResourceLogs().AppendEmpty()
	sl := rl.ScopeLogs().AppendEmpty()
	lr := sl.LogRecords().AppendEmpty()

	lr.Body().SetStr(bodyContent)

	for k, v := range attributes {
		setAttributeValue(lr.Attributes(), k, v)
	}

	return logs
}

func setAttributeValue(attrs pcommon.Map, key string, value interface{}) {
	switch v := value.(type) {
	case string:
		attrs.PutStr(key, v)
	case int64:
		attrs.PutInt(key, v)
	case float64:
		attrs.PutDouble(key, v)
	case bool:
		attrs.PutBool(key, v)
	}
}

func TestLogProcessor_ProcessLogs(t *testing.T) {
	tests := []struct {
		name         string
		inputLogs    plog.Logs
		mockResponse string
		expectedBody string
		expectError  bool
	}{
		{
			name:         "successful redaction",
			inputLogs:    createTestLogs("Hello John Doe", nil),
			mockResponse: "Hello <REDACTED>",
			expectedBody: "Hello <REDACTED>",
			expectError:  false,
		},
		{
			name:         "empty log body",
			inputLogs:    createTestLogs("", nil),
			mockResponse: "",
			expectedBody: "",
			expectError:  false,
		},
		{
			name: "multiple log records",
			inputLogs: func() plog.Logs {
				logs := plog.NewLogs()
				rl := logs.ResourceLogs().AppendEmpty()
				sl := rl.ScopeLogs().AppendEmpty()

				lr1 := sl.LogRecords().AppendEmpty()
				lr1.Body().SetStr("Hello John")

				lr2 := sl.LogRecords().AppendEmpty()
				lr2.Body().SetStr("Hi Jane")

				return logs
			}(),
			mockResponse: "Hello <REDACTED>",
			expectedBody: "Hello <REDACTED>",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockPresidioClient{
				processTextFunc: func(ctx context.Context, text string) (string, error) {
					return tt.mockResponse, nil
				},
			}

			processor := createTestProcessor(mockClient)
			processedLogs, err := processor.ProcessLogs(context.Background(), tt.inputLogs)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, processedLogs)

			// Verify the first log record's body
			if processedLogs.ResourceLogs().Len() > 0 {
				rl := processedLogs.ResourceLogs().At(0)
				if rl.ScopeLogs().Len() > 0 {
					sl := rl.ScopeLogs().At(0)
					if sl.LogRecords().Len() > 0 {
						lr := sl.LogRecords().At(0)
						assert.Equal(t, tt.expectedBody, lr.Body().Str())
					}
				}
			}
		})
	}
}

func TestNewPresidioLogRedaction(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *presidioclient.PresidioRedactionProcessorConfig
		wantNil bool
	}{
		{
			name:    "valid empty configuration",
			cfg:     &presidioclient.PresidioRedactionProcessorConfig{},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zap.NewNop()
			processor := NewPresidioLogRedaction(
				context.Background(),
				tt.cfg,
				component.TelemetrySettings{},
				logger,
			)

			if tt.wantNil {
				assert.Nil(t, processor, "processor should be nil")
			} else {
				assert.NotNil(t, processor, "processor should not be nil")
				assert.NotNil(t, processor.Client, "client should not be nil")
				assert.NotNil(t, processor.Logger, "logger should not be nil")
			}
		})
	}
}
