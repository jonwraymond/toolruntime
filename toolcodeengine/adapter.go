// Package toolcodeengine provides an adapter that implements toolcode.Engine
// using toolruntime.Runtime for execution.
package toolcodeengine

import (
	"context"
	"errors"
	"fmt"

	"github.com/jonwraymond/toolcode"
	"github.com/jonwraymond/toolruntime"
)

// Config configures an Engine.
type Config struct {
	// Runtime is the toolruntime.Runtime to use for execution.
	Runtime toolruntime.Runtime

	// Profile is the security profile to use for execution.
	Profile toolruntime.SecurityProfile
}

// Engine implements toolcode.Engine using a toolruntime.Runtime backend.
type Engine struct {
	runtime toolruntime.Runtime
	profile toolruntime.SecurityProfile
}

// New creates a new Engine with the given configuration.
func New(cfg Config) *Engine {
	if cfg.Runtime == nil {
		// Fail fast rather than panic later on Execute.
		panic("toolcodeengine: runtime is required")
	}

	profile := cfg.Profile
	if profile == "" {
		profile = toolruntime.ProfileStandard
	}

	return &Engine{
		runtime: cfg.Runtime,
		profile: profile,
	}
}

// Execute implements toolcode.Engine by delegating to the underlying runtime.
func (e *Engine) Execute(ctx context.Context, params toolcode.ExecuteParams, tools toolcode.Tools) (toolcode.ExecuteResult, error) {
	if e.runtime == nil {
		return toolcode.ExecuteResult{}, toolruntime.ErrRuntimeUnavailable
	}

	// Wrap Tools into a ToolGateway
	gateway := WrapTools(tools)

	// Map toolcode.ExecuteParams to toolruntime.ExecuteRequest
	req := toolruntime.ExecuteRequest{
		Language: params.Language,
		Code:     params.Code,
		Timeout:  params.Timeout,
		Limits: toolruntime.Limits{
			MaxToolCalls: params.MaxToolCalls,
		},
		Profile: e.profile,
		Gateway: gateway,
	}

	// Execute via the runtime
	result, err := e.runtime.Execute(ctx, req)

	// Map errors
	if err != nil {
		return mapResult(result), mapError(err)
	}

	return mapResult(result), nil
}

// mapResult converts toolruntime.ExecuteResult to toolcode.ExecuteResult.
func mapResult(r toolruntime.ExecuteResult) toolcode.ExecuteResult {
	toolCalls := make([]toolcode.ToolCallRecord, len(r.ToolCalls))
	for i, tc := range r.ToolCalls {
		toolCalls[i] = toolcode.ToolCallRecord{
			ToolID:      tc.ToolID,
			BackendKind: tc.BackendKind,
			DurationMs:  tc.Duration.Milliseconds(),
			ErrorOp:     tc.ErrorOp,
		}
	}

	return toolcode.ExecuteResult{
		Value:      r.Value,
		Stdout:     r.Stdout,
		Stderr:     r.Stderr,
		ToolCalls:  toolCalls,
		DurationMs: r.Duration.Milliseconds(),
	}
}

// mapError converts toolruntime errors to toolcode errors.
func mapError(err error) error {
	if err == nil {
		return nil
	}

	// Map timeout and resource limit errors to ErrLimitExceeded
	if errors.Is(err, toolruntime.ErrTimeout) {
		return fmt.Errorf("%w: %v", toolcode.ErrLimitExceeded, err)
	}
	if errors.Is(err, toolruntime.ErrResourceLimit) {
		return fmt.Errorf("%w: %v", toolcode.ErrLimitExceeded, err)
	}

	// Map sandbox violation to ErrCodeExecution
	if errors.Is(err, toolruntime.ErrSandboxViolation) {
		return fmt.Errorf("%w: %v", toolcode.ErrCodeExecution, err)
	}

	// Return other errors as-is
	return err
}
