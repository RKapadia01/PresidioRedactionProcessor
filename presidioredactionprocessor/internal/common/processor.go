package common

import (
	"context"
	"fmt"

	"github.com/RKapadia01/PresidioRedactionProcessor/presidioredactionprocessor/internal/presidioclient"
	"github.com/RKapadia01/PresidioRedactionProcessor/presidioredactionprocessor/internal/presidioclient/grpcclient"
	"github.com/RKapadia01/PresidioRedactionProcessor/presidioredactionprocessor/internal/presidioclient/httpclient"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspan"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type PresidioRedaction struct {
	Config             *presidioclient.PresidioRedactionProcessorConfig
	Logger             *zap.Logger
	Client             presidioclient.PresidioClient
	ConcurrencyLimiter chan struct{}
	TraceConditions    []ottl.Condition[ottlspan.TransformContext]
	LogConditions      []ottl.Condition[ottllog.TransformContext]
}

func CreateBaseRedaction(cfg *presidioclient.PresidioRedactionProcessorConfig, logger *zap.Logger) (*PresidioRedaction, error) {
	var client presidioclient.PresidioClient

	if IsStringGRPCUrl(cfg.PresidioServiceConfig.AnalyzerEndpoint) &&
		IsStringGRPCUrl(cfg.PresidioServiceConfig.AnonymizerEndpoint) {
		client = grpcclient.NewPresidioGrpcClient(cfg)
	} else if IsStringHTTPUrl(cfg.PresidioServiceConfig.AnalyzerEndpoint) &&
		IsStringHTTPUrl(cfg.PresidioServiceConfig.AnonymizerEndpoint) {
		client = httpclient.NewPresidioHttpClient(cfg, logger)
	} else {
		return nil, fmt.Errorf("invalid Presidio service endpoints")
	}

	return &PresidioRedaction{
		Config:             cfg,
		Logger:             logger,
		Client:             client,
		ConcurrencyLimiter: make(chan struct{}, cfg.PresidioServiceConfig.ConcurrencyLimit),
	}, nil
}

func (s *PresidioRedaction) ProcessAttribute(ctx context.Context, attributes pcommon.Map) error {
	var errs error
	attributes.Range(func(k string, v pcommon.Value) bool {
		valueStr := v.Str()
		if len(valueStr) == 0 {
			return true
		}

		redactedValue, redactionErr := s.Client.ProcessText(ctx, valueStr)

		if redactionErr != nil {
			errs = multierr.Append(errs, fmt.Errorf("error redacting value: %w", redactionErr))
			return true
		}
		attributes.PutStr(k, redactedValue)
		return true
	})

	return errs
}

func (s *PresidioRedaction) GetErrorMode() ottl.ErrorMode {
	return s.Config.ErrorMode
}

func (s *PresidioRedaction) GetLogger() *zap.Logger {
	return s.Logger
}
