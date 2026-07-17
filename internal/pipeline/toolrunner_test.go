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

	t.Run("Error - Empty Args", func(t *testing.T) {
		runner := NewToolRunner(&mockCommander{})

		result, err := runner.Run(ctx, []string{})

		assert.EqualError(t, err, "The args cannot be empty.")
		assert.Equal(t, "", result)
	})

	t.Run("Error - Nil Args", func(t *testing.T) {
		runner := NewToolRunner(&mockCommander{})

		result, err := runner.Run(ctx, nil)

		assert.EqualError(t, err, "The args cannot be empty.")
		assert.Equal(t, "", result)
	})

	t.Run("Success", func(t *testing.T) {
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

	t.Run("Error - Command Failed", func(t *testing.T) {
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

	t.Run("Captures Name and Args", func(t *testing.T) {
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

	t.Run("Called Exactly Once", func(t *testing.T) {
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

	t.Run("Returns Non-nil", func(t *testing.T) {
		runner := NewToolRunner(&mockCommander{})

		assert.NotNil(t, runner)
	})

	t.Run("Error - Full Command in Message", func(t *testing.T) {
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

	t.Run("Success - No Arguments", func(t *testing.T) {
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

	t.Run("Error - Empty String Args", func(t *testing.T) {
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

	t.Run("Error - Wraps Underlying", func(t *testing.T) {
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

	t.Run("Error - Multiple Args", func(t *testing.T) {
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

func TestBuildBlastXCmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Success", func(t *testing.T) {
		result := runner.BuildBlastXCmd("nr", "contigs.fa", "blastx_out.txt")

		assert.Equal(t, []string{
			"blastx", "-db", "nr", "-query", "contigs.fa",
			"-evalue", "0.001", "-out", "blastx_out.txt",
		}, result)
	})

	t.Run("Empty blastDB", func(t *testing.T) {
		assert.Nil(t, runner.BuildBlastXCmd("", "contigs.fa", "blastx_out.txt"))
	})

	t.Run("Empty inputFile", func(t *testing.T) {
		assert.Nil(t, runner.BuildBlastXCmd("nr", "", "blastx_out.txt"))
	})

	t.Run("Empty outputFile", func(t *testing.T) {
		assert.Nil(t, runner.BuildBlastXCmd("nr", "contigs.fa", ""))
	})
}

func TestBuildFastQCCmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Success", func(t *testing.T) {
		result := runner.BuildFastQCCmd("fastqc", "read1.fq", "read2.fq", "/out")

		assert.Equal(t, []string{"fastqc", "--quiet", "read1.fq", "read2.fq", "--outdir", "/out"}, result)
	})

	t.Run("Empty fastqcCmd", func(t *testing.T) {
		assert.Nil(t, runner.BuildFastQCCmd("", "read1.fq", "read2.fq", "/out"))
	})

	t.Run("Empty read1", func(t *testing.T) {
		assert.Nil(t, runner.BuildFastQCCmd("fastqc", "", "read2.fq", "/out"))
	})

	t.Run("Empty read2", func(t *testing.T) {
		assert.Nil(t, runner.BuildFastQCCmd("fastqc", "read1.fq", "", "/out"))
	})

	t.Run("Empty outputDir", func(t *testing.T) {
		assert.Nil(t, runner.BuildFastQCCmd("fastqc", "read1.fq", "read2.fq", ""))
	})
}

func TestBuildUnicyclerCmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Success", func(t *testing.T) {
		result := runner.BuildUnicyclerCmd("unicycler", "r1.fq", "r2.fq", "/out", "4", "/spades")

		assert.Equal(t, []string{
			"unicycler", "-1", "r1.fq", "-2", "r2.fq", "-o", "/out",
			"--min_fasta_length", "500", "--mode", "conservative",
			"-t", "4", "--spades_path", "/spades",
		}, result)
	})

	t.Run("Empty unicyclerCmd", func(t *testing.T) {
		assert.Nil(t, runner.BuildUnicyclerCmd("", "r1.fq", "r2.fq", "/out", "4", "/spades"))
	})

	t.Run("Empty read1", func(t *testing.T) {
		assert.Nil(t, runner.BuildUnicyclerCmd("unicycler", "", "r2.fq", "/out", "4", "/spades"))
	})

	t.Run("Empty threads", func(t *testing.T) {
		assert.Nil(t, runner.BuildUnicyclerCmd("unicycler", "r1.fq", "r2.fq", "/out", "", "/spades"))
	})
}

func TestBuildProkkaCmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Success", func(t *testing.T) {
		result := runner.BuildProkkaCmd("prokka", "/out", "sample", "contigs.fa", "8")

		assert.Equal(t, []string{
			"prokka", "--outdir", "/out", "--prefix", "sample",
			"contigs.fa", "--force", "--cpus", "8",
		}, result)
	})

	t.Run("Empty prokkaCmd", func(t *testing.T) {
		assert.Nil(t, runner.BuildProkkaCmd("", "/out", "sample", "contigs.fa", "8"))
	})

	t.Run("Empty outputDir", func(t *testing.T) {
		assert.Nil(t, runner.BuildProkkaCmd("prokka", "", "sample", "contigs.fa", "8"))
	})

	t.Run("Empty prefix", func(t *testing.T) {
		assert.Nil(t, runner.BuildProkkaCmd("prokka", "/out", "", "contigs.fa", "8"))
	})
}

func TestBuildCheckMLineageCmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Success", func(t *testing.T) {
		result := runner.BuildCheckMLineageCmd("checkm", "/in", "/out", "4")

		assert.Equal(t, []string{
			"checkm", "lineage_wf", "-x", "fasta", "/in", "/out",
			"--threads", "4", "--pplacer_threads", "1",
		}, result)
	})

	t.Run("Empty checkmCmd", func(t *testing.T) {
		assert.Nil(t, runner.BuildCheckMLineageCmd("", "/in", "/out", "4"))
	})

	t.Run("Empty inputDir", func(t *testing.T) {
		assert.Nil(t, runner.BuildCheckMLineageCmd("checkm", "", "/out", "4"))
	})
}

func TestBuildCheckMQACmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Success", func(t *testing.T) {
		result := runner.BuildCheckMQACmd("checkm", "/dir", "sample1", "4")

		assert.Equal(t, []string{
			"checkm", "qa", "-o", "2", "-f",
			"/dir/sample1_results",
			"--tab_table", "/dir/lineage.ms",
			"/dir", "--threads", "4",
		}, result)
	})

	t.Run("Empty checkmCmd", func(t *testing.T) {
		assert.Nil(t, runner.BuildCheckMQACmd("", "/dir", "sample1", "4"))
	})

	t.Run("Empty sample", func(t *testing.T) {
		assert.Nil(t, runner.BuildCheckMQACmd("checkm", "/dir", "", "4"))
	})
}

func TestBuildKraken2Cmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Success", func(t *testing.T) {
		result := runner.BuildKraken2Cmd("kraken2", "/db", "/out", "4", "contigs.fa")

		assert.Equal(t, []string{
			"kraken2", "--db", "/db", "--use-names",
			"--output", "/out/out_kraken",
			"--threads", "4", "contigs.fa",
		}, result)
	})

	t.Run("Empty krakenCmd", func(t *testing.T) {
		assert.Nil(t, runner.BuildKraken2Cmd("", "/db", "/out", "4", "contigs.fa"))
	})

	t.Run("Empty dbPath", func(t *testing.T) {
		assert.Nil(t, runner.BuildKraken2Cmd("kraken2", "", "/out", "4", "contigs.fa"))
	})
}

func TestBuildSplitterCmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Success", func(t *testing.T) {
		result := runner.BuildSplitterCmd("4", "input.fq", "prefix_")

		assert.Equal(t, []string{
			"split", "--numeric-suffixes=1", "-n", "l/4",
			"input.fq", "prefix_",
		}, result)
	})

	t.Run("Empty threads", func(t *testing.T) {
		assert.Nil(t, runner.BuildSplitterCmd("", "input.fq", "prefix_"))
	})

	t.Run("Empty inputFile", func(t *testing.T) {
		assert.Nil(t, runner.BuildSplitterCmd("4", "", "prefix_"))
	})
}

func TestBuildFastANICmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Success", func(t *testing.T) {
		result := runner.BuildFastANICmd("fastani", "query.fna", "ref.txt", "/out", "4")

		assert.Equal(t, []string{
			"fastani", "-q", "query.fna", "--rl", "ref.txt",
			"-o", "/out", "--threads", "4",
		}, result)
	})

	t.Run("Empty fastaniCmd", func(t *testing.T) {
		assert.Nil(t, runner.BuildFastANICmd("", "query.fna", "ref.txt", "/out", "4"))
	})

	t.Run("Empty query", func(t *testing.T) {
		assert.Nil(t, runner.BuildFastANICmd("fastani", "", "ref.txt", "/out", "4"))
	})
}

func TestBuildAbricateCmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Success", func(t *testing.T) {
		result := runner.BuildAbricateCmd("abricate", "resfinder", "input.fna", "out.txt", "4")

		assert.Equal(t, []string{
			"sh", "-c",
			"abricate --db resfinder input.fna > out.txt --threads 4",
		}, result)
	})

	t.Run("Empty abricateCmd", func(t *testing.T) {
		assert.Nil(t, runner.BuildAbricateCmd("", "resfinder", "input.fna", "out.txt", "4"))
	})

	t.Run("Empty db", func(t *testing.T) {
		assert.Nil(t, runner.BuildAbricateCmd("abricate", "", "input.fna", "out.txt", "4"))
	})
}

func TestBuildMLSTCmd(t *testing.T) {
	runner := &toolRunner{}

	t.Run("Success", func(t *testing.T) {
		result := runner.BuildMLSTCmd("mlst", "4", "contigs.fa", "out.txt")

		assert.Equal(t, []string{
			"sh", "-c",
			"mlst --threads 4 --exclude abaumannii --csv contigs.fa > out.txt",
		}, result)
	})

	t.Run("Empty mlstCmd", func(t *testing.T) {
		assert.Nil(t, runner.BuildMLSTCmd("", "4", "contigs.fa", "out.txt"))
	})

	t.Run("Empty threads", func(t *testing.T) {
		assert.Nil(t, runner.BuildMLSTCmd("mlst", "", "contigs.fa", "out.txt"))
	})
}
