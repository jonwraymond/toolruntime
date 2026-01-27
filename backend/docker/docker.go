// Package docker provides a backend that executes code in Docker containers
// with configurable security profiles and resource limits.
package docker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jonwraymond/toolruntime"
)

// Errors for Docker backend operations.
var (
	// ErrDockerNotAvailable is returned when Docker is not available.
	ErrDockerNotAvailable = errors.New("docker not available")

	// ErrImageNotFound is returned when the execution image is not found.
	ErrImageNotFound = errors.New("docker image not found")

	// ErrContainerFailed is returned when container creation/execution fails.
	ErrContainerFailed = errors.New("container execution failed")
)

// ContainerOptions represents container configuration for execution.
type ContainerOptions struct {
	// NetworkDisabled disables network access.
	NetworkDisabled bool

	// ReadOnlyRootfs makes the root filesystem read-only.
	ReadOnlyRootfs bool

	// MemoryLimit is the memory limit in bytes.
	MemoryLimit int64

	// CPUQuota is the CPU quota in microseconds.
	CPUQuota int64

	// PidsLimit is the maximum number of processes.
	PidsLimit int64

	// SeccompProfile is the path to a seccomp profile.
	SeccompProfile string

	// User is the user to run as (non-root).
	User string
}

// Config configures a Docker backend.
type Config struct {
	// ImageName is the Docker image to use for execution.
	// Default: toolruntime-sandbox:latest
	ImageName string

	// SeccompPath is the path to a custom seccomp profile for hardened mode.
	SeccompPath string

	// Logger is an optional logger for backend events.
	Logger Logger
}

// Logger is the interface for logging.
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// Backend executes code in Docker containers with security isolation.
type Backend struct {
	imageName   string
	seccompPath string
	logger      Logger
}

// New creates a new Docker backend with the given configuration.
func New(cfg Config) *Backend {
	imageName := cfg.ImageName
	if imageName == "" {
		imageName = "toolruntime-sandbox:latest"
	}

	return &Backend{
		imageName:   imageName,
		seccompPath: cfg.SeccompPath,
		logger:      cfg.Logger,
	}
}

// Kind returns the backend kind identifier.
func (b *Backend) Kind() toolruntime.BackendKind {
	return toolruntime.BackendDocker
}

// Execute runs code in a Docker container with security isolation.
func (b *Backend) Execute(ctx context.Context, req toolruntime.ExecuteRequest) (toolruntime.ExecuteResult, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return toolruntime.ExecuteResult{}, err
	}

	// Apply timeout
	timeout := req.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return toolruntime.ExecuteResult{}, ctx.Err()
	default:
	}

	start := time.Now()

	// Get container options based on profile and limits
	profile := req.Profile
	if profile == "" {
		profile = toolruntime.ProfileStandard
	}
	opts := b.containerOptions(profile, req.Limits)

	// Log execution
	if b.logger != nil {
		b.logger.Info("executing in Docker container",
			"profile", profile,
			"image", b.imageName,
			"networkDisabled", opts.NetworkDisabled,
			"readOnlyRootfs", opts.ReadOnlyRootfs)
	}

	// NOTE: Full Docker integration would:
	// 1. Create container with opts
	// 2. Copy code and gateway proxy server into container
	// 3. Start container and wait for completion
	// 4. Capture stdout/stderr
	// 5. Remove container

	// For now, return an error indicating Docker is not implemented
	result := toolruntime.ExecuteResult{
		Duration: time.Since(start),
		Backend: toolruntime.BackendInfo{
			Kind: toolruntime.BackendDocker,
			Details: map[string]any{
				"image":   b.imageName,
				"profile": string(profile),
			},
		},
	}

	return result, fmt.Errorf("%w: Docker backend not fully implemented", ErrDockerNotAvailable)
}

// containerOptions returns ContainerOptions based on the security profile and limits.
func (b *Backend) containerOptions(profile toolruntime.SecurityProfile, limits toolruntime.Limits) ContainerOptions {
	opts := ContainerOptions{
		User: "nobody:nogroup", // Always run as non-root
	}

	switch profile {
	case toolruntime.ProfileDev:
		// Dev mode: minimal restrictions
		opts.NetworkDisabled = false
		opts.ReadOnlyRootfs = false

	case toolruntime.ProfileStandard:
		// Standard: no network, read-only rootfs
		opts.NetworkDisabled = true
		opts.ReadOnlyRootfs = true

	case toolruntime.ProfileHardened:
		// Hardened: all restrictions plus seccomp
		opts.NetworkDisabled = true
		opts.ReadOnlyRootfs = true
		if b.seccompPath != "" {
			opts.SeccompProfile = b.seccompPath
		}
	}

	// Apply resource limits
	if limits.MemoryBytes > 0 {
		opts.MemoryLimit = limits.MemoryBytes
	}
	if limits.CPUQuotaMillis > 0 {
		// Convert milliseconds to microseconds for Docker
		opts.CPUQuota = limits.CPUQuotaMillis * 1000
	}
	if limits.PidsMax > 0 {
		opts.PidsLimit = limits.PidsMax
	}

	return opts
}
