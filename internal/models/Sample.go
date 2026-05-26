package models

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/google/uuid"
)

type Gender string

const (
	Male        Gender = "Male"
	Female      Gender = "Female"
	Unspecified Gender = "Unspecified"
)

var (
	genderTranslations = map[Gender]map[string]string{
		Male: {
			"en": "Male", "es": "Masculino", "pt": "Masculino",
		},
		Female: {
			"en": "Female", "es": "Femenino", "pt": "Feminino",
		},
		Unspecified: {
			"en": "Unspecified", "es": "No especificado", "pt": "Não especificado",
		},
	}
)

func ToGender(value string) Gender {
	value = strings.ToLower(value)

	switch value {
	case "masculino", "male":
		return Male
	case "feminino", "female", "femenino":
		return Female
	default:
		return Unspecified
	}
}

func (g Gender) IsValid() bool {
	switch g {
	case Male, Female, Unspecified:
		return true
	default:
		return false
	}
}

func (g *Gender) ToTranslatedString(language string) *string {
	if g == nil {
		return nil
	}

	language = translation.ParseLanguage(language)

	translations, ok := genderTranslations[*g]
	if !ok {
		val := genderTranslations[Unspecified][language]
		return &val
	}

	val := translations[language]
	return &val
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
	Fastq1         *string    `gorm:"type:varchar(255);default:null" json:"fastq1,omitempty"`
	Fastq2         *string    `gorm:"type:varchar(255);default:null" json:"fastq2,omitempty"`
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
	Gender         *string    `json:"gender"`
	DateOfBirth    *time.Time `json:"date_of_birth"`
	Fastq1         *string    `json:"fastq1"`
	Fastq2         *string    `json:"fastq2"`
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
	language = translation.ParseLanguage(language)

	sequencer := fmt.Sprintf("%s - %s", s.Sequencer.Brand, s.Sequencer.Model)
	species := fmt.Sprintf("%s %s", s.Microorganism.Species,
		s.Microorganism.Variety[language])
	gender := s.Gender.ToTranslatedString(language)

	var fastq1Path, fastq2Path, fastaPath *string

	if s.Fastq1 != nil {
		path := filepath.Base(*s.Fastq1)
		fastq1Path = &path
	}
	if s.Fastq2 != nil {
		path := filepath.Base(*s.Fastq2)
		fastq2Path = &path
	}
	if s.Fasta != nil {
		path := filepath.Base(*s.Fasta)
		fastaPath = &path
	}

	return SampleResponse{
		ID:             s.ID,
		Name:           s.Name,
		CollectionDate: s.CollectionDate,
		RunNumber:      s.RunNumber,
		RunDate:        s.RunDate,
		City:           s.City,
		OriginCode:     s.OriginCode,
		Gender:         gender,
		DateOfBirth:    s.DateOfBirth,
		Fastq1:         fastq1Path,
		Fastq2:         fastq2Path,
		Fasta:          fastaPath,
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
	RunNumber      string     `json:"run_number" binding:"required,max=50"`
	RunDate        time.Time  `json:"run_date" binding:"required" time_format:"2006-01-02"`
	City           *string    `json:"city,omitempty" binding:"omitempty,max=255"`
	OriginCode     *string    `json:"origin_code,omitempty" binding:"omitempty,max=255"`
	Gender         *Gender    `json:"gender,omitempty" binding:"omitempty"`
	DateOfBirth    *time.Time `json:"date_of_birth,omitempty" binding:"omitempty" time_format:"2006-01-02"`
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
	Name           *string    `json:"name" binding:"omitempty,min=3,max=100"`
	CollectionDate *time.Time `json:"collection_date" binding:"omitempty" time_format:"2006-01-02"`
	RunNumber      *string    `json:"run_number,omitempty" binding:"omitempty,max=50"`
	RunDate        *time.Time `json:"run_date,omitempty" binding:"omitempty" time_format:"2006-01-02"`
	City           *string    `json:"city,omitempty" binding:"omitempty,max=255"`
	OriginCode     *string    `json:"origin_code,omitempty" binding:"omitempty,max=255"`
	Gender         *Gender    `json:"gender,omitempty" binding:"omitempty"`
	DateOfBirth    *time.Time `json:"date_of_birth,omitempty" binding:"omitempty" time_format:"2006-01-02"`
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

type SampleAttachmentInput struct {
	Fastq1 *string `json:"fastq1" binding:"max=255"`
	Fastq2 *string `json:"fastq2" binding:"max=255"`
	Fasta  *string `json:"fasta" binding:"max=255"`
}
