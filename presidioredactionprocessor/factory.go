// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:generate mdatagen metadata.yaml

package presidioredactionprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"

	"github.com/RKapadia01/PresidioRedactionProcessor/presidioredactionprocessor/internal/logs"
	"github.com/RKapadia01/PresidioRedactionProcessor/presidioredactionprocessor/internal/metadata"
	"github.com/RKapadia01/PresidioRedactionProcessor/presidioredactionprocessor/internal/presidioclient"
	"github.com/RKapadia01/PresidioRedactionProcessor/presidioredactionprocessor/internal/traces"
)

func NewFactory() processor.Factory {
	return processor.NewFactory(
		metadata.Type,
		presidioclient.CreateDefaultConfig,
		processor.WithTraces(createTracesProcessor, metadata.TracesStability),
		processor.WithLogs(createLogsProcessor, metadata.LogsStability),
	)
}

func createTracesProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	next consumer.Traces,
) (processor.Traces, error) {
	oCfg := cfg.(*presidioclient.PresidioRedactionProcessorConfig)

	if err := oCfg.Validate(); err != nil {
		return nil, err
	}

	presidioRedaction := traces.NewPresidioTraceRedaction(ctx, oCfg, set.TelemetrySettings, set.Logger)

	return processorhelper.NewTraces(
		ctx,
		set,
		cfg,
		next,
		presidioRedaction.ProcessTraces,
		processorhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}))
}

func createLogsProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	next consumer.Logs,
) (processor.Logs, error) {
	oCfg := cfg.(*presidioclient.PresidioRedactionProcessorConfig)

	if err := oCfg.Validate(); err != nil {
		return nil, err
	}

	presidioRedaction := logs.NewPresidioLogRedaction(ctx, oCfg, set.TelemetrySettings, set.Logger)

	return processorhelper.NewLogs(
		ctx,
		set,
		cfg,
		next,
		presidioRedaction.ProcessLogs,
		processorhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}))
}
