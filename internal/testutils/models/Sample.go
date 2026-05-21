package models

import (
	"time"

	rModels "github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type Sample struct {
	ID             string          `gorm:"primaryKey;default:(hex(randomblob(16)))" json:"id"`
	Name           string          `gorm:"type:varchar(255);not null" json:"name"`
	CollectionDate time.Time       `gorm:"type:date;not null" json:"collection_date"`
	RunNumber      string          `gorm:"type:varchar(255);not null" json:"run_number"`
	RunDate        time.Time       `gorm:"type:date;not null" json:"run_date"`
	City           *string         `gorm:"type:varchar(255);default:null" json:"city,omitempty"`
	OriginCode     *string         `gorm:"type:varchar(255);default:null" json:"origin_code,omitempty"`
	Gender         *rModels.Gender `gorm:"type:varchar(15);default:null" json:"gender,omitempty"`
	DateOfBirth    *time.Time      `gorm:"type:date;default:null" json:"date_of_birth,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Fastq1         *string         `gorm:"type:varchar(255);default:null" json:"fastq1,omitempty"`
	Fastq2         *string         `gorm:"type:varchar(255);default:null" json:"fastq2,omitempty"`
	Fasta          *string         `gorm:"type:varchar(255);default:null" json:"fasta,omitempty"`
	// Foreign Keys
	CountryID       uint                  `gorm:"not null" json:"-"`
	Country         rModels.Country       `gorm:"foreignKey:CountryID;references:ID"`
	UserID          string                `gorm:"not null" json:"-"`
	User            rModels.User          `gorm:"foreignKey:UserID;references:ID"`
	OriginID        string                `gorm:"not null" json:"-"`
	Origin          rModels.Origin        `gorm:"foreignKey:OriginID;references:ID"`
	SampleSourceID  string                `gorm:"not null" json:"-"`
	SampleSource    rModels.SampleSource  `gorm:"foreignKey:SampleSourceID;references:ID"`
	MicroorganismID string                `gorm:"not null" json:"-"`
	Microorganism   rModels.Microorganism `gorm:"foreignKey:MicroorganismID;references:ID"`
	SequencerID     string                `gorm:"not null" json:"-"`
	Sequencer       rModels.Sequencer     `gorm:"foreignKey:SequencerID;references:ID"`
	LaboratoryID    string                `gorm:"not null" json:"-"`
	Laboratory      rModels.Laboratory    `gorm:"foreignKey:LaboratoryID;references:ID"`
	HealthServiceID string                `gorm:"not null" json:"-"`
	HealthService   rModels.HealthService `gorm:"foreignKey:HealthServiceID;references:ID"`
}

func NewSample(
	ID, name string, collectionDate time.Time, runNumber string,
	runDate time.Time, city, originCode string, gender rModels.Gender,
	dateOfBirth time.Time, fastq1, fastq2, fasta string, country rModels.Country,
	user rModels.User, origin rModels.Origin, sampleSource rModels.SampleSource,
	microorganism rModels.Microorganism, sequencer rModels.Sequencer,
	laboratory rModels.Laboratory, healthService rModels.HealthService,
) rModels.Sample {
	return rModels.Sample{
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

func CreateMockSample() rModels.Sample {
	mockUser := NewLoginUser()
	mockCountry := mockUser.Country
	mockOrigin := NewOrigin(uuid.New().String(),
		map[string]string{"pt": "Humano", "en": "Human", "es": "Humano"},
		true,
	)
	mockSampleSource := NewSampleSource(
		uuid.NewString(),
		map[string]string{"pt": "Aspirado",
			"en": "Aspirated", "es": "Aspirado"},
		map[string]string{"pt": "Trato respiratório",
			"en": "Respiratory tract", "es": "Vías respiratorias"},
		true,
	)
	mockMicro := NewMicroorganism(
		uuid.NewString(), rModels.Bacteria, "Neisseria meningitidis",
		map[string]string{
			"pt": "Sorogrupo B", "en": "Serogroup B", "es": "Serogrupo B"}, true,
	)
	mockSequencer := rModels.Sequencer{
		ID:       uuid.New(),
		Brand:    "Illumina",
		Model:    "MySeq",
		IsActive: true,
	}
	mockLab := rModels.Laboratory{
		ID:           uuid.New(),
		Name:         "Laboratorio Central do Rio de Janeiro",
		Abbreviation: "LACEN/RJ",
		IsActive:     true,
	}
	mockHealthService := rModels.HealthService{
		ID:           uuid.New(),
		Name:         "Laboratorio Central do Rio de Janeiro",
		Type:         "Public",
		CountryID:    mockCountry.ID,
		Country:      mockCountry,
		City:         nil,
		Contactant:   nil,
		ContactEmail: nil,
		ContactPhone: nil,
		IsActive:     true,
	}

	id := uuid.New()
	date := time.Date(2024, time.May, 11, 0, 0, 0, 0, time.UTC)
	mockSample := NewSample(
		id.String(), "sample 1", date, "R1", date, "", "A01", rModels.Male,
		date, "read1.fastq", "read2.fastq", "", mockCountry, mockUser,
		mockOrigin, mockSampleSource, mockMicro, mockSequencer, mockLab,
		mockHealthService,
	)

	return mockSample
}

func NewSampleCreateInput(sample rModels.Sample) rModels.SampleCreateInput {
	return rModels.SampleCreateInput{
		Name:            sample.Name,
		CollectionDate:  sample.CollectionDate,
		RunNumber:       sample.RunNumber,
		RunDate:         sample.RunDate,
		City:            sample.City,
		OriginCode:      sample.OriginCode,
		Gender:          sample.Gender,
		DateOfBirth:     sample.DateOfBirth,
		CountryCode:     sample.Country.Code,
		UserID:          sample.UserID,
		OriginID:        sample.OriginID,
		SampleSourceID:  sample.SampleSourceID,
		MicroorganismID: sample.MicroorganismID,
		SequencerID:     sample.SequencerID,
		LaboratoryID:    sample.LaboratoryID,
		HealthServiceID: sample.HealthServiceID,
	}
}
