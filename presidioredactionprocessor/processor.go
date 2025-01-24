// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package presidioredactionprocessor

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type presidioRedaction struct {
	config             *PresidioRedactionProcessorConfig
	logger             *zap.Logger
	client             *http.Client
	concurrencyLimiter chan struct{}
}

func newPresidioRedaction(_ context.Context, cfg *PresidioRedactionProcessorConfig, logger *zap.Logger) *presidioRedaction {
	if cfg.PresidioServiceConfig.ConcurrencyLimit <= 0 {
		cfg.PresidioServiceConfig.ConcurrencyLimit = 1
	}

	return &presidioRedaction{
		config:             cfg,
		logger:             logger,
		client:             &http.Client{},
		concurrencyLimiter: make(chan struct{}, cfg.PresidioServiceConfig.ConcurrencyLimit),
	}
}

func (s *presidioRedaction) processTraces(ctx context.Context, batch ptrace.Traces) (ptrace.Traces, error) {
	start := time.Now() // Start timer
	defer func() {
		duration := time.Since(start)
		s.logger.Info("processTraces completed", zap.Duration("duration", duration))
	}()

	for i := 0; i < batch.ResourceSpans().Len(); i++ {
		rs := batch.ResourceSpans().At(i)
		s.processResourceSpan(ctx, rs)
	}

	return batch, nil
}

func (s *presidioRedaction) processLogs(ctx context.Context, logs plog.Logs) (plog.Logs, error) {
	for i := 0; i < logs.ResourceLogs().Len(); i++ {
		rl := logs.ResourceLogs().At(i)
		s.processResourceLog(ctx, rl)
	}

	return logs, nil
}

func (s *presidioRedaction) processResourceLog(ctx context.Context, rl plog.ResourceLogs) {
	for j := 0; j < rl.ScopeLogs().Len(); j++ {
		ils := rl.ScopeLogs().At(j)
		for k := 0; k < ils.LogRecords().Len(); k++ {
			log := ils.LogRecords().At(k)

			s.processAttribute(ctx, log.Attributes())

			logBodyStr := log.Body().Str()

			if len(logBodyStr) == 0 {
				continue
			}

			redactedBody, err := s.getRedactedValue(ctx, logBodyStr)
			if err != nil {
				s.logger.Error("Error calling presidio service", zap.Error(err))
				continue
			}

			log.Body().SetStr(redactedBody)
		}
	}
}

func (s *presidioRedaction) processResourceSpan(ctx context.Context, rs ptrace.ResourceSpans) {
	rsAttrs := rs.Resource().Attributes()
	s.processAttribute(ctx, rsAttrs)

	for j := 0; j < rs.ScopeSpans().Len(); j++ {
		ils := rs.ScopeSpans().At(j)
		for k := 0; k < ils.Spans().Len(); k++ {
			span := ils.Spans().At(k)
			spanAttrs := span.Attributes()

			s.processAttribute(ctx, spanAttrs)
		}
	}
}

func (s *presidioRedaction) processAttribute(ctx context.Context, attributes pcommon.Map) {
	attributes.Range(func(k string, v pcommon.Value) bool {
		valueStr := v.Str()
		if len(valueStr) == 0 {
			return true
		}

		redactedValue, err := s.getRedactedValue(ctx, valueStr)

		if err != nil {
			s.logger.Error("Error retrieving the redacted value", zap.Error(err))
			return true
		}
		attributes.PutStr(k, redactedValue)
		return true
	})
}

func (s *presidioRedaction) getRedactedValue(ctx context.Context, value string) (string, error) {
	if isStringGRPCUrl(s.config.PresidioServiceConfig.AnalyzerEndpoint) &&
		isStringGRPCUrl(s.config.PresidioServiceConfig.AnonymizerEndpoint) {
		anonymizerResult, err := s.callPresidioGRPC(ctx, value)
		if err != nil {
			return "", err
		}
		return anonymizerResult.Text, nil
	}
	
	analysisResults, err := s.callPresidioAnalyzer(ctx, value)
	if err != nil {
		return "", err
	}

	if len(analysisResults) == 0 {
		return value, nil
	}

	anonymizerResult, err := s.callPresidioAnonymizer(ctx, value, analysisResults)
	if err != nil {
		return "", err
	}

	return anonymizerResult.Text, nil
}
