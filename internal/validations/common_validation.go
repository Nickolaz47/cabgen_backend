package validations

import "github.com/CABGenOrg/cabgen_backend/internal/models"

func ApplySampleUpdate(sample *models.Sample, input *models.SampleUpdateInput) {
	if input.Name != nil {
		sample.Name = *input.Name
	}

	if input.CollectionDate != nil {
		sample.CollectionDate = *input.CollectionDate
	}

	if input.RunNumber != nil {
		sample.RunNumber = *input.RunNumber
	}

	if input.RunDate != nil {
		sample.RunDate = *input.RunDate
	}

	if input.City != nil {
		sample.City = input.City
	}

	if input.OriginCode != nil {
		sample.OriginCode = input.OriginCode
	}

	if input.Gender != nil {
		sample.Gender = input.Gender
	}

	if input.DateOfBirth != nil {
		sample.DateOfBirth = input.DateOfBirth
	}
}

func ApplySampleFilesUpdate(sample *models.Sample,
	input *models.SampleAttachmentInput) {
	if input.Fastq1 != nil {
		sample.Fastq1 = input.Fastq1
	}

	if input.Fastq2 != nil {
		sample.Fastq2 = input.Fastq2
	}

	if input.Fasta != nil {
		sample.Fasta = input.Fasta
	}
}
