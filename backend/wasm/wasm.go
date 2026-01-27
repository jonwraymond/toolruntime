// Package wasm provides a backend that executes code compiled to WebAssembly.
// Provides strong in-process isolation; requires constrained SDK surface.
package wasm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jonwraymond/toolruntime"
)

// Errors for WASM backend operations.
var (
	// ErrWASMRuntimeNotAvailable is returned when WASM runtime is not available.
	ErrWASMRuntimeNotAvailable = errors.New("wasm runtime not available")

	// ErrModuleCompilationFailed is returned when WASM module compilation fails.
	ErrModuleCompilationFailed = errors.New("wasm module compilation failed")

	// ErrModuleExecutionFailed is returned when WASM module execution fails.
	ErrModuleExecutionFailed = errors.New("wasm module execution failed")

	// ErrUnsupportedLanguage is returned when the language cannot be compiled to WASM.
	ErrUnsupportedLanguage = errors.New("language not supported for wasm compilation")
)

// Logger is the interface for logging.
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// Config configures a WASM backend.
type Config struct {
	// Runtime specifies the WASM runtime to use.
	// Options: wasmtime, wasmer, wazero
	// Default: wazero
	Runtime string

	// MaxMemoryPages is the maximum memory pages (64KB each).
	// Default: 256 (16MB)
	MaxMemoryPages int

	// EnableWASI enables WASI (WebAssembly System Interface).
	// Default: true
	EnableWASI bool

	// AllowedHostFunctions lists host functions the WASM module can call.
	AllowedHostFunctions []string

	// Logger is an optional logger for backend events.
	Logger Logger
}

// Backend executes code compiled to WebAssembly.
type Backend struct {
	runtime              string
	maxMemoryPages       int
	enableWASI           bool
	allowedHostFunctions []string
	logger               Logger
}

// New creates a new WASM backend with the given configuration.
func New(cfg Config) *Backend {
	runtime := cfg.Runtime
	if runtime == "" {
		runtime = "wazero"
	}

	maxMemoryPages := cfg.MaxMemoryPages
	if maxMemoryPages <= 0 {
		maxMemoryPages = 256 // 16MB
	}

	return &Backend{
		runtime:              runtime,
		maxMemoryPages:       maxMemoryPages,
		enableWASI:           cfg.EnableWASI,
		allowedHostFunctions: cfg.AllowedHostFunctions,
		logger:               cfg.Logger,
	}
}

// Kind returns the backend kind identifier.
func (b *Backend) Kind() toolruntime.BackendKind {
	return toolruntime.BackendWASM
}

// Execute runs code compiled to WebAssembly.
func (b *Backend) Execute(_ context.Context, req toolruntime.ExecuteRequest) (toolruntime.ExecuteResult, error) {
	if err := req.Validate(); err != nil {
		return toolruntime.ExecuteResult{}, err
	}

	start := time.Now()

	result := toolruntime.ExecuteResult{
		Duration: time.Since(start),
		Backend: toolruntime.BackendInfo{
			Kind: toolruntime.BackendWASM,
			Details: map[string]any{
				"runtime":        b.runtime,
				"maxMemoryPages": b.maxMemoryPages,
				"enableWASI":     b.enableWASI,
			},
		},
	}

	return result, fmt.Errorf("%w: wasm backend not fully implemented", ErrWASMRuntimeNotAvailable)
}

var _ toolruntime.Backend = (*Backend)(nil)
