package models

import (
	"encoding/json"
	"time"

	rModels "github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Analysis struct {
	ID string `gorm:"primaryKey;default:(hex(randomblob(16)))"`

	// Pipeline Control
	Type   rModels.AnalysisType   `gorm:"type:varchar(20);not null"`
	Status rModels.AnalysisStatus `gorm:"type:varchar(20);not null;default:'PENDING'"`

	// Paths
	FastQC1 *string `gorm:"type:varchar(255)"`
	FastQC2 *string `gorm:"type:varchar(255)"`

	// Results
	Metrics datatypes.JSON `gorm:"type:jsonb"`

	// Run Metadata
	ErrorMessage *string `gorm:"type:text"`
	StartedAt    *time.Time
	FinishedAt   *time.Time

	// Datetime
	CreatedAt time.Time
	UpdatedAt time.Time

	// Foreign Keys
	SampleID string         `gorm:"type:not null;index"`
	Sample   rModels.Sample `gorm:"foreignKey:SampleID;references:ID"`
	UserID   string         `gorm:"type:not null;index"`
	User     rModels.User   `gorm:"foreignKey:UserID;references:ID"`
}

func NewAnalysis(
	id string, analysisType rModels.AnalysisType,
	analysisStatus rModels.AnalysisStatus, fastqc1, fastqc2 string,
	metrics map[string]any, errorMessage string, startedAt,
	finishedAt time.Time, sample rModels.Sample,
	user rModels.User) rModels.Analysis {

	var metricsBytes datatypes.JSON
	if metrics != nil {
		bytes, err := json.Marshal(metrics)
		if err != nil {
			return rModels.Analysis{}
		}
		metricsBytes = datatypes.JSON(bytes)
	}

	return rModels.Analysis{
		ID:           uuid.MustParse(id),
		Type:         analysisType,
		Status:       analysisStatus,
		FastQC1:      &fastqc1,
		FastQC2:      &fastqc2,
		Metrics:      metricsBytes,
		ErrorMessage: &errorMessage,
		StartedAt:    &startedAt,
		FinishedAt:   &finishedAt,
		SampleID:     sample.ID,
		Sample:       sample,
		UserID:       user.ID,
		User:         user,
	}
}

func CreateMockAnalysis() rModels.Analysis {
	fastqc1 := "/result/fastqc_reads1.html"
	fastqc2 := "/result/fastqc_reads2.html"
	startedAt := time.Date(2024, time.May, 11, 0, 0, 0, 0, time.UTC)
	sample := CreateMockSample()
	user := NewLoginUser()
	metrics := map[string]any{
		"completeness":        95.89,
		"species":             "Acinetobacter sp",
		"mlst":                502,
		"acquired_resistance": "fosA6_1 (resistance to fosfomycin) (allele confidence 98.81)",
	}

	var metricsBytes datatypes.JSON
	bytes, err := json.Marshal(metrics)
	if err != nil {
		return rModels.Analysis{}
	}
	metricsBytes = datatypes.JSON(bytes)

	return rModels.Analysis{
		ID:        uuid.New(),
		Type:      rModels.AnalysisTypeComplete,
		Status:    rModels.AnalysisStatusDone,
		Metrics:   metricsBytes,
		FastQC1:   &fastqc1,
		FastQC2:   &fastqc2,
		StartedAt: &startedAt,
		SampleID:  sample.ID,
		Sample:    sample,
		UserID:    user.ID,
		User:      user,
	}
}

func NewAnalysisCreateDTO(analysis rModels.Analysis) rModels.AnalysisCreateDTO {
	return rModels.AnalysisCreateDTO{
		Type:     analysis.Type,
		SampleID: analysis.Sample.ID,
		UserID:   analysis.User.ID,
	}
}
