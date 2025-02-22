// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package presidioredactionprocessor

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/multierr"
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

func createBaseRedaction(cfg *PresidioRedactionProcessorConfig, logger *zap.Logger) *presidioRedaction {
	return &presidioRedaction{
		config:             cfg,
		logger:             logger,
		client:             &http.Client{},
		concurrencyLimiter: make(chan struct{}, cfg.PresidioServiceConfig.ConcurrencyLimit),
	}
}

func newPresidioLogRedaction(ctx context.Context, cfg *PresidioRedactionProcessorConfig, settings component.TelemetrySettings, logger *zap.Logger) *presidioRedaction {
	base := createBaseRedaction(cfg, logger)
	if base == nil {
		return nil
	}

	logParser, err := ottllog.NewParser(ottlfuncs.StandardFuncs[ottllog.TransformContext](), settings)
	if err != nil {
		logger.Error("Error creating log parser", zap.Error(err))
		return nil
	}

	base.logConditions = parseConditions(cfg.PresidioServiceConfig.LogConditions, logParser, logger)
	return base
}

func newPresidioTraceRedaction(ctx context.Context, cfg *PresidioRedactionProcessorConfig, settings component.TelemetrySettings, logger *zap.Logger) *presidioRedaction {
	base := createBaseRedaction(cfg, logger)
	if base == nil {
		return nil
	}

	spanParser, err := ottlspan.NewParser(ottlfuncs.StandardFuncs[ottlspan.TransformContext](), settings)
	if err != nil {
		logger.Error("Error creating span parser", zap.Error(err))
		return nil
	}

	base.traceConditions = parseConditions(cfg.PresidioServiceConfig.TraceConditions, spanParser, logger)
	return base
}

func (s *presidioRedaction) processTraces(ctx context.Context, batch ptrace.Traces) (ptrace.Traces, error) {
	var errs error
	for i := 0; i < batch.ResourceSpans().Len(); i++ {
		rs := batch.ResourceSpans().At(i)
		resourceSpanErr := s.processResourceSpan(ctx, rs)
		if resourceSpanErr != nil {
			errs = multierr.Append(errs, fmt.Errorf("error processing resource span: %w", resourceSpanErr))
		}
	}

	return batch, s.handleProcessingError(errs, "traces")
}

func (s *presidioRedaction) processLogs(ctx context.Context, logs plog.Logs) (plog.Logs, error) {
	var errs error
	for i := 0; i < logs.ResourceLogs().Len(); i++ {
		rl := logs.ResourceLogs().At(i)
		resourceLogErr := s.processResourceLog(ctx, rl)
		if resourceLogErr != nil {
			errs = multierr.Append(errs, fmt.Errorf("error processing resource log: %w", resourceLogErr))
		}
	}

	return logs, s.handleProcessingError(errs, "logs")
}

func (s *presidioRedaction) processResourceLog(ctx context.Context, rl plog.ResourceLogs) error {
	var errs error
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
				matches, logCdnErr := condition.Eval(ctx, lCtx)
				if logCdnErr != nil {
					errs = multierr.Append(errs, fmt.Errorf("error evaluating log condition: %w", logCdnErr))
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

			attrErr := s.processAttribute(ctx, log.Attributes())
			if attrErr != nil {
				errs = multierr.Append(errs, fmt.Errorf("error processing log attributes: %w", attrErr))
			}

			logBodyStr := log.Body().Str()

			if len(logBodyStr) == 0 {
				continue
			}

			redactedBody, redactionErr := s.getRedactedValue(ctx, logBodyStr)
			if redactionErr != nil {
				errs = multierr.Append(errs, fmt.Errorf("error redacting log body: %w", redactionErr))
				continue
			}

			log.Body().SetStr(redactedBody)
		}
	}
	return errs
}

func (s *presidioRedaction) processResourceSpan(ctx context.Context, rs ptrace.ResourceSpans) error {
	var errs error
	rsAttrs := rs.Resource().Attributes()
	rsAttrsErr := s.processAttribute(ctx, rsAttrs)
	if rsAttrsErr != nil {
		errs = multierr.Append(errs, fmt.Errorf("error processing resource attributes: %w", rsAttrsErr))
	}

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
				matches, traceCdnErr := condition.Eval(ctx, tCtx)
				if traceCdnErr != nil {
					errs = multierr.Append(errs, fmt.Errorf("error evaluating trace condition: %w", traceCdnErr))
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
			spanAttrErr := s.processAttribute(ctx, spanAttrs)
			if spanAttrErr != nil {
				errs = multierr.Append(errs, fmt.Errorf("error processing span attributes: %w", spanAttrErr))
			}
		}
	}
	return errs
}

func (s *presidioRedaction) processAttribute(ctx context.Context, attributes pcommon.Map) error {
	var errs error
	attributes.Range(func(k string, v pcommon.Value) bool {
		valueStr := v.Str()
		if len(valueStr) == 0 {
			return true
		}

		redactedValue, redactionErr := s.getRedactedValue(ctx, valueStr)

		if redactionErr != nil {
			errs = multierr.Append(errs, fmt.Errorf("error redacting value: %w", redactionErr))
			return true
		}
		attributes.PutStr(k, redactedValue)
		return true
	})

	return errs
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

func parseConditions[T any](conditions []string, parser ottl.Parser[T], logger *zap.Logger) []ottl.Condition[T] {
	parsed := make([]ottl.Condition[T], 0, len(conditions))
	for _, condition := range conditions {
		expr, err := parser.ParseCondition(condition)
		if err != nil {
			logger.Error("Error parsing condition", zap.Error(err))
			continue
		}
		parsed = append(parsed, *expr)
	}
	return parsed
}

func (s *presidioRedaction) handleProcessingError(err error, operation string) error {
	if err != nil {
		switch s.config.ErrorMode {
		case ottl.IgnoreError:
			s.logger.Error("failed to process "+operation, zap.Error(err))
			return nil
		case ottl.PropagateError:
			s.logger.Error("failed to process "+operation, zap.Error(err))
			return err
		}
	}
	return nil
}
