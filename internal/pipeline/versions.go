package pipeline

import (
	"bytes"
	"context"
	"regexp"
)

type toolSpec struct {
	Name    string
	Cmd     []string
	Pattern string
}

type ToolVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var tools = []toolSpec{
	{"FastQC", []string{"fastqc", "--version"}, `v([\d.]+)`},
	{"Unicycler", []string{"unicycler", "--version"}, `v([\d.]+)`},
	{"Prokka", []string{"prokka", "--version"}, `([\d.]+)`},
	{"CheckM", []string{"checkm", "-h"}, `v([\d.]+)`},
	{"Kraken2", []string{"kraken2", "--version"}, `version ([\d.]+)`},
	{"FastANI", []string{"fastANI", "--version"}, `([\d.]+)`},
	{"Abricate", []string{"abricate", "--version"}, `([\d.]+)`},
	{"MLST", []string{"mlst", "--version"}, `([\d.]+)`},
	{"Blast", []string{"blastx", "-version"}, `blastx: ([\d.]+)`},
}

func GetBioinfoProgramVersions(ctx context.Context,
	commander Commander) []ToolVersion {
	result := make([]ToolVersion, 0, len(tools))

	for _, t := range tools {
		name := t.Cmd[0]
		args := t.Cmd[1:]

		cmd := commander.Command(ctx, name, args...)
		re := regexp.MustCompile(t.Pattern)

		var stdout, stderr bytes.Buffer
		cmd.SetStdout(&stdout)
		cmd.SetStderr(&stderr)

		version := "unknown"
		cmd.Run()

		if m := re.FindStringSubmatch(stderr.String() + stdout.String()); len(m) > 1 {
			version = m[1]
		}

		result = append(result, ToolVersion{Name: t.Name, Version: version})
	}

	return result
}
