package pipeline

import (
	"bytes"
	"context"
	"fmt"
	"strings"
)

type ToolRunner interface {
	Run(ctx context.Context, args []string) (string, error)
}

type toolRunner struct {
	commander Commander
}

func NewToolRunner(commander Commander) ToolRunner {
	return &toolRunner{
		commander: commander,
	}
}

func (r *toolRunner) buildBlastXCmd(blastDB, inputFile,
	outputFile string) []string {
	if blastDB == "" || inputFile == "" || outputFile == "" {
		return nil
	}

	return []string{
		"blastx", "-db", blastDB, "-query", inputFile, "-evalue", "0.001",
		"-out", outputFile,
	}
}

func (r *toolRunner) buildFastQCCmd(fastqcCmd, read1, read2,
	outputDir string) []string {
	if fastqcCmd == "" || read1 == "" || read2 == "" || outputDir == "" {
		return nil
	}

	return []string{fastqcCmd, "--quiet", read1, read2, "--outdir", outputDir}
}

func (r *toolRunner) buildUnicyclerCmd(unicyclerCmd, read1, read2, outputDir,
	threads, spadesPath string) []string {
	if unicyclerCmd == "" || read1 == "" || read2 == "" || outputDir == "" ||
		threads == "" {
		return nil
	}

	if spadesPath != "" {
		return []string{
			unicyclerCmd, "-1", read1, "-2", read2, "-o", outputDir,
			"--min_fasta_length", "500", "--mode", "conservative",
			"-t", threads, "--spades_path", spadesPath,
		}
	}

	return []string{
		unicyclerCmd, "-1", read1, "-2", read2, "-o", outputDir,
		"--min_fasta_length", "500", "--mode", "conservative",
		"-t", threads,
	}
}

func (r *toolRunner) buildProkkaCmd(prokkaCmd, outputDir, prefix,
	assemblyPath, threads string) []string {
	if prokkaCmd == "" || outputDir == "" || prefix == "" ||
		assemblyPath == "" || threads == "" {
		return nil
	}

	return []string{
		prokkaCmd, "--outdir", outputDir, "--prefix", prefix,
		assemblyPath, "--force", "--cpus", threads,
	}
}

func (r *toolRunner) buildCheckMLineageCmd(checkmCmd, inputDir,
	outputDir, threads string) []string {
	if checkmCmd == "" || inputDir == "" || outputDir == "" || threads == "" {
		return nil
	}

	return []string{
		checkmCmd, "lineage_wf", "-x", "fasta", inputDir, outputDir,
		"--threads", threads, "--pplacer_threads", "1",
	}
}

func (r *toolRunner) buildCheckMQACmd(checkmCmd, checkmDir,
	sample, threads string) []string {
	if checkmCmd == "" || checkmDir == "" || sample == "" || threads == "" {
		return nil
	}

	return []string{
		checkmCmd, "qa", "-o", "2", "-f",
		checkmDir + "/" + sample + "_resultados",
		"--tab_table", checkmDir + "/lineage.ms",
		checkmDir, "--threads", threads,
	}
}

func (r *toolRunner) buildKraken2Cmd(krakenCmd, dbPath, outputDir,
	threads, assemblyPath string) []string {
	if krakenCmd == "" || dbPath == "" || outputDir == "" ||
		threads == "" || assemblyPath == "" {
		return nil
	}

	return []string{
		krakenCmd, "--db", dbPath, "--use-names",
		"--output", outputDir + "/out_kraken",
		"--threads", threads, assemblyPath,
	}
}

func (r *toolRunner) buildSplitterCmd(threads, inputFile,
	outputFilePrefix string) []string {
	if threads == "" || inputFile == "" || outputFilePrefix == "" {
		return nil
	}

	return []string{
		"split", "--numeric-suffixes=1", "-n", "l/" + threads,
		inputFile, outputFilePrefix,
	}
}

func (r *toolRunner) buildFastANICmd(fastaniCmd, query, refList,
	output, threads string) []string {
	if fastaniCmd == "" || query == "" || refList == "" || output == "" ||
		threads == "" {
		return nil
	}

	return []string{
		fastaniCmd, "-q", query, "--rl", refList,
		"-o", output, "--threads", threads,
	}
}

func (r *toolRunner) buildAbricateCmd(abricateCmd, db, inputFile,
	outputFile, threads string) []string {
	if abricateCmd == "" || db == "" || inputFile == "" || outputFile == "" ||
		threads == "" {
		return nil
	}

	return []string{
		"sh", "-c",
		fmt.Sprintf("%s --db %s %s > %s --threads %s",
			abricateCmd, db, inputFile, outputFile, threads),
	}
}

func (r *toolRunner) buildMLSTCmd(mlstCmd, threads, assemblyPath,
	outputFile string) []string {
	if mlstCmd == "" || threads == "" || assemblyPath == "" ||
		outputFile == "" {
		return nil
	}

	return []string{
		"sh", "-c",
		fmt.Sprintf("%s --threads %s --exclude abaumannii --csv %s > %s",
			mlstCmd, threads, assemblyPath, outputFile),
	}
}

func (r *toolRunner) Run(ctx context.Context, args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", fmt.Errorf("The args cannot be empty.")
	}

	name := args[0]
	cmdArgs := args[1:]

	cmd := r.commander.Command(ctx, name, cmdArgs...)

	var stdout, stderr bytes.Buffer
	cmd.SetStdout(&stdout)
	cmd.SetStderr(&stderr)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf(
			"Command '%s' failed with return code %w. Output: %s. Error: %s",
			strings.Join(args, " "), err, stdout.String(), stderr.String())
	}

	return stdout.String(), nil
}
