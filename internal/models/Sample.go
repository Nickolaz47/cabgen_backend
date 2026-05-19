package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Gender string

const (
	Male        Gender = "Male"
	Female      Gender = "Female"
	Unspecified Gender = "Unspecified"
)

func (g Gender) IsValid() bool {
	switch g {
	case Male, Female, Unspecified:
		return true
	default:
		return false
	}
}

type Sample struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name           string     `gorm:"type:varchar(255);not null" json:"name"`
	CollectionDate time.Time  `gorm:"type:date;not null" json:"collection_date"`
	RunNumber      string     `gorm:"type:varchar(255);not null" json:"run_number"`
	RunDate        time.Time  `gorm:"type:date;not null" json:"run_date"`
	City           *string    `gorm:"type:varchar(255);default:null" json:"city,omitempty"`
	OriginCode     *string    `gorm:"type:varchar(255);default:null" json:"origin_code,omitempty"`
	Gender         *Gender    `gorm:"type:varchar(15);default:null" json:"gender,omitempty"`
	DateOfBirth    *time.Time `gorm:"type:date;default:null" json:"date_of_birth,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Fastq1         string     `gorm:"type:varchar(255);default:not null" json:"fastq1"`
	Fastq2         string     `gorm:"type:varchar(255);default:not null" json:"fastq2"`
	Fasta          *string    `gorm:"type:varchar(255);default:null" json:"fasta,omitempty"`
	// Foreign Keys
	CountryID       uint          `gorm:"not null" json:"-"`
	Country         Country       `gorm:"foreignKey:CountryID;references:ID"`
	UserID          uuid.UUID     `gorm:"not null" json:"-"`
	User            User          `gorm:"foreignKey:UserID;references:ID"`
	OriginID        uuid.UUID     `gorm:"not null" json:"-"`
	Origin          Origin        `gorm:"foreignKey:OriginID;references:ID"`
	SampleSourceID  uuid.UUID     `gorm:"not null" json:"-"`
	SampleSource    SampleSource  `gorm:"foreignKey:SampleSourceID;references:ID"`
	MicroorganismID uuid.UUID     `gorm:"not null" json:"-"`
	Microorganism   Microorganism `gorm:"foreignKey:MicroorganismID;references:ID"`
	SequencerID     uuid.UUID     `gorm:"not null" json:"-"`
	Sequencer       Sequencer     `gorm:"foreignKey:SequencerID;references:ID"`
	LaboratoryID    uuid.UUID     `gorm:"not null" json:"-"`
	Laboratory      Laboratory    `gorm:"foreignKey:LaboratoryID;references:ID"`
	HealthServiceID uuid.UUID     `gorm:"not null" json:"-"`
	HealthService   HealthService `gorm:"foreignKey:HealthServiceID;references:ID"`
}

type SampleResponse struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	CollectionDate time.Time  `json:"collection_date"`
	RunNumber      string     `json:"run_number"`
	RunDate        time.Time  `json:"run_date"`
	City           *string    `json:"city"`
	OriginCode     *string    `json:"origin_code"`
	Gender         *Gender    `json:"gender"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	Fastq1         string     `json:"fastq1"`
	Fastq2         string     `json:"fastq2"`
	Fasta          *string    `json:"fasta"`
	// Foreign Keys
	CountryCode   string `json:"country_code"`
	User          string `json:"user"`
	Origin        string `json:"origin"`
	SampleSource  string `json:"sample_source"`
	Microorganism string `json:"microorganism"`
	Sequencer     string `json:"sequencer"`
	Laboratory    string `json:"laboratory"`
	HealthService string `json:"health_service"`
}

func (s *Sample) ToResponse(language string) SampleResponse {
	if language == "" {
		language = "en"
	}

	sequencer := fmt.Sprintf("%s - %s", s.Sequencer.Brand, s.Sequencer.Model)
	species := fmt.Sprintf("%s %s", s.Microorganism.Species,
		s.Microorganism.Variety[language])

	return SampleResponse{
		ID:             s.ID,
		Name:           s.Name,
		CollectionDate: s.CollectionDate,
		RunNumber:      s.RunNumber,
		RunDate:        s.RunDate,
		City:           s.City,
		OriginCode:     s.OriginCode,
		Gender:         s.Gender,
		DateOfBirth:    s.DateOfBirth,
		Fastq1:         s.Fastq1,
		Fastq2:         s.Fastq2,
		Fasta:          s.Fasta,
		CountryCode:    s.Country.Code,
		User:           s.User.Username,
		Origin:         s.Origin.Names[language],
		SampleSource:   s.SampleSource.Names[language],
		Microorganism:  species,
		Sequencer:      sequencer,
		Laboratory:     s.Laboratory.Name,
		HealthService:  s.HealthService.Name,
	}
}

type SampleCreateInput struct {
	Name           string     `json:"name" binding:"required,min=3,max=100"`
	CollectionDate time.Time  `json:"collection_date" binding:"required" time_format:"2006-01-02"`
	RunNumber      string     `json:"run_number" binding:"required"`
	RunDate        time.Time  `json:"run_date" binding:"required" time_format:"2006-01-02"`
	City           *string    `json:"city,omitempty" binding:"omitempty,max=255"`
	OriginCode     *string    `json:"origin_code,omitempty" binding:"omitempty,max=255"`
	Gender         *Gender    `json:"gender,omitempty" binding:"omitempty,max=255"`
	DateOfBirth    *time.Time `json:"date_of_birth,omitempty" binding:"omitempty,max=255"`
	Fastq1         string     `json:"fastq1" binding:"required,min=4,max=255"`
	Fastq2         string     `json:"fastq2" binding:"required,min=4,max=255"`
	// Foreign Keys
	CountryCode     string    `json:"country_code" binding:"required,len=3"`
	UserID          uuid.UUID `json:"user_id" binding:"required"`
	OriginID        uuid.UUID `json:"origin_id" binding:"required"`
	SampleSourceID  uuid.UUID `json:"sample_source_id" binding:"required"`
	MicroorganismID uuid.UUID `json:"microorganism_id" binding:"required"`
	SequencerID     uuid.UUID `json:"sequencer_id" binding:"required"`
	LaboratoryID    uuid.UUID `json:"laboratory_id" binding:"required"`
	HealthServiceID uuid.UUID `json:"health_service_id" binding:"required"`
}

type SampleUpdateInput struct {
	Name           string     `json:"name" binding:"omitempty,min=3,max=100"`
	CollectionDate time.Time  `json:"collection_date" binding:"omitempty" time_format:"2006-01-02"`
	RunNumber      *string    `json:"run_number,omitempty" binding:"omitempty"`
	RunDate        *time.Time `json:"run_date,omitempty" binding:"omitempty" time_format:"2006-01-02"`
	City           *string    `json:"city,omitempty" binding:"omitempty,max=255"`
	OriginCode     *string    `json:"origin_code,omitempty" binding:"omitempty,max=255"`
	Gender         *Gender    `json:"gender,omitempty" binding:"omitempty,max=255"`
	DateOfBirth    *time.Time `json:"date_of_birth,omitempty" binding:"omitempty,max=255"`
	Fastq1         *string    `json:"fastq1,omitempty" binding:"omitempty,min=4,max=255"`
	Fastq2         *string    `json:"fastq2,omitempty" binding:"omitempty,min=4,max=255"`
	// Foreign Keys
	CountryCode     *string    `json:"country_code,omitempty" binding:"omitempty,len=3"`
	UserID          *uuid.UUID `json:"user_id,omitempty" binding:"omitempty"`
	OriginID        *uuid.UUID `json:"origin_id,omitempty" binding:"omitempty"`
	SampleSourceID  *uuid.UUID `json:"sample_source_id,omitempty" binding:"omitempty"`
	MicroorganismID *uuid.UUID `json:"microorganism_id,omitempty" binding:"omitempty"`
	SequencerID     *uuid.UUID `json:"sequencer_id,omitempty" binding:"omitempty"`
	LaboratoryID    *uuid.UUID `json:"laboratory_id,omitempty" binding:"omitempty"`
	HealthServiceID *uuid.UUID `json:"health_service_id,omitempty" binding:"omitempty"`
}
