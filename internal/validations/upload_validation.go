package validations

import (
	"slices"
	"strings"
)

var uploadFieldExts = map[string][]string{
	"fastq1": {".fastq", ".fastq.gz", ".fq", ".fq.gz"},
	"fastq2": {".fastq", ".fastq.gz", ".fq", ".fq.gz"},
	"fasta":  {".fasta", ".fna", ".fa"},
}

func IsUploadField(name string) bool {
	_, ok := uploadFieldExts[name]
	return ok
}

func IsAllowedUploadFile(field, fileName string) bool {
	if len(fileName) > 255 {
		return false
	}
	name := strings.ToLower(fileName)
	return slices.ContainsFunc(uploadFieldExts[field],
		func(ext string) bool { return strings.HasSuffix(name, ext) })
}
