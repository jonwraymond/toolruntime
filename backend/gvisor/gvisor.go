// Package gvisor provides a GVisorBackend that executes code with gVisor (runsc).
// Provides stronger isolation than plain containers; appropriate for untrusted multi-tenant execution.
package gvisor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jonwraymond/toolruntime"
)

// Errors for gVisor backend operations.
var (
	// ErrGVisorNotAvailable is returned when gVisor/runsc is not available.
	ErrGVisorNotAvailable = errors.New("gvisor not available")

	// ErrSandboxCreationFailed is returned when sandbox creation fails.
	ErrSandboxCreationFailed = errors.New("sandbox creation failed")

	// ErrSandboxExecutionFailed is returned when sandbox execution fails.
	ErrSandboxExecutionFailed = errors.New("sandbox execution failed")
)

// Logger is the interface for logging.
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// Config configures a GVisorBackend.
type Config struct {
	// RunscPath is the path to the runsc binary.
	// Default: runsc (uses PATH)
	RunscPath string

	// RootDir is the root directory for gVisor state.
	// Default: /var/run/gvisor
	RootDir string

	// Platform is the gVisor platform to use.
	// Options: ptrace, kvm, systrap
	// Default: systrap
	Platform string

	// NetworkMode specifies the network configuration.
	// Options: none, sandbox, host
	// Default: none
	NetworkMode string

	// Logger is an optional logger for backend events.
	Logger Logger
}

// GVisorBackend executes code with gVisor for stronger isolation.
type GVisorBackend struct {
	runscPath   string
	rootDir     string
	platform    string
	networkMode string
	logger      Logger
}

// New creates a new GVisorBackend with the given configuration.
func New(cfg Config) *GVisorBackend {
	runscPath := cfg.RunscPath
	if runscPath == "" {
		runscPath = "runsc"
	}

	rootDir := cfg.RootDir
	if rootDir == "" {
		rootDir = "/var/run/gvisor"
	}

	platform := cfg.Platform
	if platform == "" {
		platform = "systrap"
	}

	networkMode := cfg.NetworkMode
	if networkMode == "" {
		networkMode = "none"
	}

	return &GVisorBackend{
		runscPath:   runscPath,
		rootDir:     rootDir,
		platform:    platform,
		networkMode: networkMode,
		logger:      cfg.Logger,
	}
}

// Kind returns the backend kind identifier.
func (b *GVisorBackend) Kind() toolruntime.BackendKind {
	return toolruntime.BackendGVisor
}

// Execute runs code with gVisor isolation.
func (b *GVisorBackend) Execute(ctx context.Context, req toolruntime.ExecuteRequest) (toolruntime.ExecuteResult, error) {
	if err := req.Validate(); err != nil {
		return toolruntime.ExecuteResult{}, err
	}

	start := time.Now()

	result := toolruntime.ExecuteResult{
		Duration: time.Since(start),
		Backend: toolruntime.BackendInfo{
			Kind: toolruntime.BackendGVisor,
			Details: map[string]any{
				"platform":    b.platform,
				"networkMode": b.networkMode,
			},
		},
	}

	return result, fmt.Errorf("%w: gvisor backend not fully implemented", ErrGVisorNotAvailable)
}

var _ toolruntime.Backend = (*GVisorBackend)(nil)
