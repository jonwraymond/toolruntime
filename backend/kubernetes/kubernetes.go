// Package kubernetes provides a KubernetesBackend that executes code in Kubernetes pods/jobs.
// Best for scheduling, quotas, and multi-tenant controls; isolation depends on runtime class.
package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jonwraymond/toolruntime"
)

// Errors for Kubernetes backend operations.
var (
	// ErrKubernetesNotAvailable is returned when Kubernetes is not available.
	ErrKubernetesNotAvailable = errors.New("kubernetes not available")

	// ErrPodCreationFailed is returned when pod creation fails.
	ErrPodCreationFailed = errors.New("pod creation failed")

	// ErrPodExecutionFailed is returned when pod execution fails.
	ErrPodExecutionFailed = errors.New("pod execution failed")
)

// Logger is the interface for logging.
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// Config configures a KubernetesBackend.
type Config struct {
	// Namespace is the Kubernetes namespace for execution pods.
	// Default: default
	Namespace string

	// Image is the container image to use for execution.
	// Default: toolruntime-sandbox:latest
	Image string

	// RuntimeClassName is the optional runtime class for stronger isolation.
	// Examples: gvisor, kata
	RuntimeClassName string

	// ServiceAccount is the service account for execution pods.
	ServiceAccount string

	// Logger is an optional logger for backend events.
	Logger Logger
}

// KubernetesBackend executes code in Kubernetes pods/jobs.
type KubernetesBackend struct {
	namespace        string
	image            string
	runtimeClassName string
	serviceAccount   string
	logger           Logger
}

// New creates a new KubernetesBackend with the given configuration.
func New(cfg Config) *KubernetesBackend {
	namespace := cfg.Namespace
	if namespace == "" {
		namespace = "default"
	}

	image := cfg.Image
	if image == "" {
		image = "toolruntime-sandbox:latest"
	}

	return &KubernetesBackend{
		namespace:        namespace,
		image:            image,
		runtimeClassName: cfg.RuntimeClassName,
		serviceAccount:   cfg.ServiceAccount,
		logger:           cfg.Logger,
	}
}

// Kind returns the backend kind identifier.
func (b *KubernetesBackend) Kind() toolruntime.BackendKind {
	return toolruntime.BackendKubernetes
}

// Execute runs code in a Kubernetes pod.
func (b *KubernetesBackend) Execute(ctx context.Context, req toolruntime.ExecuteRequest) (toolruntime.ExecuteResult, error) {
	if err := req.Validate(); err != nil {
		return toolruntime.ExecuteResult{}, err
	}

	start := time.Now()

	result := toolruntime.ExecuteResult{
		Duration: time.Since(start),
		Backend: toolruntime.BackendInfo{
			Kind: toolruntime.BackendKubernetes,
			Details: map[string]any{
				"namespace":        b.namespace,
				"image":            b.image,
				"runtimeClassName": b.runtimeClassName,
			},
		},
	}

	return result, fmt.Errorf("%w: kubernetes backend not fully implemented", ErrKubernetesNotAvailable)
}

var _ toolruntime.Backend = (*KubernetesBackend)(nil)
