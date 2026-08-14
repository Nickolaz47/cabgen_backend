package models

import (
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/pipeline"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const AnalysesByBatch = 50

var errorMessageTranslations = map[error]map[string]string{
	pipeline.ErrFileNotFound: {
		"en": "Input file not found.",
		"pt": "Arquivo de entrada não encontrado.",
		"es": "Archivo de entrada no encontrado.",
	},
	pipeline.ErrCorruptedInput: {
		"en": "The input file is corrupted or truncated.",
		"pt": "O arquivo de entrada está corrompido ou truncado.",
		"es": "El archivo de entrada está corrupto o truncado.",
	},
	pipeline.ErrEmptyReads: {
		"en": "No reads found in the input file.",
		"pt": "Nenhum read encontrado no arquivo de entrada.",
		"es": "No se encontraron reads en el archivo de entrada.",
	},
	pipeline.ErrInvalidFormat: {
		"en": "Invalid file format.",
		"pt": "Formato de arquivo inválido.",
		"es": "Formato de archivo inválido.",
	},
	pipeline.ErrFastQC: {
		"en": "The FastQC step failed. Create a new analysis.",
		"pt": "A etapa do FastQC falhou. Crie uma nova análise.",
		"es": "El paso de FastQC falló. Cree un nuevo análisis.",
	},
	pipeline.ErrUnicycler: {
		"en": "The Unicycler step failed. Create a new analysis.",
		"pt": "A etapa do Unicycler falhou. Crie uma nova análise.",
		"es": "El paso de Unicycler falló. Cree un nuevo análisis.",
	},
	pipeline.ErrProkka: {
		"en": "The Prokka step failed. Create a new analysis.",
		"pt": "A etapa do Prokka falhou. Crie uma nova análise.",
		"es": "El paso de Prokka falló. Cree un nuevo análisis.",
	},
	pipeline.ErrCheckM: {
		"en": "The CheckM step failed. Create a new analysis.",
		"pt": "A etapa do CheckM falhou. Crie uma nova análise.",
		"es": "El paso de CheckM falló. Cree un nuevo análisis.",
	},
	pipeline.ErrKraken2: {
		"en": "The Kraken2 step failed. Create a new analysis.",
		"pt": "A etapa do Kraken2 falhou. Crie uma nova análise.",
		"es": "El paso de Kraken2 falló. Cree un nuevo análisis.",
	},
	pipeline.ErrSpecies: {
		"en": "The Species identification (mlst, fastani) step failed. Create a new analysis.",
		"pt": "A etapa de identificação de espécie (mlst, fastani) falhou. Crie uma nova análise.",
		"es": "El paso de identificación de especie (mlst, fastani) falló. Cree un nuevo análisis.",
	},
	pipeline.ErrAbricate: {
		"en": "The Abricate step failed. Create a new analysis.",
		"pt": "A etapa do Abricate falhou. Crie uma nova análise.",
		"es": "El paso de Abricate falló. Cree un nuevo análisis.",
	},
	pipeline.ErrPrepareFolders: {
		"en": "Folder preparation failed. Create a new analysis.",
		"pt": "Falha ao preparar os arquivos da análise. Crie uma nova análise.",
		"es": "Error al preparar los archivos del análisis. Cree un nuevo análisis.",
	},
}

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

type AnalysisStep string

const (
	StepFastQC    AnalysisStep = "FastQC"
	StepUnicycler AnalysisStep = "Unicycler"
	StepProkka    AnalysisStep = "Prokka"
	StepCheckM    AnalysisStep = "CheckM"
	StepKraken2   AnalysisStep = "Kraken2"
	StepSpecies   AnalysisStep = "Species"
	StepAbricate  AnalysisStep = "Abricate"
	StepCoverage  AnalysisStep = "Coverage"
)

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
	Step   AnalysisStep   `gorm:"type:varchar(20);default:''"`

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

type AnalysisResponse struct {
	ID             uuid.UUID      `json:"id"`
	Type           AnalysisType   `json:"type"`
	Status         AnalysisStatus `json:"status"`
	Step           AnalysisStep   `json:"step"`
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

func (a *Analysis) ToResponse(language string) AnalysisResponse {
	lang := translation.ParseLanguage(language)

	var errorMsg *string
	if a.ErrorMessage != nil {
		msg := *a.ErrorMessage
		for errVal, translations := range errorMessageTranslations {
			if errVal.Error() == msg {
				if translated, ok := translations[lang]; ok {
					errorMsg = &translated
				} else {
					errorMsg = a.ErrorMessage
				}
				break
			}
		}
		if errorMsg == nil {
			errorMsg = a.ErrorMessage
		}
	}

	return AnalysisResponse{
		ID:             a.ID,
		Type:           a.Type,
		Status:         a.Status,
		Step:           a.Step,
		ErrorMessage:   errorMsg,
		Sample:         a.Sample.OriginCode,
		SampleID:       a.SampleID,
		User:           a.User.Username,
		UserID:         a.UserID,
		Metrics:        a.Metrics,
		ResultsZipPath: a.ResultsZipPath,
		FastQC1:        a.FastQC1,
		FastQC2:        a.FastQC2,
		StartedAt:      a.StartedAt,
		FinishedAt:     a.FinishedAt,
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
