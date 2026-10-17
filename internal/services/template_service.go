package services

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

var templateFreeColumns = map[string]string{
	"A": "origin_code",
	"B": "collection_date",
	"C": "run_number",
	"D": "run_date",
	"G": "date_of_birth",
}

var templateRequiredColumns = map[string]struct{}{
	"A": {}, "B": {}, "C": {}, "D": {}, "E": {}, "H": {}, "I": {},
	"J": {}, "K": {}, "L": {}, "M": {}, "N": {}, "O": {},
}

var templateColumnDescriptions = map[string]map[string]string{
	"A": {
		"en": "Sample identification code at the originating laboratory.",
		"es": "Código de identificación de la muestra en el laboratorio de origen.",
		"pt": "Código de identificação da amostra no laboratório de origem.",
	},
	"B": {
		"en": "Date the sample was collected (dd/mm/yyyy).",
		"es": "Fecha de recolección de la muestra (dd/mm/aaaa).",
		"pt": "Data da coleta da amostra (dd/mm/aaaa).",
	},
	"C": {
		"en": "Number or identifier of the sequencing run.",
		"es": "Número o identificador de la corrida de secuenciación.",
		"pt": "Número ou identificador da corrida de sequenciamento.",
	},
	"D": {
		"en": "Date of the sequencing run (dd/mm/yyyy).",
		"es": "Fecha de la corrida de secuenciación (dd/mm/aaaa).",
		"pt": "Data da corrida de sequenciamento (dd/mm/aaaa).",
	},
	"E": {
		"en": "City where the sample was collected.",
		"es": "Ciudad donde se recolectó la muestra.",
		"pt": "Cidade da coleta da amostra.",
	},
	"F": {
		"en": "Patient sex.",
		"es": "Sexo del paciente.",
		"pt": "Sexo do paciente.",
	},
	"G": {
		"en": "Patient date of birth (dd/mm/yyyy).",
		"es": "Fecha de nacimiento del paciente (dd/mm/aaaa).",
		"pt": "Data de nascimento do paciente (dd/mm/aaaa).",
	},
	"H": {
		"en": "Whether to include this sample in the genomic network. Samples included in the network are shown in the monthly surveillance Dashboard.",
		"es": "Si desea incluir esta muestra en la red genómica. Las muestras incluidas en la red se presentan en el Dashboard mensual de vigilancia.",
		"pt": "Se deseja incluir essa amostra na rede genômica. As amostras incluídas na rede são apresentadas no Dashboard mensal de vigilância.",
	},
	"I": {
		"en": "Country where the sample was collected.",
		"es": "País donde se recolectó la muestra.",
		"pt": "País de coleta da amostra.",
	},
	"J": {
		"en": "Sample origin.",
		"es": "Origen de la muestra.",
		"pt": "Origem da amostra.",
	},
	"K": {
		"en": "Sample source (biological material).",
		"es": "Fuente de la muestra (material biológico).",
		"pt": "Fonte da amostra (material biológico).",
	},
	"L": {
		"en": "Microorganism identified in the sample.",
		"es": "Microorganismo identificado en la muestra.",
		"pt": "Microrganismo identificado na amostra.",
	},
	"M": {
		"en": "Equipment used for sequencing.",
		"es": "Equipo utilizado en la secuenciación.",
		"pt": "Equipamento utilizado no sequenciamento.",
	},
	"N": {
		"en": "Laboratory responsible for sequencing.",
		"es": "Laboratorio responsable de la secuenciación.",
		"pt": "Laboratório responsável pelo sequenciamento.",
	},
	"O": {
		"en": "Health service where the sample came from.",
		"es": "Servicio de salud de origen de la muestra.",
		"pt": "Serviço de saúde de origem da amostra.",
	},
}

func toMap(options []models.SelectOption) map[string]string {
	m := make(map[string]string)
	for _, opt := range options {
		m[opt.Label] = opt.Value
	}

	return m
}

// Empty label keeps the zero value; value outside the list = error with row/column
func validateDropdown(raw string, row int, header string,
	col utils.Column) (string, error) {
	if raw == "" {
		return "", nil
	}

	v, ok := col.FindRealValue(raw)
	if !ok {
		return "", &TableValueError{Row: row, Col: header}
	}

	return v, nil
}

// True when every cell is empty
func emptyRow(row []string) bool {
	for _, c := range row {
		if strings.TrimSpace(c) != "" {
			return false
		}
	}

	return true
}

// Serial (raw cell number) wins over text; text tries ISO then DD/MM (BR)
func parseTableDate(raw string) (models.Date, bool) {
	raw = strings.TrimSpace(raw)
	if f, err := strconv.ParseFloat(raw, 64); err == nil && f > 0 &&
		f <= maxExcelDateSerial {
		if t, tErr := excelize.ExcelDateToTime(f, false); tErr == nil {
			y, m, d := t.Date()
			return models.Date{
				Time: time.Date(y, m, d, 0, 0, 0, 0, time.UTC),
			}, true
		}
	}
	for _, layout := range []string{"2006-01-02", "02/01/2006"} {
		if t, err := time.Parse(layout, raw); err == nil {
			return models.Date{Time: t}, true
		}
	}

	return models.Date{}, false
}

// 9999-12-31 in the Excel 1900 date system
const maxExcelDateSerial = 2958465

type TemplateService interface {
	BuildColumns(ctx context.Context, language string) ([]utils.Column,
		map[string]struct{}, error)
	CreateTemplateTable(ctx context.Context, language string) (
		*excelize.File, error)
	ValidateTemplateTable(ctx context.Context, language string,
		file *excelize.File) ([]models.SampleCreateInput, error)
}

type templateService struct {
	CountryService       CountryService
	CityService          CityService
	SelectOptionsService SelectOptionsService
	Logger               *zap.Logger
}

func NewTemplateService(countryService CountryService, cityService CityService,
	selectOptionsService SelectOptionsService,
	logger *zap.Logger) TemplateService {
	return &templateService{
		CountryService:       countryService,
		CityService:          cityService,
		SelectOptionsService: selectOptionsService,
		Logger:               logger,
	}
}

func (s *templateService) BuildColumns(ctx context.Context,
	language string) ([]utils.Column, map[string]struct{}, error) {
	cities, err := s.CityService.FindAll(ctx, language)
	if err != nil {
		s.Logger.Error(
			"Service Error", logging.ServiceLogging(ctx,
				"TemplateService", "CreateTemplateTable",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, nil, ErrInternal
	}

	countries, err := s.CountryService.FindAll(ctx, language)
	if err != nil {
		s.Logger.Error(
			"Service Error", logging.ServiceLogging(ctx,
				"TemplateService", "CreateTemplateTable",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, nil, ErrInternal
	}

	selectOptions, err := s.SelectOptionsService.FindAllFormSelects(ctx,
		language)
	if err != nil {
		s.Logger.Error(
			"Service Error", logging.ServiceLogging(ctx,
				"TemplateService", "CreateTemplateTable",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, nil, ErrInternal
	}

	columns := make([]utils.Column, 0)
	for letter, header := range templateFreeColumns {
		col := utils.NewColumn(letter, header, nil)
		columns = append(columns, col)
	}

	cityColumn := utils.NewColumn("E", "city", toMap(cities))
	genderColumn := utils.NewColumn("F", "gender", toMap(selectOptions.Genders))
	inNetworkOpts := []models.SelectOption{
		{Label: models.ToTranslatedInNetwork(true, language), Value: "true"},
		{Label: models.ToTranslatedInNetwork(false, language), Value: "false"},
	}
	inNetworkColumn := utils.NewColumn("H", "in_network",
		toMap(inNetworkOpts))
	countryColumn := utils.NewColumn("I", "country_code", toMap(countries))
	originColumn := utils.NewColumn("J", "origin_id",
		toMap(selectOptions.Origins))
	sampleSourceColumn := utils.NewColumn("K", "sample_source_id",
		toMap(selectOptions.SampleSources))
	microorganismColumn := utils.NewColumn("L", "microorganism_id",
		toMap(selectOptions.Microorganisms))
	sequencerColumn := utils.NewColumn("M", "sequencer_id",
		toMap(selectOptions.Sequencers))
	laboratoryColumn := utils.NewColumn("N", "laboratory_id",
		toMap(selectOptions.Laboratories))
	healthServiceColumn := utils.NewColumn("O", "health_service_id",
		toMap(selectOptions.HealthServices))

	columns = append(columns, cityColumn, genderColumn, inNetworkColumn,
		countryColumn, originColumn, sampleSourceColumn, microorganismColumn,
		sequencerColumn, laboratoryColumn, healthServiceColumn)

	lang := translation.ParseLanguage(language)
	for i := range columns {
		col := &columns[i]
		if desc, ok := templateColumnDescriptions[col.Letter]; ok {
			col.Description = desc[lang]
		}
		if _, ok := templateRequiredColumns[col.Letter]; ok {
			col.Required = true
		}
	}

	columnHeaders := make(map[string]struct{})
	for _, col := range columns {
		columnHeaders[col.Header] = struct{}{}
	}

	return columns, columnHeaders, nil
}

func (s *templateService) CreateTemplateTable(ctx context.Context,
	language string) (*excelize.File, error) {
	columns, _, err := s.BuildColumns(ctx, language)
	if err != nil {
		return nil, err
	}

	f, err := utils.GenerateMetadataTemplate(columns)
	if err != nil {
		s.Logger.Error(
			"Service Error", logging.ServiceLogging(ctx,
				"TemplateService", "CreateTemplateTable",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	return f, nil
}

func (s *templateService) ValidateTemplateTable(ctx context.Context,
	language string, file *excelize.File) ([]models.SampleCreateInput, error) {
	columns, columnHeaders, err := s.BuildColumns(ctx, language)
	if err != nil {
		s.Logger.Error(
			"Service Error", logging.ServiceLogging(ctx,
				"TemplateService", "ValidateTemplateTable",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, err
	}

	// Accept dropdown labels in every supported language: download and upload
	// may carry different Accept-Language headers
	for _, lang := range translation.Languages {
		if lang == language {
			continue
		}
		other, _, oErr := s.BuildColumns(ctx, lang)
		if oErr != nil {
			return nil, oErr
		}
		for i := range columns {
			for _, o := range other {
				if o.Header == columns[i].Header {
					columns[i].MergeLabels(o)
					break
				}
			}
		}
	}

	// Read the sample sheet; missing or unreadable = invalid file
	rows, err := file.GetRows(utils.SheetName,
		excelize.Options{RawCellValue: true})
	if err != nil {
		s.Logger.Error(
			"Service Error", logging.ServiceLogging(ctx,
				"TemplateService", "ValidateTemplateTable",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInvalidTable
	}

	// Sheet with no rows = empty table
	if len(rows) == 0 {
		s.Logger.Error(
			"Service Error", logging.ServiceLogging(ctx,
				"TemplateService", "ValidateTemplateTable",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrEmptyTable
	}

	// More data rows than allowed = rejected
	if len(rows)-1 > utils.MaxTableRows {
		s.Logger.Error(
			"Service Error", logging.ServiceLogging(ctx,
				"TemplateService", "ValidateTemplateTable",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrTooManyRows
	}

	// Compare submitted headers against the contract: missing or unknown = error.
	// Required headers may carry a trailing mandatory marker (asterisk)
	normalizeHeader := func(h string) string {
		h = strings.TrimSuffix(strings.TrimSpace(h), "*")
		return strings.ToLower(strings.TrimSpace(h))
	}

	present := make(map[string]struct{}, len(rows[0]))
	for _, h := range rows[0] {
		present[normalizeHeader(h)] = struct{}{}
	}

	var invalid []string
	for header := range columnHeaders {
		if _, ok := present[normalizeHeader(header)]; !ok {
			invalid = append(invalid, header)
		}
	}
	for h := range present {
		if _, ok := columnHeaders[h]; !ok {
			invalid = append(invalid, h)
		}
	}
	if len(invalid) > 0 {
		sort.Strings(invalid)
		return s.logValidationError(ctx,
			&TableHeadersError{Invalid: invalid})
	}

	// Index header name -> column position (column order in the sheet is free)
	headerIdx := map[string]int{}
	for i, h := range rows[0] {
		headerIdx[normalizeHeader(h)] = i
	}
	columnsByHeader := make(map[string]utils.Column, len(columns))
	for _, col := range columns {
		key := normalizeHeader(col.Header)
		columnsByHeader[key] = col
	}

	// Iterate data rows (row 1 = header; Excel rows start at 2)
	inputs := make([]models.SampleCreateInput, 0, len(rows)-1)
	for r, row := range rows[1:] {
		// Fully empty rows are not samples
		if emptyRow(row) {
			continue
		}
		line := r + 2

		cell := func(header string) string {
			i := headerIdx[header]
			if i >= len(row) {
				return ""
			}
			return strings.TrimSpace(row[i])
		}

		var in models.SampleCreateInput
		var err error

		// Free-text columns pass raw to Gin binding
		in.OriginCode = cell("origin_code")
		in.RunNumber = cell("run_number")

		// Dropdown columns must match an accepted value before converting to the field type
		if in.City, err = validateDropdown(cell("city"), line, "city",
			columnsByHeader["city"]); err != nil {
			return s.logValidationError(ctx, err)
		}
		if in.CountryCode, err = validateDropdown(cell("country_code"),
			line, "country_code",
			columnsByHeader["country_code"]); err != nil {
			return s.logValidationError(ctx, err)
		}

		var gender string
		if gender, err = validateDropdown(cell("gender"), line, "gender",
			columnsByHeader["gender"]); err != nil {
			return s.logValidationError(ctx, err)
		}
		if gender != "" {
			g := models.Gender(gender)
			in.Gender = &g
		}

		var inNetwork string
		if inNetwork, err = validateDropdown(cell("in_network"), line,
			"in_network", columnsByHeader["in_network"]); err != nil {
			return s.logValidationError(ctx, err)
		}
		if inNetwork != "" {
			b, bErr := strconv.ParseBool(inNetwork)
			if bErr != nil {
				return s.logValidationError(ctx,
					&TableValueError{Row: line, Col: "in_network"})
			}
			in.InNetwork = &b
		}

		// ID columns: label -> accepted value -> uuid
		ids := map[string]*uuid.UUID{
			"origin_id":         &in.OriginID,
			"sample_source_id":  &in.SampleSourceID,
			"microorganism_id":  &in.MicroorganismID,
			"sequencer_id":      &in.SequencerID,
			"laboratory_id":     &in.LaboratoryID,
			"health_service_id": &in.HealthServiceID,
		}
		for header, target := range ids {
			v, vErr := validateDropdown(cell(header), line, header,
				columnsByHeader[header])
			if vErr != nil {
				return s.logValidationError(ctx, vErr)
			}
			if v == "" {
				continue
			}
			u, pErr := uuid.Parse(v)
			if pErr != nil {
				return s.logValidationError(ctx,
					&TableValueError{Row: line, Col: header})
			}
			*target = u
		}

		// Dates parse here; empty required stays zero for Gin binding
		if raw := cell("collection_date"); raw != "" {
			d, ok := parseTableDate(raw)
			if !ok {
				return s.logValidationError(ctx,
					&TableValueError{Row: line, Col: "collection_date"})
			}
			in.CollectionDate = d
		}
		if raw := cell("run_date"); raw != "" {
			d, ok := parseTableDate(raw)
			if !ok {
				return s.logValidationError(ctx,
					&TableValueError{Row: line, Col: "run_date"})
			}
			in.RunDate = d
		}
		if raw := cell("date_of_birth"); raw != "" {
			d, ok := parseTableDate(raw)
			if !ok {
				return s.logValidationError(ctx,
					&TableValueError{Row: line, Col: "date_of_birth"})
			}
			in.DateOfBirth = &d
		}
		inputs = append(inputs, in)
	}

	return inputs, nil
}

func (s *templateService) logValidationError(ctx context.Context,
	err error) ([]models.SampleCreateInput, error) {
	s.Logger.Error(
		"Service Error", logging.ServiceLogging(ctx,
			"TemplateService", "ValidateTemplateTable",
			logging.ValidationError, err,
		)...)

	return nil, err
}
