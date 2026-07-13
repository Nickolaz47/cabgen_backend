package pipeline

import (
	"bytes"
	"context"
	"fmt"
)

type ToolRunner interface {
	Run(ctx context.Context, commandLine string) (string, error)
}

type toolRunner struct {
	commander Commander
}

func NewToolRunner(commander Commander) ToolRunner {
	return &toolRunner{
		commander: commander,
	}
}

func (r *toolRunner) Run(ctx context.Context,
	commandLine string) (string, error) {
	if commandLine == "" {
		return "", fmt.Errorf("The command_line argument cannot be empty.")
	}

	cmd := r.commander.Command(ctx, commandLine)

	var stdout, stderr bytes.Buffer
	cmd.SetStdout(&stdout)
	cmd.SetStderr(&stderr)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf(
			"Command '%s' failed with return code %w. Output: %s. Error: %s",
			commandLine, err, stdout.String(), stderr.String())
	}

	return stdout.String(), nil
}
