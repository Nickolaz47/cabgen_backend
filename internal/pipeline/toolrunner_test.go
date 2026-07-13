package pipeline

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCmd struct {
	runFunc func() error
}

func (m *mockCmd) Start() error          { return nil }
func (m *mockCmd) Run() error            { return m.runFunc() }
func (m *mockCmd) Wait() error           { return nil }
func (m *mockCmd) SetStdout(w io.Writer) {}
func (m *mockCmd) SetStderr(w io.Writer) {}
func (m *mockCmd) SetStdin(r io.Reader)  {}

type mockCommander struct {
	cmdFunc func(ctx context.Context, name string, args ...string) Cmd
}

func (m *mockCommander) Command(ctx context.Context, name string,
	args ...string) Cmd {
	return m.cmdFunc(ctx, name, args...)
}

func TestToolRunnerRun(t *testing.T) {
	ctx := context.Background()

	t.Run("Empty command line returns error", func(t *testing.T) {
		runner := NewToolRunner(&mockCommander{})

		result, err := runner.Run(ctx, "")

		assert.EqualError(t, err, "The command_line argument cannot be empty.")
		assert.Equal(t, "", result)
	})

	t.Run("Successful command returns stdout", func(t *testing.T) {
		runner := NewToolRunner(&mockCommander{
			cmdFunc: func(ctx context.Context, name string, args ...string) Cmd {
				return &mockCmd{
					runFunc: func() error {
						return nil
					},
				}
			},
		})

		result, err := runner.Run(ctx, "echo hello")

		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("Failed command returns error with details", func(t *testing.T) {
		runner := NewToolRunner(&mockCommander{
			cmdFunc: func(ctx context.Context, name string, args ...string) Cmd {
				return &mockCmd{
					runFunc: func() error {
						return fmt.Errorf("exit status 1")
					},
				}
			},
		})

		result, err := runner.Run(ctx, "failing-cmd arg1 arg2")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Command 'failing-cmd arg1 arg2' failed")
		assert.Contains(t, err.Error(), "exit status 1")
		assert.Equal(t, "", result)
	})

	t.Run("Command receives the full command line as name", func(t *testing.T) {
		var capturedName string
		var capturedArgs []string

		runner := NewToolRunner(&mockCommander{
			cmdFunc: func(ctx context.Context, name string, args ...string) Cmd {
				capturedName = name
				capturedArgs = args
				return &mockCmd{
					runFunc: func() error {
						return nil
					},
				}
			},
		})

		_, err := runner.Run(ctx, "tool --verbose --output=file.txt")

		assert.NoError(t, err)
		assert.Equal(t, "tool --verbose --output=file.txt", capturedName)
		assert.Empty(t, capturedArgs)
	})

	t.Run("Commander.Command is called exactly once", func(t *testing.T) {
		callCount := 0
		runner := NewToolRunner(&mockCommander{
			cmdFunc: func(ctx context.Context, name string, args ...string) Cmd {
				callCount++
				return &mockCmd{
					runFunc: func() error {
						return nil
					},
				}
			},
		})

		_, err := runner.Run(ctx, "any-command")

		assert.NoError(t, err)
		assert.Equal(t, 1, callCount)
	})

	t.Run("NewToolRunner returns non-nil ToolRunner", func(t *testing.T) {
		runner := NewToolRunner(&mockCommander{})

		assert.NotNil(t, runner)
	})

	t.Run("Error message includes command line", func(t *testing.T) {
		runner := NewToolRunner(&mockCommander{
			cmdFunc: func(ctx context.Context, name string, args ...string) Cmd {
				return &mockCmd{
					runFunc: func() error {
						return fmt.Errorf("permission denied")
					},
				}
			},
		})

		_, err := runner.Run(ctx, "restricted-cmd")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "restricted-cmd")
		assert.Contains(t, err.Error(), "permission denied")
	})

	t.Run("Single command with no arguments", func(t *testing.T) {
		var capturedName string
		runner := NewToolRunner(&mockCommander{
			cmdFunc: func(ctx context.Context, name string, args ...string) Cmd {
				capturedName = name
				return &mockCmd{
					runFunc: func() error {
						return nil
					},
				}
			},
		})

		_, err := runner.Run(ctx, "standalone")

		assert.NoError(t, err)
		assert.Equal(t, "standalone", capturedName)
	})

	t.Run("Empty command line does not call Commander", func(t *testing.T) {
		callCount := 0
		runner := NewToolRunner(&mockCommander{
			cmdFunc: func(ctx context.Context, name string, args ...string) Cmd {
				callCount++
				return &mockCmd{
					runFunc: func() error {
						return nil
					},
				}
			},
		})

		_, err := runner.Run(ctx, "")

		assert.Error(t, err)
		assert.Equal(t, 0, callCount)
	})

	t.Run("Error wraps underlying command error", func(t *testing.T) {
		runner := NewToolRunner(&mockCommander{
			cmdFunc: func(ctx context.Context, name string, args ...string) Cmd {
				return &mockCmd{
					runFunc: func() error {
						return fmt.Errorf("signal: killed")
					},
				}
			},
		})

		_, err := runner.Run(ctx, "long-running-cmd")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signal: killed")
	})
}
