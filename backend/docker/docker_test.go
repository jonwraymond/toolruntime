package docker

import (
	"context"
	"errors"
	"testing"

	"github.com/jonwraymond/tooldocs"
	"github.com/jonwraymond/toolindex"
	"github.com/jonwraymond/toolrun"
	"github.com/jonwraymond/toolruntime"
)

// mockGateway implements toolruntime.ToolGateway for testing
type mockGateway struct{}

func (m *mockGateway) SearchTools(_ context.Context, _ string, _ int) ([]toolindex.Summary, error) {
	return nil, nil
}

func (m *mockGateway) ListNamespaces(_ context.Context) ([]string, error) {
	return nil, nil
}

func (m *mockGateway) DescribeTool(_ context.Context, _ string, _ tooldocs.DetailLevel) (tooldocs.ToolDoc, error) {
	return tooldocs.ToolDoc{}, nil
}

func (m *mockGateway) ListToolExamples(_ context.Context, _ string, _ int) ([]tooldocs.ToolExample, error) {
	return nil, nil
}

func (m *mockGateway) RunTool(_ context.Context, _ string, _ map[string]any) (toolrun.RunResult, error) {
	return toolrun.RunResult{}, nil
}

func (m *mockGateway) RunChain(_ context.Context, _ []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error) {
	return toolrun.RunResult{}, nil, nil
}

// TestBackendImplementsInterface verifies Backend satisfies toolruntime.Backend
func TestBackendImplementsInterface(t *testing.T) {
	t.Helper()
	var _ toolruntime.Backend = (*Backend)(nil)
}

func TestBackendKind(t *testing.T) {
	b := New(Config{})

	if b.Kind() != toolruntime.BackendDocker {
		t.Errorf("Kind() = %v, want %v", b.Kind(), toolruntime.BackendDocker)
	}
}

func TestBackendRequiresGateway(t *testing.T) {
	b := New(Config{})

	ctx := context.Background()
	req := toolruntime.ExecuteRequest{
		Code:    "print('hello')",
		Gateway: nil,
	}

	_, err := b.Execute(ctx, req)
	if !errors.Is(err, toolruntime.ErrMissingGateway) {
		t.Errorf("Execute() without gateway error = %v, want %v", err, toolruntime.ErrMissingGateway)
	}
}

func TestBackendRequiresCode(t *testing.T) {
	b := New(Config{})

	ctx := context.Background()
	req := toolruntime.ExecuteRequest{
		Code:    "",
		Gateway: &mockGateway{},
	}

	_, err := b.Execute(ctx, req)
	if !errors.Is(err, toolruntime.ErrMissingCode) {
		t.Errorf("Execute() without code error = %v, want %v", err, toolruntime.ErrMissingCode)
	}
}

func TestBackendContractCompliance(t *testing.T) {
	toolruntime.RunBackendContractTests(t, toolruntime.BackendContract{
		NewBackend: func() toolruntime.Backend {
			return New(Config{})
		},
		NewGateway: func() toolruntime.ToolGateway {
			return &mockGateway{}
		},
		ExpectedKind:       toolruntime.BackendDocker,
		SkipExecutionTests: true, // Docker may not be available
	})
}

func TestBackendProfileRestrictions(t *testing.T) {
	b := New(Config{})

	tests := []struct {
		profile         toolruntime.SecurityProfile
		expectNetwork   bool
		expectReadOnly  bool
	}{
		{toolruntime.ProfileStandard, false, true},
		{toolruntime.ProfileHardened, false, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.profile), func(t *testing.T) {
			// This is a design validation test - actual Docker restrictions
			// would be tested in integration tests
			opts := b.containerOptions(tt.profile, toolruntime.Limits{})

			if opts.NetworkDisabled != !tt.expectNetwork {
				t.Errorf("profile %v NetworkDisabled = %v, want %v",
					tt.profile, opts.NetworkDisabled, !tt.expectNetwork)
			}
			if opts.ReadOnlyRootfs != tt.expectReadOnly {
				t.Errorf("profile %v ReadOnlyRootfs = %v, want %v",
					tt.profile, opts.ReadOnlyRootfs, tt.expectReadOnly)
			}
		})
	}
}

func TestBackendResourceLimits(t *testing.T) {
	b := New(Config{})

	limits := toolruntime.Limits{
		MemoryBytes:    256 * 1024 * 1024, // 256MB
		CPUQuotaMillis: 1000,              // 1 CPU
		PidsMax:        100,
	}

	opts := b.containerOptions(toolruntime.ProfileStandard, limits)

	if opts.MemoryLimit != limits.MemoryBytes {
		t.Errorf("MemoryLimit = %d, want %d", opts.MemoryLimit, limits.MemoryBytes)
	}
	if opts.CPUQuota != limits.CPUQuotaMillis*1000 {
		t.Errorf("CPUQuota = %d, want %d", opts.CPUQuota, limits.CPUQuotaMillis*1000)
	}
	if opts.PidsLimit != limits.PidsMax {
		t.Errorf("PidsLimit = %d, want %d", opts.PidsLimit, limits.PidsMax)
	}
}

func TestBackendRequiresClient(t *testing.T) {
	b := New(Config{}) // No client configured

	ctx := context.Background()
	req := toolruntime.ExecuteRequest{
		Code:    "print('hello')",
		Gateway: &mockGateway{},
	}

	_, err := b.Execute(ctx, req)
	if !errors.Is(err, ErrClientNotConfigured) {
		t.Errorf("Execute() without client error = %v, want %v", err, ErrClientNotConfigured)
	}
}

func TestBackendWithMockClient(t *testing.T) {
	mockRunner := &MockContainerRunner{
		RunFunc: func(_ context.Context, spec ContainerSpec) (ContainerResult, error) {
			// Verify spec is built correctly
			if spec.Image != "toolruntime-sandbox:latest" {
				t.Errorf("spec.Image = %q, want %q", spec.Image, "toolruntime-sandbox:latest")
			}
			return ContainerResult{
				ExitCode: 0,
				Stdout:   "hello world",
				Stderr:   "",
			}, nil
		},
	}

	b := New(Config{
		Client: mockRunner,
	})

	ctx := context.Background()
	req := toolruntime.ExecuteRequest{
		Code:    "print('hello')",
		Gateway: &mockGateway{},
	}

	result, err := b.Execute(ctx, req)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.Stdout != "hello world" {
		t.Errorf("Stdout = %q, want %q", result.Stdout, "hello world")
	}
}

func TestBackendWithHealthChecker(t *testing.T) {
	t.Run("healthy daemon", func(t *testing.T) {
		mockRunner := &MockContainerRunner{
			RunFunc: func(_ context.Context, _ ContainerSpec) (ContainerResult, error) {
				return ContainerResult{ExitCode: 0}, nil
			},
		}
		mockHealth := &MockHealthChecker{
			PingFunc: func(_ context.Context) error {
				return nil
			},
		}

		b := New(Config{
			Client:        mockRunner,
			HealthChecker: mockHealth,
		})

		ctx := context.Background()
		req := toolruntime.ExecuteRequest{
			Code:    "print('hello')",
			Gateway: &mockGateway{},
		}

		_, err := b.Execute(ctx, req)
		if err != nil {
			t.Errorf("Execute() error = %v", err)
		}
	})

	t.Run("unhealthy daemon", func(t *testing.T) {
		mockRunner := &MockContainerRunner{}
		mockHealth := &MockHealthChecker{
			PingFunc: func(_ context.Context) error {
				return errors.New("connection refused")
			},
		}

		b := New(Config{
			Client:        mockRunner,
			HealthChecker: mockHealth,
		})

		ctx := context.Background()
		req := toolruntime.ExecuteRequest{
			Code:    "print('hello')",
			Gateway: &mockGateway{},
		}

		_, err := b.Execute(ctx, req)
		if !errors.Is(err, ErrDaemonUnavailable) {
			t.Errorf("Execute() error = %v, want %v", err, ErrDaemonUnavailable)
		}
	})
}

func TestBackendWithImageResolver(t *testing.T) {
	resolvedImage := ""
	mockRunner := &MockContainerRunner{
		RunFunc: func(_ context.Context, spec ContainerSpec) (ContainerResult, error) {
			resolvedImage = spec.Image
			return ContainerResult{ExitCode: 0}, nil
		},
	}
	mockResolver := &MockImageResolver{
		ResolveFunc: func(_ context.Context, image string) (string, error) {
			return image + "@sha256:abc123", nil
		},
	}

	b := New(Config{
		ImageName:     "my-image:v1",
		Client:        mockRunner,
		ImageResolver: mockResolver,
	})

	ctx := context.Background()
	req := toolruntime.ExecuteRequest{
		Code:    "print('hello')",
		Gateway: &mockGateway{},
	}

	_, err := b.Execute(ctx, req)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if resolvedImage != "my-image:v1@sha256:abc123" {
		t.Errorf("resolved image = %q, want %q", resolvedImage, "my-image:v1@sha256:abc123")
	}
}

func TestBackendBuildSpec(t *testing.T) {
	b := New(Config{
		SeccompPath: "/path/to/seccomp.json",
	})

	req := toolruntime.ExecuteRequest{
		Code:    "print('hello')",
		Gateway: &mockGateway{},
		Profile: toolruntime.ProfileHardened,
		Limits: toolruntime.Limits{
			MemoryBytes:    256 * 1024 * 1024,
			CPUQuotaMillis: 1000,
			PidsMax:        100,
		},
	}

	spec, err := b.buildSpec("test-image:latest", req, toolruntime.ProfileHardened)
	if err != nil {
		t.Fatalf("buildSpec() error = %v", err)
	}

	// Verify image
	if spec.Image != "test-image:latest" {
		t.Errorf("Image = %q, want %q", spec.Image, "test-image:latest")
	}

	// Verify security
	if spec.Security.NetworkMode != "none" {
		t.Errorf("Security.NetworkMode = %q, want %q", spec.Security.NetworkMode, "none")
	}
	if !spec.Security.ReadOnlyRootfs {
		t.Error("Security.ReadOnlyRootfs = false, want true")
	}
	if spec.Security.SeccompProfile != "/path/to/seccomp.json" {
		t.Errorf("Security.SeccompProfile = %q, want %q", spec.Security.SeccompProfile, "/path/to/seccomp.json")
	}

	// Verify resources
	if spec.Resources.MemoryBytes != 256*1024*1024 {
		t.Errorf("Resources.MemoryBytes = %d, want %d", spec.Resources.MemoryBytes, 256*1024*1024)
	}
	if spec.Resources.CPUQuota != 1000*1000 { // milliseconds to microseconds
		t.Errorf("Resources.CPUQuota = %d, want %d", spec.Resources.CPUQuota, 1000*1000)
	}
	if spec.Resources.PidsLimit != 100 {
		t.Errorf("Resources.PidsLimit = %d, want %d", spec.Resources.PidsLimit, 100)
	}

	// Verify labels
	if spec.Labels["toolruntime.profile"] != "hardened" {
		t.Errorf("Labels[toolruntime.profile] = %q, want %q", spec.Labels["toolruntime.profile"], "hardened")
	}
	if spec.Labels["toolruntime.backend"] != "docker" {
		t.Errorf("Labels[toolruntime.backend] = %q, want %q", spec.Labels["toolruntime.backend"], "docker")
	}
}

func TestClientError(t *testing.T) {
	t.Run("with container ID", func(t *testing.T) {
		err := &ClientError{
			Op:          "start",
			Image:       "alpine:latest",
			ContainerID: "abc123",
			Err:         errors.New("permission denied"),
		}
		expected := "docker start alpine:latest (abc123): permission denied"
		if err.Error() != expected {
			t.Errorf("Error() = %q, want %q", err.Error(), expected)
		}
	})

	t.Run("without container ID", func(t *testing.T) {
		err := &ClientError{
			Op:    "pull",
			Image: "alpine:latest",
			Err:   errors.New("not found"),
		}
		expected := "docker pull alpine:latest: not found"
		if err.Error() != expected {
			t.Errorf("Error() = %q, want %q", err.Error(), expected)
		}
	})

	t.Run("unwrap", func(t *testing.T) {
		innerErr := errors.New("inner error")
		err := &ClientError{
			Op:    "create",
			Image: "alpine:latest",
			Err:   innerErr,
		}
		if !errors.Is(err, innerErr) {
			t.Error("Unwrap() should return inner error")
		}
	})
}
