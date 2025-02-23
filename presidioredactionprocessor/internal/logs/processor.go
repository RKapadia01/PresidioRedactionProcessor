package logs

import (
	"context"
	"fmt"

	"github.com/RKapadia01/presidioredactionprocessor/internal/common"
	"github.com/RKapadia01/presidioredactionprocessor/internal/presidioclient"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type LogProcessor struct {
	*common.PresidioRedaction
}

func NewPresidioLogRedaction(ctx context.Context, cfg *presidioclient.PresidioRedactionProcessorConfig, settings component.TelemetrySettings, logger *zap.Logger) *LogProcessor {
	base, err := common.CreateBaseRedaction(cfg, logger)
	if err != nil {
		logger.Error("Error creating base redaction", zap.Error(err))
		return nil
	}

	if base == nil {
		return nil
	}

	logParser, err := ottllog.NewParser(ottlfuncs.StandardFuncs[ottllog.TransformContext](), settings)
	if err != nil {
		logger.Error("Error creating log parser", zap.Error(err))
		return nil
	}

	base.LogConditions = common.ParseConditions(cfg.PresidioServiceConfig.LogConditions, logParser, logger)
	return &LogProcessor{PresidioRedaction: base}
}

func (s *LogProcessor) ProcessLogs(ctx context.Context, logs plog.Logs) (plog.Logs, error) {
	var errs error
	for i := 0; i < logs.ResourceLogs().Len(); i++ {
		rl := logs.ResourceLogs().At(i)
		resourceLogErr := s.processResourceLog(ctx, rl)
		if resourceLogErr != nil {
			errs = multierr.Append(errs, fmt.Errorf("error processing resource log: %w", resourceLogErr))
		}
	}

	return logs, common.HandleProcessingError(s, errs, "logs")
}

func (s *LogProcessor) processResourceLog(ctx context.Context, rl plog.ResourceLogs) error {
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

			for _, condition := range s.LogConditions {
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

			if !shouldProcess && len(s.LogConditions) > 0 {
				continue
			}

			attrErr := s.ProcessAttribute(ctx, log.Attributes())
			if attrErr != nil {
				errs = multierr.Append(errs, fmt.Errorf("error processing log attributes: %w", attrErr))
			}

			logBodyStr := log.Body().Str()

			if len(logBodyStr) == 0 {
				continue
			}

			redactedBody, redactionErr := s.Client.ProcessText(ctx, logBodyStr)
			if redactionErr != nil {
				errs = multierr.Append(errs, fmt.Errorf("error redacting log body: %w", redactionErr))
				continue
			}

			log.Body().SetStr(redactedBody)
		}
	}
	return errs
}
