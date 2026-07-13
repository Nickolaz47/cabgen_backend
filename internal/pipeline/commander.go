package pipeline

import (
	"context"
	"io"
	"os/exec"
)

type Cmd interface {
	Start() error
	Run() error
	Wait() error
	SetStdout(io.Writer)
	SetStderr(io.Writer)
	SetStdin(io.Reader)
}

type Commander interface {
	Command(ctx context.Context, name string, args ...string) Cmd
}

type RealCmd struct {
	cmd *exec.Cmd
}

func (r *RealCmd) Start() error           { return r.cmd.Start() }
func (r *RealCmd) Run() error             { return r.cmd.Run() }
func (r *RealCmd) Wait() error            { return r.cmd.Wait() }
func (r *RealCmd) SetStdout(w io.Writer)  { r.cmd.Stdout = w }
func (r *RealCmd) SetStderr(w io.Writer)  { r.cmd.Stderr = w }
func (r *RealCmd) SetStdin(rdr io.Reader) { r.cmd.Stdin = rdr }

type RealCommander struct{}

func (r *RealCommander) Command(ctx context.Context, name string,
	args ...string) Cmd {
	return &RealCmd{cmd: exec.CommandContext(ctx, name, args...)}
}
