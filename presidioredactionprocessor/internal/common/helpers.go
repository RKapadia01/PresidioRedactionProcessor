package common

import (
	"strings"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"go.uber.org/zap"
)

type ErrorHandler interface {
	GetErrorMode() ottl.ErrorMode
	GetLogger() *zap.Logger
}

func HandleProcessingError(handler ErrorHandler, err error, operation string) error {
	if err != nil {
		switch handler.GetErrorMode() {
		case ottl.IgnoreError:
			handler.GetLogger().Error("failed to process "+operation, zap.Error(err))
			return nil
		case ottl.PropagateError:
			handler.GetLogger().Error("failed to process "+operation, zap.Error(err))
			return err
		}
	}
	return nil
}

func IsStringHTTPUrl(endpoint string) bool {
	return strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://")
}

func IsStringGRPCUrl(endpoint string) bool {
	return strings.HasPrefix(endpoint, "grpc://") || strings.HasPrefix(endpoint, "grpcs://")
}

func ParseConditions[T any](conditions []string, parser ottl.Parser[T], logger *zap.Logger) []ottl.Condition[T] {
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
