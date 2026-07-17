package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const AnalysesByBatch = 50

type AnalysisStatus string

const (
	AnalysisStatusPending AnalysisStatus = "PENDING"
	AnalysisStatusRunning AnalysisStatus = "RUNNING"
	AnalysisStatusDone    AnalysisStatus = "DONE"
	AnalysisStatusFailed  AnalysisStatus = "FAILED"
)

func (a AnalysisStatus) IsValid() bool {
	switch a {
	case AnalysisStatusPending, AnalysisStatusRunning, AnalysisStatusDone,
		AnalysisStatusFailed:
		return true
	default:
		return false
	}
}

type AnalysisType string

const (
	AnalysisTypeFastQC   AnalysisType = "FASTQC"
	AnalysisTypeGenome   AnalysisType = "GENOME"
	AnalysisTypeComplete AnalysisType = "COMPLETE"
)

func (a AnalysisType) IsValid() bool {
	switch a {
	case AnalysisTypeFastQC, AnalysisTypeGenome, AnalysisTypeComplete:
		return true
	default:
		return false
	}
}

var AnalysisTypes = []AnalysisType{AnalysisTypeFastQC, AnalysisTypeGenome,
	AnalysisTypeComplete}

type AnalysisResults struct {
	// --- Genomic Coverage ---
	Coverage float64 `json:"coverage,omitempty"`

	// --- Assembly Quality (CheckM) ---
	CheckMCompleteness  string `json:"completeness,omitempty"`
	CheckMContamination string `json:"contamination,omitempty"`
	CheckMGenomeSize    string `json:"genome_size,omitempty"`
	CheckMN50           string `json:"n50,omitempty"`

	// --- Taxonomy and Typing ---
	PrimarySpeciesName   string `json:"primary_species,omitempty"`
	SecondarySpeciesName string `json:"secondary_species,omitempty"`
	MLST                 string `json:"mlst,omitempty"`

	// --- Identified Mutations ---
	PoliMutations  []string `json:"poli_mutations,omitempty"`
	OtherMutations []string `json:"other_mutations,omitempty"`

	// --- Virulence (Abricate) ---
	ResfinderGenes []string `json:"gene,omitempty"`
	ResfinderBlast []string `json:"resfinder,omitempty"`
	VFDB           []string `json:"vfdb,omitempty"`
	PlasmidFinder  []string `json:"plasmid,omitempty"`
}

type Analysis struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	// Pipeline Control
	Type   AnalysisType   `gorm:"type:varchar(20);not null"`
	Status AnalysisStatus `gorm:"type:varchar(20);not null;default:'PENDING'"`

	// Paths
	FastQC1 *string `gorm:"type:varchar(255)"`
	FastQC2 *string `gorm:"type:varchar(255)"`

	// Results
	Metrics        datatypes.JSON `gorm:"type:jsonb"`
	ResultsZipPath *string        `gorm:"type:varchar(255)"`

	// Run Metadata
	ErrorMessage *string `gorm:"type:text"`
	StartedAt    *time.Time
	FinishedAt   *time.Time

	// Datetime
	CreatedAt time.Time
	UpdatedAt time.Time

	// Foreign Keys
	SampleID uuid.UUID `gorm:"type:uuid;not null;index"`
	Sample   Sample    `gorm:"foreignKey:SampleID;references:ID"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index"`
	User     User      `gorm:"foreignKey:UserID;references:ID"`
}

type AnalysisAdminResponse struct {
	ID             uuid.UUID      `json:"id"`
	Type           AnalysisType   `json:"type"`
	Status         AnalysisStatus `json:"status"`
	ErrorMessage   *string        `json:"error_message"`
	Sample         string         `json:"sample"`
	SampleID       uuid.UUID      `json:"sample_id"`
	User           string         `json:"user"`
	UserID         uuid.UUID      `json:"user_id"`
	Metrics        datatypes.JSON `json:"metrics"`
	ResultsZipPath *string        `json:"results_zip_path"`
	FastQC1        *string        `json:"fastqc1"`
	FastQC2        *string        `json:"fastqc2"`
	StartedAt      *time.Time     `json:"started_at"`
	FinishedAt     *time.Time     `json:"finished_at"`
}

func (a *Analysis) ToAdminResponse() AnalysisAdminResponse {
	return AnalysisAdminResponse{
		ID:           a.ID,
		Type:         a.Type,
		Status:       a.Status,
		ErrorMessage: a.ErrorMessage,
		Sample:       a.Sample.Name,
		SampleID:     a.Sample.ID,
		User:         a.User.Username,
		UserID:       a.UserID,
		Metrics:      a.Metrics,
		FastQC1:      a.FastQC1,
		FastQC2:      a.FastQC2,
		StartedAt:    a.StartedAt,
		FinishedAt:   a.FinishedAt,
	}
}

type AnalysisResponse struct {
	ID             uuid.UUID      `json:"id"`
	Type           AnalysisType   `json:"type"`
	Status         AnalysisStatus `json:"status"`
	ErrorMessage   *string        `json:"error_message"`
	Sample         string         `json:"sample"`
	SampleID       uuid.UUID      `json:"sample_id"`
	Metrics        datatypes.JSON `json:"metrics"`
	ResultsZipPath *string        `json:"results_zip_path"`
	FastQC1        *string        `json:"fastqc1"`
	FastQC2        *string        `json:"fastqc2"`
	StartedAt      *time.Time     `json:"started_at"`
	FinishedAt     *time.Time     `json:"finished_at"`
}

func (a *Analysis) ToResponse() AnalysisResponse {
	return AnalysisResponse{
		ID:           a.ID,
		Type:         a.Type,
		Status:       a.Status,
		ErrorMessage: a.ErrorMessage,
		Sample:       a.Sample.Name,
		SampleID:     a.SampleID,
		Metrics:      a.Metrics,
		FastQC1:      a.FastQC1,
		FastQC2:      a.FastQC2,
		StartedAt:    a.StartedAt,
		FinishedAt:   a.FinishedAt,
	}
}

type AdminAnalysisCreateInput struct {
	Type     AnalysisType `json:"type" binding:"required"`
	SampleID uuid.UUID    `json:"sample_id" binding:"required"`
	UserID   uuid.UUID    `json:"user_id" binding:"required"`
}

type AnalysisCreateInput struct {
	Type     AnalysisType `json:"type" binding:"required"`
	SampleID uuid.UUID    `json:"sample_id" binding:"required"`
}

type AnalysisCreateDTO struct {
	Type     AnalysisType
	SampleID uuid.UUID
	UserID   uuid.UUID
}

func AnalysisCreateInputToDTO(i AnalysisCreateInput,
	userID uuid.UUID) AnalysisCreateDTO {
	return AnalysisCreateDTO{
		Type:     i.Type,
		SampleID: i.SampleID,
		UserID:   userID,
	}
}

type AdminAnalysisUpdateInput struct {
	Status         *AnalysisStatus `json:"status" binding:"omitempty"`
	Metrics        *datatypes.JSON `json:"metrics" binding:"omitempty"`
	FastQC1        *string         `json:"fastqc1" binding:"omitempty"`
	FastQC2        *string         `json:"fastqc2" binding:"omitempty"`
	ResultsZipPath *string         `json:"results_zip_path" binding:"omitempty"`
	ErrorMessage   *string         `json:"error_message" binding:"omitempty"`
}

type AnalysisTSVDownloadInput struct {
	IDs []uuid.UUID `json:"ids" binding:"required,min=1"`
}
