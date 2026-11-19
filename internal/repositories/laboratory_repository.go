package repositories

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/db/sqlc"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
)

type LaboratoryRepository interface {
	GetLaboratories(ctx context.Context) ([]models.Laboratory, error)
	GetActiveLaboratories(ctx context.Context) ([]models.Laboratory, error)
	GetLaboratoryByID(ctx context.Context, ID uuid.UUID) (*models.Laboratory, error)
	GetLaboratoriesByNameOrAbbreviation(ctx context.Context, input string) ([]models.Laboratory, error)
	GetLaboratoryDuplicate(ctx context.Context, name string, ID uuid.UUID) (*models.Laboratory, error)
	CreateLaboratory(ctx context.Context, lab *models.Laboratory) error
	UpdateLaboratory(ctx context.Context, lab *models.Laboratory) error
	DeleteLaboratory(ctx context.Context, lab *models.Laboratory) error
}

type laboratoryRepo struct {
	DB db.Database
}

func NewLaboratoryRepo(db db.Database) LaboratoryRepository {
	return &laboratoryRepo{DB: db}
}

func (r *laboratoryRepo) GetLaboratories(
	ctx context.Context) ([]models.Laboratory, error) {
	queries := r.DB.Queries()

	sqlcLabs, err := queries.GetLaboratories(ctx)
	if err != nil {
		return nil, err
	}

	labs := make([]models.Laboratory, len(sqlcLabs))
	for i, sl := range sqlcLabs {
		labs[i] = r.toDomain(sl)
	}

	return labs, nil
}

func (r *laboratoryRepo) GetActiveLaboratories(
	ctx context.Context) ([]models.Laboratory, error) {
	queries := r.DB.Queries()

	sqlcLabs, err := queries.GetActiveLaboratories(ctx)
	if err != nil {
		return nil, err
	}

	labs := make([]models.Laboratory, len(sqlcLabs))
	for i, sl := range sqlcLabs {
		labs[i] = r.toDomain(sl)
	}

	return labs, nil
}

func (r *laboratoryRepo) GetLaboratoryByID(
	ctx context.Context, ID uuid.UUID) (*models.Laboratory, error) {
	queries := r.DB.Queries()

	sqlcLab, err := queries.GetLaboratoryByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	lab := r.toDomain(sqlcLab)
	return &lab, nil
}

func (r *laboratoryRepo) GetLaboratoriesByNameOrAbbreviation(
	ctx context.Context, input string) ([]models.Laboratory, error) {
	queries := r.DB.Queries()

	sqlcLabs, err := queries.GetLaboratoriesByNameOrAbbreviation(ctx, input)
	if err != nil {
		return nil, err
	}

	labs := make([]models.Laboratory, len(sqlcLabs))
	for i, sl := range sqlcLabs {
		labs[i] = r.toDomain(sl)
	}

	return labs, nil
}

func (r *laboratoryRepo) GetLaboratoryDuplicate(ctx context.Context,
	name string, ID uuid.UUID) (*models.Laboratory, error) {
	queries := r.DB.Queries()

	sqlcLab, err := queries.GetLaboratoryDuplicate(ctx,
		sqlc.GetLaboratoryDuplicateParams{Name: name, ExcludeID: ID})
	if err != nil {
		return nil, err
	}

	lab := r.toDomain(sqlcLab)
	return &lab, nil
}

func (r *laboratoryRepo) CreateLaboratory(ctx context.Context,
	lab *models.Laboratory) error {
	queries := r.DB.Queries()

	_, err := queries.CreateLaboratory(ctx, sqlc.CreateLaboratoryParams{
		Name: lab.Name, Abbreviation: lab.Abbreviation, IsActive: lab.IsActive,
	})

	return err
}

func (r *laboratoryRepo) UpdateLaboratory(ctx context.Context,
	lab *models.Laboratory) error {
	queries := r.DB.Queries()

	_, err := queries.UpdateLaboratory(ctx, sqlc.UpdateLaboratoryParams{
		ID:   lab.ID,
		Name: lab.Name, Abbreviation: lab.Abbreviation, IsActive: lab.IsActive,
	})

	return err
}

func (r *laboratoryRepo) DeleteLaboratory(ctx context.Context,
	lab *models.Laboratory) error {
	queries := r.DB.Queries()

	return queries.DeleteLaboratory(ctx, lab.ID)
}

func (r *laboratoryRepo) toDomain(sqlcModel sqlc.Laboratory) models.Laboratory {
	return models.Laboratory{
		ID:           sqlcModel.ID,
		Name:         sqlcModel.Name,
		Abbreviation: sqlcModel.Abbreviation,
		IsActive:     sqlcModel.IsActive,
	}
}
