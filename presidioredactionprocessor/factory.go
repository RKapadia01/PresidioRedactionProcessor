// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:generate mdatagen metadata.yaml

package presidioredactionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/redactionprocessor"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"

	"github.com/RKapadia01/PresidioRedactionProcessor/internal/logs"
	"github.com/RKapadia01/PresidioRedactionProcessor/internal/metadata"
	"github.com/RKapadia01/PresidioRedactionProcessor/internal/presidioclient"
	"github.com/RKapadia01/PresidioRedactionProcessor/internal/traces"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
)

func NewFactory() processor.Factory {
	return processor.NewFactory(
		metadata.Type,
		createDefaultConfig,
		processor.WithTraces(createTracesProcessor, metadata.TracesStability),
		processor.WithLogs(createLogsProcessor, metadata.LogsStability),
	)
}

func createDefaultConfig() component.Config {
	return &presidioclient.PresidioRedactionProcessorConfig{
		PresidioRunMode: "embedded",
		ErrorMode:       ottl.PropagateError,
		PresidioServiceConfig: presidioclient.PresidioServiceConfig{
			AnalyzerEndpoint:   "grpc://localhost:50051",
			AnonymizerEndpoint: "grpc://localhost:50052",
		},
	}
}

func createTracesProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	next consumer.Traces,
) (processor.Traces, error) {
	oCfg := cfg.(*presidioclient.PresidioRedactionProcessorConfig)

	//TODO: Fix this
	// if err := oCfg.validate(); err != nil {
	// 	return nil, err
	// }

	configurePresidioEndpoints(oCfg)

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

	// TODO: Fix this
	// if err := oCfg.validate(); err != nil {
	// 	return nil, err
	// }

	configurePresidioEndpoints(oCfg)

	presidioRedaction := logs.NewPresidioLogRedaction(ctx, oCfg, set.TelemetrySettings, set.Logger)

	return processorhelper.NewLogs(
		ctx,
		set,
		cfg,
		next,
		presidioRedaction.ProcessLogs,
		processorhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}))
}

func configurePresidioEndpoints(cfg *presidioclient.PresidioRedactionProcessorConfig) {
	if cfg.PresidioRunMode == "embedded" {
		cfg.PresidioServiceConfig.AnalyzerEndpoint = "grpc://localhost:50051"
		cfg.PresidioServiceConfig.AnonymizerEndpoint = "grpc://localhost:500052"
	}
}
