// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package presidioredactionprocessor

import (
	"context"
	"net/http"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspan"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"
)

type presidioRedaction struct {
	config             *PresidioRedactionProcessorConfig
	logger             *zap.Logger
	client             *http.Client
	concurrencyLimiter chan struct{}
	traceConditions    []ottl.Condition[ottlspan.TransformContext]
	logConditions      []ottl.Condition[ottllog.TransformContext]
}

func newPresidioLogRedaction(ctx context.Context, cfg *PresidioRedactionProcessorConfig, settings component.TelemetrySettings, logger *zap.Logger) *presidioRedaction {
	logParser, err := ottllog.NewParser(ottlfuncs.StandardFuncs[ottllog.TransformContext](), settings)
	if err != nil {
		logger.Error("Error creating log parser", zap.Error(err))
		return nil
	}

	logConditions := make([]ottl.Condition[ottllog.TransformContext], 0, len(cfg.PresidioServiceConfig.TraceConditions))

	for _, condition := range cfg.PresidioServiceConfig.LogConditions {
		expr, err := logParser.ParseCondition(condition)

		if err != nil {
			logger.Error("Error parsing log condition", zap.Error(err))
			continue
		}

		logConditions = append(logConditions, *expr)
	}

	return &presidioRedaction{
		config:             cfg,
		logger:             logger,
		client:             &http.Client{},
		concurrencyLimiter: make(chan struct{}, cfg.PresidioServiceConfig.ConcurrencyLimit),
		logConditions:      logConditions,
	}
}

func newPresidioTraceRedaction(ctx context.Context, cfg *PresidioRedactionProcessorConfig, settings component.TelemetrySettings, logger *zap.Logger) *presidioRedaction {
	parser, err := ottlspan.NewParser(ottlfuncs.StandardFuncs[ottlspan.TransformContext](), settings)
	if err != nil {
		logger.Error("Error creating span parser", zap.Error(err))
		return nil
	}

	traceConditions := make([]ottl.Condition[ottlspan.TransformContext], 0, len(cfg.PresidioServiceConfig.TraceConditions))

	for _, condition := range cfg.PresidioServiceConfig.TraceConditions {
		expr, err := parser.ParseCondition(condition)

		if err != nil {
			logger.Error("Error parsing trace condition", zap.Error(err))
			continue
		}

		traceConditions = append(traceConditions, *expr)
	}

	return &presidioRedaction{
		config:             cfg,
		logger:             logger,
		client:             &http.Client{},
		concurrencyLimiter: make(chan struct{}, cfg.PresidioServiceConfig.ConcurrencyLimit),
		traceConditions:    traceConditions,
	}
}

func (s *presidioRedaction) processTraces(ctx context.Context, batch ptrace.Traces) (ptrace.Traces, error) {
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

			shouldProcess := false
			lCtx := ottllog.NewTransformContext(
				log,
				ils.Scope(),
				rl.Resource(),
				ils,
				rl,
			)

			for _, condition := range s.logConditions {
				matches, err := condition.Eval(ctx, lCtx)
				if err != nil {
					s.logger.Error("Error evaluating log condition", zap.Error(err))
					continue
				}
				if matches {
					shouldProcess = true
					break
				}
			}

			if !shouldProcess && len(s.logConditions) > 0 {
				continue
			}

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
			shouldProcess := false
			tCtx := ottlspan.NewTransformContext(
				span,
				ils.Scope(),
				rs.Resource(),
				ils,
				rs,
			)

			for _, condition := range s.traceConditions {
				matches, err := condition.Eval(ctx, tCtx)
				if err != nil {
					s.logger.Error("Error evaluating trace condition", zap.Error(err))
					continue
				}
				if matches {
					shouldProcess = true
					break
				}
			}

			if len(s.traceConditions) > 0 && !shouldProcess {
				continue
			}

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
