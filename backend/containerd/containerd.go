// Package containerd provides a ContainerdBackend that executes code via containerd.
// Similar to Docker but more infrastructure-native for servers/agents already using containerd.
package containerd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jonwraymond/toolruntime"
)

// Errors for containerd backend operations.
var (
	// ErrContainerdNotAvailable is returned when containerd is not available.
	ErrContainerdNotAvailable = errors.New("containerd not available")

	// ErrImageNotFound is returned when the execution image is not found.
	ErrImageNotFound = errors.New("image not found")

	// ErrContainerFailed is returned when container creation/execution fails.
	ErrContainerFailed = errors.New("container execution failed")
)

// Logger is the interface for logging.
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// Config configures a ContainerdBackend.
type Config struct {
	// ImageRef is the image reference to use for execution.
	// Default: toolruntime-sandbox:latest
	ImageRef string

	// Namespace is the containerd namespace to use.
	// Default: default
	Namespace string

	// SocketPath is the path to the containerd socket.
	// Default: /run/containerd/containerd.sock
	SocketPath string

	// Logger is an optional logger for backend events.
	Logger Logger
}

// ContainerdBackend executes code via containerd with security isolation.
type ContainerdBackend struct {
	imageRef   string
	namespace  string
	socketPath string
	logger     Logger
}

// New creates a new ContainerdBackend with the given configuration.
func New(cfg Config) *ContainerdBackend {
	imageRef := cfg.ImageRef
	if imageRef == "" {
		imageRef = "toolruntime-sandbox:latest"
	}

	namespace := cfg.Namespace
	if namespace == "" {
		namespace = "default"
	}

	socketPath := cfg.SocketPath
	if socketPath == "" {
		socketPath = "/run/containerd/containerd.sock"
	}

	return &ContainerdBackend{
		imageRef:   imageRef,
		namespace:  namespace,
		socketPath: socketPath,
		logger:     cfg.Logger,
	}
}

// Kind returns the backend kind identifier.
func (b *ContainerdBackend) Kind() toolruntime.BackendKind {
	return toolruntime.BackendContainerd
}

// Execute runs code via containerd with security isolation.
func (b *ContainerdBackend) Execute(ctx context.Context, req toolruntime.ExecuteRequest) (toolruntime.ExecuteResult, error) {
	if err := req.Validate(); err != nil {
		return toolruntime.ExecuteResult{}, err
	}

	start := time.Now()

	result := toolruntime.ExecuteResult{
		Duration: time.Since(start),
		Backend: toolruntime.BackendInfo{
			Kind: toolruntime.BackendContainerd,
			Details: map[string]any{
				"imageRef":  b.imageRef,
				"namespace": b.namespace,
			},
		},
	}

	return result, fmt.Errorf("%w: containerd backend not fully implemented", ErrContainerdNotAvailable)
}

var _ toolruntime.Backend = (*ContainerdBackend)(nil)
