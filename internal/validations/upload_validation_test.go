package validations_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/stretchr/testify/assert"
)

func TestIsUploadField(t *testing.T) {
	t.Run("Known Fields", func(t *testing.T) {
		assert.True(t, validations.IsUploadField("fastq1"))
		assert.True(t, validations.IsUploadField("fastq2"))
		assert.True(t, validations.IsUploadField("fasta"))
	})

	t.Run("Unknown Field", func(t *testing.T) {
		assert.False(t, validations.IsUploadField("evil"))
		assert.False(t, validations.IsUploadField(""))
	})
}

func TestIsAllowedUploadFile(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		fileName string
		want     bool
	}{
		{"fastq plain", "fastq1", "reads_R1.fastq", true},
		{"fastq gz", "fastq1", "reads_R1.fastq.gz", true},
		{"fq plain", "fastq2", "reads_R2.fq", true},
		{"fq gz", "fastq2", "reads_R2.fq.gz", true},
		{"fasta plain", "fasta", "contigs.fasta", true},
		{"fna", "fasta", "contigs.fna", true},
		{"fa", "fasta", "contigs.fa", true},
		{"fasta gz rejected", "fasta", "contigs.fasta.gz", false},
		{"exe rejected", "fastq1", "virus.exe", false},
		{"no ext", "fastq1", "reads", false},
		{"long name rejected", "fastq1",
			"a" + string(make([]byte, 254)) + ".fastq", false},
		{"exact 255", "fastq1",
			"b" + string(make([]byte, 248)) + ".fastq", true},
		{"over 255", "fastq1",
			"c" + string(make([]byte, 250)) + ".fastq", false},
		{"case insensitive", "fastq1", "reads_R1.FASTQ", true},
		{"case insensitive gz", "fastq1", "reads_R1.FQ.GZ", true},
		{"unknown field", "evil", "virus.exe", false},
		{"empty field", "", "reads.fastq", false},
		{"empty filename", "fastq1", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validations.IsAllowedUploadFile(tt.field, tt.fileName)
			assert.Equal(t, tt.want, got)
		})
	}
}
