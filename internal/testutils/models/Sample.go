package models

import (
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type Sample struct {
	ID             string         `gorm:"primaryKey;default:(hex(randomblob(16)))" json:"id"`
	Name           string         `gorm:"type:varchar(255);not null" json:"name"`
	CollectionDate time.Time      `gorm:"type:date;not null" json:"collection_date"`
	RunNumber      string         `gorm:"type:varchar(255);not null" json:"run_number"`
	RunDate        time.Time      `gorm:"type:date;not null" json:"run_date"`
	City           *string        `gorm:"type:varchar(255);default:null" json:"city,omitempty"`
	OriginCode     *string        `gorm:"type:varchar(255);default:null" json:"origin_code,omitempty"`
	Gender         *models.Gender `gorm:"type:varchar(15);default:null" json:"gender,omitempty"`
	DateOfBirth    *time.Time     `gorm:"type:date;default:null" json:"date_of_birth,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Fastq1         *string        `gorm:"type:varchar(255);default:null" json:"fastq1,omitempty"`
	Fastq2         *string        `gorm:"type:varchar(255);default:null" json:"fastq2,omitempty"`
	Fasta          *string        `gorm:"type:varchar(255);default:null" json:"fasta,omitempty"`
	// Foreign Keys
	CountryID       uint                 `gorm:"not null" json:"-"`
	Country         models.Country       `gorm:"foreignKey:CountryID;references:ID"`
	UserID          string               `gorm:"not null" json:"-"`
	User            models.User          `gorm:"foreignKey:UserID;references:ID"`
	OriginID        string               `gorm:"not null" json:"-"`
	Origin          models.Origin        `gorm:"foreignKey:OriginID;references:ID"`
	SampleSourceID  string               `gorm:"not null" json:"-"`
	SampleSource    models.SampleSource  `gorm:"foreignKey:SampleSourceID;references:ID"`
	MicroorganismID string               `gorm:"not null" json:"-"`
	Microorganism   models.Microorganism `gorm:"foreignKey:MicroorganismID;references:ID"`
	SequencerID     string               `gorm:"not null" json:"-"`
	Sequencer       models.Sequencer     `gorm:"foreignKey:SequencerID;references:ID"`
	LaboratoryID    string               `gorm:"not null" json:"-"`
	Laboratory      models.Laboratory    `gorm:"foreignKey:LaboratoryID;references:ID"`
	HealthServiceID string               `gorm:"not null" json:"-"`
	HealthService   models.HealthService `gorm:"foreignKey:HealthServiceID;references:ID"`
}

func NewSample(
	ID, name string, collectionDate time.Time, runNumber string,
	runDate time.Time, city, originCode string, gender models.Gender,
	dateOfBirth time.Time, fastq1, fastq2, fasta string, country models.Country,
	user models.User, origin models.Origin, sampleSource models.SampleSource,
	microorganism models.Microorganism, sequencer models.Sequencer,
	laboratory models.Laboratory, healthService models.HealthService,
) models.Sample {
	return models.Sample{
		ID:              uuid.MustParse(ID),
		Name:            name,
		CollectionDate:  collectionDate,
		RunNumber:       runNumber,
		RunDate:         runDate,
		City:            &city,
		OriginCode:      &originCode,
		Gender:          &gender,
		DateOfBirth:     &dateOfBirth,
		Fastq1:          &fastq1,
		Fastq2:          &fastq2,
		Fasta:           &fasta,
		CountryID:       country.ID,
		Country:         country,
		UserID:          user.ID,
		User:            user,
		OriginID:        origin.ID,
		Origin:          origin,
		SampleSourceID:  sampleSource.ID,
		SampleSource:    sampleSource,
		MicroorganismID: microorganism.ID,
		Microorganism:   microorganism,
		SequencerID:     sequencer.ID,
		Sequencer:       sequencer,
		LaboratoryID:    laboratory.ID,
		Laboratory:      laboratory,
		HealthServiceID: healthService.ID,
		HealthService:   healthService,
	}
}
