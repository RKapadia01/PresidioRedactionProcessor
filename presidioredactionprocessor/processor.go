// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package presidioredactionprocessor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspan"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"
)

type presidioRedaction struct {
	config             *Config
	logger             *zap.Logger
	client             *http.Client
	concurrencyLimiter chan struct{}
	traceConditions    []ottl.Condition[ottlspan.TransformContext]
	logConditions      []ottl.Condition[ottllog.TransformContext]
}

func newPresidioRedaction(ctx context.Context, cfg *Config, settings component.TelemetrySettings, logger *zap.Logger) *presidioRedaction {
	//Create a span level parser
	parser, err := ottlspan.NewParser(ottlfuncs.StandardFuncs[ottlspan.TransformContext](), settings)
	if err != nil {
		logger.Error("Error creating span parser", zap.Error(err))
		return nil
	}

	//save the trace conditions into memory
	traceConditions := make([]ottl.Condition[ottlspan.TransformContext], 0, len(cfg.TraceConditions))

	// Parse trace conditions
	for _, condition := range cfg.TraceConditions {
		expr, err := parser.ParseCondition(condition)

		if err != nil {
			logger.Error("Error parsing trace condition", zap.Error(err))
			continue
		}

		traceConditions = append(traceConditions, *expr)
	}

	// Parse log conditions
	// logParser, err := ottllog.NewParser(ottlfuncs.StandardFuncs[ottllog.TransformContext](), settings)
	// if err != nil {
	// 	logger.Error("Error creating log parser", zap.Error(err))
	// 	return nil
	// }

	// logConditions := make([]*ottl.Condition[ottllog.TransformContext], 0, len(cfg.TraceConditions))

	// for _, condition := range cfg.LogConditions {
	// 	expr, err := logParser.ParseCondition(condition)

	// 	if err != nil {
	// 		logger.Error("Error parsing log condition", zap.Error(err))
	// 		continue
	// 	}

	// 	logConditions = append(logConditions, expr)
	// }

	return &presidioRedaction{
		config:             cfg,
		logger:             logger,
		client:             &http.Client{},
		concurrencyLimiter: make(chan struct{}, cfg.ConcurrencyLimit),
		traceConditions:    traceConditions,
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

// func (s *presidioRedaction) processLogs(ctx context.Context, logs plog.Logs) (plog.Logs, error) {
// 	for i := 0; i < logs.ResourceLogs().Len(); i++ {
// 		rl := logs.ResourceLogs().At(i)
// 		s.processResourceLog(ctx, rl)
// 	}

// 	return logs, nil
// }

// func (s *presidioRedaction) processResourceLog(ctx context.Context, rl plog.ResourceLogs) {
// 	for j := 0; j < rl.ScopeLogs().Len(); j++ {
// 		ils := rl.ScopeLogs().At(j)
// 		for k := 0; k < ils.LogRecords().Len(); k++ {
// 			log := ils.LogRecords().At(k)

// 			// Check if any condition matches
// 			shouldProcess := false
// 			lCtx := ottllog.NewTransformContext(log, ils, rl)

// 			for _, condition := range s.logConditions {
// 				matches, err := condition.Eval(ctx, lCtx)
// 				if err != nil {
// 					s.logger.Error("Error evaluating log condition", zap.Error(err))
// 					continue
// 				}
// 				if matches {
// 					shouldProcess = true
// 					break
// 				}
// 			}

// 			// Skip if no conditions match
// 			if !shouldProcess && len(s.logConditions) > 0 {
// 				continue
// 			}

// 			s.redactAttr(ctx, log.Attributes())

// 			logBodyStr := log.Body().Str()

// 			if len(logBodyStr) == 0 {
// 				continue
// 			}

// 			redactedBody, err := s.getRedactedValue(ctx, logBodyStr)
// 			if err != nil {
// 				s.logger.Error("Error calling presidio service", zap.Error(err))
// 				continue
// 			}

// 			log.Body().SetStr(redactedBody)
// 		}
// 	}
// }

func (s *presidioRedaction) processResourceSpan(ctx context.Context, rs ptrace.ResourceSpans) {
	rsAttrs := rs.Resource().Attributes()
	s.redactAttr(ctx, rsAttrs)

	for j := 0; j < rs.ScopeSpans().Len(); j++ {
		ils := rs.ScopeSpans().At(j)
		for k := 0; k < ils.Spans().Len(); k++ {
			span := ils.Spans().At(k)

			// Check if any condition matches
			shouldProcess := false
			// tCtx := ottlspan.NewTransformContext(span, ils, rs)
			tCtx := ottlspan.NewTransformContext(
				span,
				ils.Scope(),   // Get scope
				rs.Resource(), // Get resource
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

			// Skip if no conditions match
			if len(s.traceConditions) > 0 && !shouldProcess {
				continue
			}

			spanAttrs := span.Attributes()

			s.redactAttr(ctx, spanAttrs)
		}
	}
}

func (s *presidioRedaction) redactAttr(ctx context.Context, attributes pcommon.Map) {
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

func (s *presidioRedaction) callPresidioAnalyzer(ctx context.Context, value string) ([]PresidioAnalyzerResponse, error) {
	requestPayload := PresidioAnalyzerRequest{
		Text:           value,
		Language:       "en",
		ScoreThreshold: s.config.AnalyzerConfig.ScoreThreshold,
		Entities:       s.config.AnalyzerConfig.Entities,
		Context:        s.config.AnalyzerConfig.Context,
	}

	jsonPayload, err := json.Marshal(requestPayload)
	if err != nil {
		return []PresidioAnalyzerResponse{}, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	url := s.config.AnalyzerEndpoint

	resp, err := s.sendHTTPRequest(ctx, http.MethodPost, url, jsonPayload, headers)
	if err != nil {
		return []PresidioAnalyzerResponse{}, err
	}
	defer resp.Body.Close()

	var presidioAnalyzerResponse []PresidioAnalyzerResponse
	err = json.NewDecoder(resp.Body).Decode(&presidioAnalyzerResponse)
	if err != nil {
		return []PresidioAnalyzerResponse{}, err
	}

	return presidioAnalyzerResponse, nil
}

func (s *presidioRedaction) callPresidioAnonymizer(ctx context.Context, value string, analyzerResults []PresidioAnalyzerResponse) (PresidioAnonymizerResponse, error) {
	anonymizers := make(map[string]PresidioAnonymizer)
	for _, entityAnonymizer := range s.config.AnonymizerConfig.Anonymizers {
		anonymizers[entityAnonymizer.Entity] = PresidioAnonymizer{
			Type:        strings.ToLower(entityAnonymizer.Type),
			NewValue:    entityAnonymizer.NewValue,
			MaskingChar: entityAnonymizer.MaskingChar,
			CharsToMask: entityAnonymizer.CharsToMask,
			FromEnd:     entityAnonymizer.FromEnd,
			HashType:    entityAnonymizer.HashType,
			Key:         entityAnonymizer.Key,
		}
	}

	requestPayload := PresidioAnonymizerRequest{
		Text:            value,
		Anonymizers:     anonymizers,
		AnalyzerResults: analyzerResults,
	}

	jsonPayload, err := json.Marshal(requestPayload)
	if err != nil {
		return PresidioAnonymizerResponse{}, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	url := s.config.AnonymizerEndpoint
	resp, err := s.sendHTTPRequest(ctx, http.MethodPost, url, jsonPayload, headers)
	if err != nil {
		return PresidioAnonymizerResponse{}, err
	}
	defer resp.Body.Close()

	var presidioAnonymizerResponse PresidioAnonymizerResponse
	err = json.NewDecoder(resp.Body).Decode(&presidioAnonymizerResponse)
	if err != nil {
		return PresidioAnonymizerResponse{}, err
	}

	return presidioAnonymizerResponse, nil
}

func (s *presidioRedaction) sendHTTPRequest(ctx context.Context, method, url string, payload []byte, headers map[string]string) (*http.Response, error) {
	// Set a concurrency limiter to avoid overloading the presidio service
	s.concurrencyLimiter <- struct{}{}
	defer func() { <-s.concurrencyLimiter }()

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("service returned status code %d", resp.StatusCode)
	}

	return resp, nil
}

type PresidioAnalyzerRequest struct {
	Text           string   `json:"text"`
	Language       string   `json:"language"`
	ScoreThreshold float64  `json:"score_threshold,omitempty"`
	Entities       []string `json:"entities,omitempty"`
	Context        []string `json:"context,omitempty"`
}

type PresidioAnalyzerResponse struct {
	Start      int     `json:"start"`
	End        int     `json:"end"`
	Score      float64 `json:"score"`
	EntityType string  `json:"entity_type"`
}

type PresidioAnonymizerRequest struct {
	Text            string                        `json:"text,omitempty"`
	Anonymizers     map[string]PresidioAnonymizer `json:"anonymizers,omitempty"`
	AnalyzerResults []PresidioAnalyzerResponse    `json:"analyzer_results,omitempty"`
}

type PresidioAnonymizer struct {
	Type        string `json:"type"`
	NewValue    string `json:"new_value,omitempty"`
	MaskingChar string `json:"masking_char,omitempty"`
	CharsToMask int    `json:"chars_to_mask,omitempty"`
	FromEnd     bool   `json:"from_end,omitempty"`
	HashType    string `json:"hash_type,omitempty"`
	Key         string `json:"key,omitempty"`
}

type PresidioAnonymizerResponse struct {
	Operation  string `json:"operation,omitempty"`
	EntityType string `json:"entity_type"`
	Start      int    `json:"start"`
	End        int    `json:"end"`
	Text       string `json:"text,omitempty"`
}
