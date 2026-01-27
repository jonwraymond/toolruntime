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

func (m *mockGateway) SearchTools(ctx context.Context, _ string, _ int) ([]toolindex.Summary, error) {
	return nil, nil
}

func (m *mockGateway) ListNamespaces(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (m *mockGateway) DescribeTool(ctx context.Context, _ string, _ tooldocs.DetailLevel) (tooldocs.ToolDoc, error) {
	return tooldocs.ToolDoc{}, nil
}

func (m *mockGateway) ListToolExamples(ctx context.Context, _ string, _ int) ([]tooldocs.ToolExample, error) {
	return nil, nil
}

func (m *mockGateway) RunTool(ctx context.Context, _ string, _ map[string]any) (toolrun.RunResult, error) {
	return toolrun.RunResult{}, nil
}

func (m *mockGateway) RunChain(ctx context.Context, _ []toolrun.ChainStep) (toolrun.RunResult, []toolrun.StepResult, error) {
	return toolrun.RunResult{}, nil, nil
}

// TestBackendImplementsInterface verifies Backend satisfies toolruntime.Backend
func TestBackendImplementsInterface(t *testing.T) {
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
