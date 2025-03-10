package traces

import (
	"context"
	"fmt"

	"github.com/RKapadia01/PresidioRedactionProcessor/presidioredactionprocessor/internal/common"
	"github.com/RKapadia01/PresidioRedactionProcessor/presidioredactionprocessor/internal/presidioclient"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspan"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type TraceProcessor struct {
	*common.PresidioRedaction
}

func NewPresidioTraceRedaction(ctx context.Context, cfg *presidioclient.PresidioRedactionProcessorConfig, settings component.TelemetrySettings, logger *zap.Logger) *TraceProcessor {
	base, err := common.CreateBaseRedaction(cfg, logger)
	if err != nil {
		logger.Error("Error creating base redaction", zap.Error(err))
		return nil
	}

	if base == nil {
		return nil
	}

	spanParser, err := ottlspan.NewParser(ottlfuncs.StandardFuncs[ottlspan.TransformContext](), settings)
	if err != nil {
		logger.Error("Error creating span parser", zap.Error(err))
		return nil
	}

	base.TraceConditions = common.ParseConditions(cfg.PresidioServiceConfig.TraceConditions, spanParser, logger)
	return &TraceProcessor{PresidioRedaction: base}
}

func (s *TraceProcessor) ProcessTraces(ctx context.Context, batch ptrace.Traces) (ptrace.Traces, error) {
	var errs error
	for i := 0; i < batch.ResourceSpans().Len(); i++ {
		rs := batch.ResourceSpans().At(i)
		resourceSpanErr := s.processResourceSpan(ctx, rs)
		if resourceSpanErr != nil {
			errs = multierr.Append(errs, fmt.Errorf("error processing resource span: %w", resourceSpanErr))
		}
	}

	return batch, common.HandleProcessingError(s, errs, "traces")
}

func (s *TraceProcessor) processResourceSpan(ctx context.Context, rs ptrace.ResourceSpans) error {
	var errs error
	rsAttrs := rs.Resource().Attributes()
	rsAttrsErr := s.ProcessAttribute(ctx, rsAttrs)
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

			for _, condition := range s.TraceConditions {
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

			if len(s.TraceConditions) > 0 && !shouldProcess {
				continue
			}

			spanAttrs := span.Attributes()
			spanAttrErr := s.ProcessAttribute(ctx, spanAttrs)
			if spanAttrErr != nil {
				errs = multierr.Append(errs, fmt.Errorf("error processing span attributes: %w", spanAttrErr))
			}
		}
	}
	return errs
}
