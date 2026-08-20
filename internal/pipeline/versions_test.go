package pipeline_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/pipeline"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetBioinfoProgramVersions(t *testing.T) {
	toolOutputs := map[string]string{
		"fastqc":    "FastQC v1.2.3",
		"unicycler": "unicycler v1.2.3",
		"prokka":    "prokka 1.2.3",
		"checkm":    "checkm v1.2.3",
		"kraken2":   "kraken2 version 1.2.3",
		"fastANI":   "fastANI 1.2.3",
		"abricate":  "abricate 1.2.3",
		"mlst":      "mlst 1.2.3",
		"blastx":    "blastx: 1.2.3",
	}

	t.Run("Success", func(t *testing.T) {
		cmd := &mocks.MockCommander{
			CommandFunc: func(_ context.Context, name string,
				_ ...string) pipeline.Cmd {
				return &mocks.MockCmd{StdoutContent: toolOutputs[name]}
			},
		}

		result := pipeline.GetBioinfoProgramVersions(
			context.Background(), cmd,
		)

		assert.Len(t, result, 9)

		expected := map[string]string{
			"fastqc": "1.2.3", "unicycler": "1.2.3",
			"prokka": "1.2.3", "checkm": "1.2.3",
			"kraken2": "1.2.3", "fastANI": "1.2.3",
			"abricate": "1.2.3", "mlst": "1.2.3",
			"blastx": "1.2.3",
		}

		for _, tv := range result {
			assert.Equal(t, expected[tv.Name], tv.Version)
		}
	})

	t.Run("Success - Version From Stderr", func(t *testing.T) {
		stderrOutputs := map[string]string{
			"checkm": "checkm v2.1.0",
		}

		cmd := &mocks.MockCommander{
			CommandFunc: func(_ context.Context, name string,
				_ ...string) pipeline.Cmd {
				if stderr, ok := stderrOutputs[name]; ok {
					return &mocks.MockCmd{StderrContent: stderr}
				}
				return &mocks.MockCmd{StdoutContent: toolOutputs[name]}
			},
		}

		result := pipeline.GetBioinfoProgramVersions(
			context.Background(), cmd,
		)

		for _, tv := range result {
			if tv.Name == "checkm" {
				assert.Equal(t, "2.1.0", tv.Version)
			} else {
				assert.Equal(t, "1.2.3", tv.Version)
			}
		}
	})

	t.Run("Partial Failure", func(t *testing.T) {
		failing := map[string]bool{
			"fastqc": true, "unicycler": true, "prokka": true,
		}

		cmd := &mocks.MockCommander{
			CommandFunc: func(_ context.Context, name string,
				_ ...string) pipeline.Cmd {
				if failing[name] {
					return &mocks.MockCmd{
						RunErr: fmt.Errorf("not found"),
					}
				}
				return &mocks.MockCmd{StdoutContent: toolOutputs[name]}
			},
		}

		result := pipeline.GetBioinfoProgramVersions(
			context.Background(), cmd,
		)

		assert.Len(t, result, 9)

		for _, tv := range result {
			if failing[tv.Name] {
				assert.Equal(t, "unknown", tv.Version)
			} else {
				assert.Equal(t, "1.2.3", tv.Version)
			}
		}
	})

	t.Run("Regex Mismatch", func(t *testing.T) {
		cmd := &mocks.MockCommander{
			CommandFunc: func(_ context.Context, _ string,
				_ ...string) pipeline.Cmd {
				return &mocks.MockCmd{StdoutContent: "no version here"}
			},
		}

		result := pipeline.GetBioinfoProgramVersions(
			context.Background(), cmd,
		)

		assert.Len(t, result, 9)

		for _, tv := range result {
			assert.Equal(t, "unknown", tv.Version)
		}
	})
}
