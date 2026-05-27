package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AnalysisStatus string

const (
	AnalysisStatusPending AnalysisStatus = "PENDING"
	AnalysisStatusRunning AnalysisStatus = "RUNNING"
	AnalysisStatusDone    AnalysisStatus = "DONE"
	AnalysisStatusFailed  AnalysisStatus = "FAILED"
)

type AnalysisType string

const (
	AnalysisTypeFastQC   AnalysisType = "FASTQC"
	AnalysisTypeGenome   AnalysisType = "GENOME"
	AnalysisTypeComplete AnalysisType = "COMPLETE"
)

type Analysis struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	// Pipeline Control
	Type   AnalysisType   `gorm:"type:varchar(20);not null"`
	Status AnalysisStatus `gorm:"type:varchar(20);not null;default:'PENDING'"`

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
	SampleID uuid.UUID `gorm:"type:uuid;not null;index"`
	Sample   Sample    `gorm:"foreignKey:SampleID;references:ID"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index"`
	User     User      `gorm:"foreignKey:UserID;references:ID"`
}

type AnalysisAdminResponse struct {
	ID           uuid.UUID      `json:"id"`
	Type         AnalysisType   `json:"type"`
	Status       AnalysisStatus `json:"status"`
	ErrorMessage *string        `json:"error_message"`
	Sample       string         `json:"sample"`
	SampleID     uuid.UUID      `json:"sample_id"`
	User         string         `json:"user"`
	Metrics      datatypes.JSON `json:"metrics"`
	FastQC1      *string        `json:"fastqc1"`
	FastQC2      *string        `json:"fastqc2"`
	StartedAt    *time.Time     `json:"started_at"`
	FinishedAt   *time.Time     `json:"finished_at"`
}

type AnalysisResponse struct {
	ID           uuid.UUID      `json:"id"`
	Type         AnalysisType   `json:"type"`
	Status       AnalysisStatus `json:"status"`
	ErrorMessage *string        `json:"error_message"`
	Sample       string         `json:"sample"`
	SampleID     uuid.UUID      `json:"sample_id"`
	Metrics      datatypes.JSON `json:"metrics"`
	FastQC1      *string        `json:"fastqc1"`
	FastQC2      *string        `json:"fastqc2"`
	StartedAt    *time.Time     `json:"started_at"`
	FinishedAt   *time.Time     `json:"finished_at"`
}

type AdminAnalysisCreateInput struct {
	Type     AnalysisType `json:"type" binding:"required,oneof=FASTQC GENOME COMPLETE"`
	SampleID uuid.UUID    `json:"sample_id" binding:"required"`
	UserID   uuid.UUID    `json:"user_id" binding:"required"`
}

type AnalysisCreateInput struct {
	Type     AnalysisType `json:"type" binding:"required,oneof=FASTQC GENOME COMPLETE"`
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
	Status       *AnalysisStatus `json:"status" binding:"omitempty,oneof=PENDING RUNNING DONE FAILED"`
	ErrorMessage *string         `json:"error_message" binding:"omitempty"`
	FastQC1      *string         `json:"fastqc1" binding:"omitempty"`
	FastQC2      *string         `json:"fastqc2" binding:"omitempty"`
}
