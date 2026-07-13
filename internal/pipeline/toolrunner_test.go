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

func (m *mockCmd) Start() error           { return nil }
func (m *mockCmd) Run() error             { return m.runFunc() }
func (m *mockCmd) Wait() error            { return nil }
func (m *mockCmd) SetStdout(w io.Writer)  {}
func (m *mockCmd) SetStderr(w io.Writer)  {}
func (m *mockCmd) SetStdin(r io.Reader)   {}

type mockCommander struct {
	cmdFunc func(ctx context.Context, name string, args ...string) Cmd
}

func (m *mockCommander) Command(ctx context.Context, name string, args ...string) Cmd {
	return m.cmdFunc(ctx, name, args...)
}

func TestToolRunnerRun(t *testing.T) {
	ctx := context.Background()

	t.Run("Empty args returns error", func(t *testing.T) {
		runner := NewToolRunner(&mockCommander{})

		result, err := runner.Run(ctx, []string{})

		assert.EqualError(t, err, "The args cannot be empty.")
		assert.Equal(t, "", result)
	})

	t.Run("Nil args returns error", func(t *testing.T) {
		runner := NewToolRunner(&mockCommander{})

		result, err := runner.Run(ctx, nil)

		assert.EqualError(t, err, "The args cannot be empty.")
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

		result, err := runner.Run(ctx, []string{"echo", "hello"})

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

		result, err := runner.Run(ctx, []string{"failing-cmd", "arg1", "arg2"})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Command 'failing-cmd arg1 arg2' failed")
		assert.Contains(t, err.Error(), "exit status 1")
		assert.Equal(t, "", result)
	})

	t.Run("Command receives name and args separately", func(t *testing.T) {
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

		_, err := runner.Run(ctx, []string{"tool", "--verbose", "--output=file.txt"})

		assert.NoError(t, err)
		assert.Equal(t, "tool", capturedName)
		assert.Equal(t, []string{"--verbose", "--output=file.txt"}, capturedArgs)
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

		_, err := runner.Run(ctx, []string{"any-command"})

		assert.NoError(t, err)
		assert.Equal(t, 1, callCount)
	})

	t.Run("NewToolRunner returns non-nil ToolRunner", func(t *testing.T) {
		runner := NewToolRunner(&mockCommander{})

		assert.NotNil(t, runner)
	})

	t.Run("Error message includes full command string", func(t *testing.T) {
		runner := NewToolRunner(&mockCommander{
			cmdFunc: func(ctx context.Context, name string, args ...string) Cmd {
				return &mockCmd{
					runFunc: func() error {
						return fmt.Errorf("permission denied")
					},
				}
			},
		})

		_, err := runner.Run(ctx, []string{"restricted-cmd", "--flag"})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "restricted-cmd --flag")
		assert.Contains(t, err.Error(), "permission denied")
	})

	t.Run("Single command with no arguments", func(t *testing.T) {
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

		_, err := runner.Run(ctx, []string{"standalone"})

		assert.NoError(t, err)
		assert.Equal(t, "standalone", capturedName)
		assert.Empty(t, capturedArgs)
	})

	t.Run("Empty string args returns error", func(t *testing.T) {
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

		_, err := runner.Run(ctx, []string{""})

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

		_, err := runner.Run(ctx, []string{"long-running-cmd"})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signal: killed")
	})

	t.Run("Multiple args are joined in error message", func(t *testing.T) {
		runner := NewToolRunner(&mockCommander{
			cmdFunc: func(ctx context.Context, name string, args ...string) Cmd {
				return &mockCmd{
					runFunc: func() error {
						return fmt.Errorf("failed")
					},
				}
			},
		})

		_, err := runner.Run(ctx, []string{"cmd", "a", "b", "c"})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cmd a b c")
	})
}

func TestBuildFastQCCmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Returns correct slice", func(t *testing.T) {
		result := runner.buildFastQCCmd("fastqc", "read1.fq", "read2.fq", "/out")

		assert.Equal(t, []string{"fastqc", "--quiet", "read1.fq", "read2.fq", "--outdir", "/out"}, result)
	})

	t.Run("Returns nil for empty fastqcCmd", func(t *testing.T) {
		assert.Nil(t, runner.buildFastQCCmd("", "read1.fq", "read2.fq", "/out"))
	})

	t.Run("Returns nil for empty read1", func(t *testing.T) {
		assert.Nil(t, runner.buildFastQCCmd("fastqc", "", "read2.fq", "/out"))
	})

	t.Run("Returns nil for empty read2", func(t *testing.T) {
		assert.Nil(t, runner.buildFastQCCmd("fastqc", "read1.fq", "", "/out"))
	})

	t.Run("Returns nil for empty outputDir", func(t *testing.T) {
		assert.Nil(t, runner.buildFastQCCmd("fastqc", "read1.fq", "read2.fq", ""))
	})
}

func TestBuildUnicyclerCmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Returns correct slice", func(t *testing.T) {
		result := runner.buildUnicyclerCmd("unicycler", "r1.fq", "r2.fq", "/out", "4", "/spades")

		assert.Equal(t, []string{
			"unicycler", "-1", "r1.fq", "-2", "r2.fq", "-o", "/out",
			"--min_fasta_length", "500", "--mode", "conservative",
			"-t", "4", "--spades_path", "/spades",
		}, result)
	})

	t.Run("Returns nil for empty unicyclerCmd", func(t *testing.T) {
		assert.Nil(t, runner.buildUnicyclerCmd("", "r1.fq", "r2.fq", "/out", "4", "/spades"))
	})

	t.Run("Returns nil for empty read1", func(t *testing.T) {
		assert.Nil(t, runner.buildUnicyclerCmd("unicycler", "", "r2.fq", "/out", "4", "/spades"))
	})

	t.Run("Returns nil for empty threads", func(t *testing.T) {
		assert.Nil(t, runner.buildUnicyclerCmd("unicycler", "r1.fq", "r2.fq", "/out", "", "/spades"))
	})
}

func TestBuildProkkaCmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Returns correct slice", func(t *testing.T) {
		result := runner.buildProkkaCmd("prokka", "/out", "sample", "contigs.fa", "8")

		assert.Equal(t, []string{
			"prokka", "--outdir", "/out", "--prefix", "sample",
			"contigs.fa", "--force", "--cpus", "8",
		}, result)
	})

	t.Run("Returns nil for empty prokkaCmd", func(t *testing.T) {
		assert.Nil(t, runner.buildProkkaCmd("", "/out", "sample", "contigs.fa", "8"))
	})

	t.Run("Returns nil for empty outputDir", func(t *testing.T) {
		assert.Nil(t, runner.buildProkkaCmd("prokka", "", "sample", "contigs.fa", "8"))
	})

	t.Run("Returns nil for empty prefix", func(t *testing.T) {
		assert.Nil(t, runner.buildProkkaCmd("prokka", "/out", "", "contigs.fa", "8"))
	})
}

func TestBuildCheckMLineageCmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Returns correct slice", func(t *testing.T) {
		result := runner.buildCheckMLineageCmd("checkm", "/in", "/out", "4")

		assert.Equal(t, []string{
			"checkm", "lineage_wf", "-x", "fasta", "/in", "/out",
			"--threads", "4", "--pplacer_threads", "1",
		}, result)
	})

	t.Run("Returns nil for empty checkmCmd", func(t *testing.T) {
		assert.Nil(t, runner.buildCheckMLineageCmd("", "/in", "/out", "4"))
	})

	t.Run("Returns nil for empty inputDir", func(t *testing.T) {
		assert.Nil(t, runner.buildCheckMLineageCmd("checkm", "", "/out", "4"))
	})
}

func TestBuildCheckMQACmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Returns correct slice", func(t *testing.T) {
		result := runner.buildCheckMQACmd("checkm", "/dir", "sample1", "4")

		assert.Equal(t, []string{
			"checkm", "qa", "-o", "2", "-f",
			"/dir/sample1_resultados",
			"--tab_table", "/dir/lineage.ms",
			"/dir", "--threads", "4",
		}, result)
	})

	t.Run("Returns nil for empty checkmCmd", func(t *testing.T) {
		assert.Nil(t, runner.buildCheckMQACmd("", "/dir", "sample1", "4"))
	})

	t.Run("Returns nil for empty sample", func(t *testing.T) {
		assert.Nil(t, runner.buildCheckMQACmd("checkm", "/dir", "", "4"))
	})
}

func TestBuildKraken2Cmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Returns correct slice", func(t *testing.T) {
		result := runner.buildKraken2Cmd("kraken2", "/db", "/out", "4", "contigs.fa")

		assert.Equal(t, []string{
			"kraken2", "--db", "/db", "--use-names",
			"--output", "/out/out_kraken",
			"--threads", "4", "contigs.fa",
		}, result)
	})

	t.Run("Returns nil for empty krakenCmd", func(t *testing.T) {
		assert.Nil(t, runner.buildKraken2Cmd("", "/db", "/out", "4", "contigs.fa"))
	})

	t.Run("Returns nil for empty dbPath", func(t *testing.T) {
		assert.Nil(t, runner.buildKraken2Cmd("kraken2", "", "/out", "4", "contigs.fa"))
	})
}

func TestBuildSplitterCmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Returns correct slice", func(t *testing.T) {
		result := runner.buildSplitterCmd("4", "input.fq", "prefix_")

		assert.Equal(t, []string{
			"split", "--numeric-suffixes=1", "-n", "l/4",
			"input.fq", "prefix_",
		}, result)
	})

	t.Run("Returns nil for empty threads", func(t *testing.T) {
		assert.Nil(t, runner.buildSplitterCmd("", "input.fq", "prefix_"))
	})

	t.Run("Returns nil for empty inputFile", func(t *testing.T) {
		assert.Nil(t, runner.buildSplitterCmd("4", "", "prefix_"))
	})
}

func TestBuildFastANICmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Returns correct slice", func(t *testing.T) {
		result := runner.buildFastANICmd("fastani", "query.fna", "ref.txt", "/out", "4")

		assert.Equal(t, []string{
			"fastani", "-q", "query.fna", "--rl", "ref.txt",
			"-o", "/out", "--threads", "4",
		}, result)
	})

	t.Run("Returns nil for empty fastaniCmd", func(t *testing.T) {
		assert.Nil(t, runner.buildFastANICmd("", "query.fna", "ref.txt", "/out", "4"))
	})

	t.Run("Returns nil for empty query", func(t *testing.T) {
		assert.Nil(t, runner.buildFastANICmd("fastani", "", "ref.txt", "/out", "4"))
	})
}

func TestBuildAbricateCmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Returns correct slice with sh -c wrapper", func(t *testing.T) {
		result := runner.buildAbricateCmd("abricate", "resfinder", "input.fna", "out.txt", "4")

		assert.Equal(t, []string{
			"sh", "-c",
			"abricate --db resfinder input.fna > out.txt --threads 4",
		}, result)
	})

	t.Run("Returns nil for empty abricateCmd", func(t *testing.T) {
		assert.Nil(t, runner.buildAbricateCmd("", "resfinder", "input.fna", "out.txt", "4"))
	})

	t.Run("Returns nil for empty db", func(t *testing.T) {
		assert.Nil(t, runner.buildAbricateCmd("abricate", "", "input.fna", "out.txt", "4"))
	})
}

func TestBuildMLSTCmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Returns correct slice with sh -c wrapper", func(t *testing.T) {
		result := runner.buildMLSTCmd("mlst", "4", "contigs.fa", "out.txt")

		assert.Equal(t, []string{
			"sh", "-c",
			"mlst --threads 4 --exclude abaumannii --csv contigs.fa > out.txt",
		}, result)
	})

	t.Run("Returns nil for empty mlstCmd", func(t *testing.T) {
		assert.Nil(t, runner.buildMLSTCmd("", "4", "contigs.fa", "out.txt"))
	})

	t.Run("Returns nil for empty threads", func(t *testing.T) {
		assert.Nil(t, runner.buildMLSTCmd("mlst", "", "contigs.fa", "out.txt"))
	})
}
