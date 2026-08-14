package pipeline

import "errors"

// Input errors
var (
	ErrFileNotFound   = errors.New("input file not found")
	ErrCorruptedInput = errors.New("input file is corrupted or truncated")
	ErrEmptyReads     = errors.New("no reads found in input file")
	ErrInvalidFormat  = errors.New("invalid file format")
)

// Pipeline step errors
var (
	ErrFastQC              = errors.New("The FastQC step failed. Create a new analysis.")
	ErrUnicycler           = errors.New("The Unicycler step failed. Create a new analysis.")
	ErrProkka              = errors.New("The Prokka step failed. Create a new analysis.")
	ErrCheckM              = errors.New("The CheckM step failed. Create a new analysis.")
	ErrKraken2             = errors.New("The Kraken2 step failed. Create a new analysis.")
	ErrSpecies             = errors.New("The Species identification (mlst, fastani) step failed. Create a new analysis.")
	ErrAbricate            = errors.New("The Abricate step failed. Create a new analysis.")
	ErrPrepareFolders      = errors.New("Folder preparation failed. Create a new analysis.")
	ErrAnalysisRun         = errors.New("analysis failed")
	ErrUnknownAnalysisType = errors.New("unknown analysis type")
)
