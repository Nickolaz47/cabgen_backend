package mocks

import (
	"context"
	"io"

	"github.com/CABGenOrg/cabgen_backend/internal/pipeline"
)

type MockCmd struct {
	StdoutContent string
	StderrContent string
	RunErr        error
	stdout        io.Writer
	stderr        io.Writer
}

func (m *MockCmd) Start() error    { return nil }
func (m *MockCmd) Wait() error     { return nil }
func (m *MockCmd) SetStdin(_ io.Reader) {}

func (m *MockCmd) SetStdout(w io.Writer) { m.stdout = w }
func (m *MockCmd) SetStderr(w io.Writer) { m.stderr = w }

func (m *MockCmd) Run() error {
	if m.RunErr != nil {
		return m.RunErr
	}
	if m.stdout != nil {
		io.WriteString(m.stdout, m.StdoutContent)
	}
	if m.stderr != nil {
		io.WriteString(m.stderr, m.StderrContent)
	}
	return nil
}

type MockCommander struct {
	CommandFunc func(ctx context.Context, name string, args ...string) pipeline.Cmd
}

func (m *MockCommander) Command(ctx context.Context, name string,
	args ...string) pipeline.Cmd {
	if m.CommandFunc != nil {
		return m.CommandFunc(ctx, name, args...)
	}
	return &MockCmd{}
}
