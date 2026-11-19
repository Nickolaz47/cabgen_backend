package repositories_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const laboratoriesMigrationsDir = "../db/migrations"

func newLaboratoryTestDB(t *testing.T) db.Database {
	t.Helper()

	database, cleanup, err := testutils.NewMockPostgresDB(
		context.Background(), laboratoriesMigrationsDir)
	if err != nil {
		t.Fatalf("failed to start test database: %v", err)
	}
	t.Cleanup(cleanup)

	return database
}

func seedLaboratory(t *testing.T, database db.Database, ID uuid.UUID,
	name, abbreviation string, isActive bool) {
	t.Helper()

	_, err := database.SQLDB().ExecContext(context.Background(),
		`INSERT INTO laboratories (id, name, abbreviation, is_active)
         VALUES ($1, $2, $3, $4)`,
		ID, name, abbreviation, isActive)
	if err != nil {
		t.Fatalf("failed to seed laboratory: %v", err)
	}
}

func dropLaboratoriesTable(t *testing.T, database db.Database) {
	t.Helper()

	_, err := database.SQLDB().ExecContext(context.Background(),
		"DROP TABLE IF EXISTS laboratories CASCADE")
	if err != nil {
		t.Fatalf("failed to drop laboratories table: %v", err)
	}
}

func assertLaboratory(t *testing.T, expected, actual models.Laboratory) {
	t.Helper()

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Abbreviation, actual.Abbreviation)
	assert.Equal(t, expected.IsActive, actual.IsActive)
}

func TestNewLaboratoryRepo(t *testing.T) {
	database := newLaboratoryTestDB(t)
	labRepo := repositories.NewLaboratoryRepo(database)

	assert.NotEmpty(t, labRepo)
}

func TestGetLaboratories(t *testing.T) {
	ctx := context.Background()
	database := newLaboratoryTestDB(t)
	labRepo := repositories.NewLaboratoryRepo(database)

	labID := uuid.New()
	seedLaboratory(t, database, labID, "Laboratory 1", "L1", true)

	t.Run("Success", func(t *testing.T) {
		labs, err := labRepo.GetLaboratories(ctx)

		assert.NoError(t, err)
		assert.Len(t, labs, 1)
		assertLaboratory(t, models.Laboratory{
			ID: labID, Name: "Laboratory 1", Abbreviation: "L1", IsActive: true,
		}, labs[0])
	})

	t.Run("Error", func(t *testing.T) {
		dropLaboratoriesTable(t, database)

		labs, err := labRepo.GetLaboratories(ctx)

		assert.Empty(t, labs)
		assert.Error(t, err)
	})
}

func TestGetActiveLaboratories(t *testing.T) {
	ctx := context.Background()
	database := newLaboratoryTestDB(t)
	labRepo := repositories.NewLaboratoryRepo(database)

	labID := uuid.New()
	inactiveID := uuid.New()
	seedLaboratory(t, database, labID, "Laboratory 1", "L1", true)
	seedLaboratory(t, database, inactiveID, "Laboratory 2", "L2", false)

	t.Run("Success", func(t *testing.T) {
		labs, err := labRepo.GetActiveLaboratories(ctx)

		assert.NoError(t, err)
		assert.Len(t, labs, 1)
		assertLaboratory(t, models.Laboratory{
			ID: labID, Name: "Laboratory 1", Abbreviation: "L1", IsActive: true,
		}, labs[0])
	})

	t.Run("Error", func(t *testing.T) {
		dropLaboratoriesTable(t, database)

		labs, err := labRepo.GetActiveLaboratories(ctx)

		assert.Empty(t, labs)
		assert.Error(t, err)
	})
}

func TestGetLaboratoryByID(t *testing.T) {
	ctx := context.Background()
	database := newLaboratoryTestDB(t)
	labRepo := repositories.NewLaboratoryRepo(database)

	labID := uuid.New()
	seedLaboratory(t, database, labID, "Laboratory 1", "L1", true)

	t.Run("Success", func(t *testing.T) {
		resultLab, err := labRepo.GetLaboratoryByID(ctx, labID)

		assert.NoError(t, err)
		assertLaboratory(t, models.Laboratory{
			ID: labID, Name: "Laboratory 1", Abbreviation: "L1", IsActive: true,
		}, *resultLab)
	})

	t.Run("Error - Record not found", func(t *testing.T) {
		lab, err := labRepo.GetLaboratoryByID(ctx, uuid.New())

		assert.Nil(t, lab)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "no rows in result set")
	})

	t.Run("Error - DB error", func(t *testing.T) {
		dropLaboratoriesTable(t, database)

		lab, err := labRepo.GetLaboratoryByID(ctx, labID)

		assert.Nil(t, lab)
		assert.Error(t, err)
	})
}

func TestGetLaboratoriesByNameOrAbbreviation(t *testing.T) {
	ctx := context.Background()
	database := newLaboratoryTestDB(t)
	labRepo := repositories.NewLaboratoryRepo(database)

	labID := uuid.New()
	lab2ID := uuid.New()
	seedLaboratory(t, database, labID, "Laboratory Alpha", "LAB-1", true)
	seedLaboratory(t, database, lab2ID, "Laboratory Beta", "LAB-2", true)

	t.Run("Success - Name", func(t *testing.T) {
		labs, err := labRepo.GetLaboratoriesByNameOrAbbreviation(ctx, "beta")

		assert.NoError(t, err)
		assert.Len(t, labs, 1)
		assertLaboratory(t, models.Laboratory{
			ID: lab2ID, Name: "Laboratory Beta", Abbreviation: "LAB-2",
			IsActive: true,
		}, labs[0])
	})

	t.Run("Success - Abbreviation", func(t *testing.T) {
		labs, err := labRepo.GetLaboratoriesByNameOrAbbreviation(ctx, "lab")

		assert.NoError(t, err)
		assert.Len(t, labs, 2)
		assertLaboratory(t, models.Laboratory{
			ID: labID, Name: "Laboratory Alpha", Abbreviation: "LAB-1",
			IsActive: true,
		}, labs[0])
		assertLaboratory(t, models.Laboratory{
			ID: lab2ID, Name: "Laboratory Beta", Abbreviation: "LAB-2",
			IsActive: true,
		}, labs[1])
	})

	t.Run("Error", func(t *testing.T) {
		dropLaboratoriesTable(t, database)

		labs, err := labRepo.GetLaboratoriesByNameOrAbbreviation(ctx, "lab")

		assert.Empty(t, labs)
		assert.Error(t, err)
	})
}

func TestGetLaboratoryDuplicate(t *testing.T) {
	ctx := context.Background()
	database := newLaboratoryTestDB(t)
	labRepo := repositories.NewLaboratoryRepo(database)

	labID := uuid.New()
	seedLaboratory(t, database, labID, "Laboratory 1", "L1", true)
	expectedLab := models.Laboratory{
		ID: labID, Name: "Laboratory 1", Abbreviation: "L1", IsActive: true,
	}

	t.Run("Success - With ID", func(t *testing.T) {
		lab, err := labRepo.GetLaboratoryDuplicate(ctx, "Laboratory 1",
			uuid.New())

		assert.NoError(t, err)
		assertLaboratory(t, expectedLab, *lab)
	})

	t.Run("Success - Without ID", func(t *testing.T) {
		lab, err := labRepo.GetLaboratoryDuplicate(ctx, "Laboratory 1",
			uuid.UUID{})

		assert.NoError(t, err)
		assertLaboratory(t, expectedLab, *lab)
	})

	t.Run("Success - Own ID excluded", func(t *testing.T) {
		lab, err := labRepo.GetLaboratoryDuplicate(ctx, "Laboratory 1", labID)

		assert.Nil(t, lab)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "no rows in result set")
	})

	t.Run("Error - Record not found", func(t *testing.T) {
		lab, err := labRepo.GetLaboratoryDuplicate(ctx,
			"Laboratory 2", uuid.UUID{})

		assert.Nil(t, lab)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "no rows in result set")
	})

	t.Run("Error - DB error", func(t *testing.T) {
		dropLaboratoriesTable(t, database)

		lab, err := labRepo.GetLaboratoryDuplicate(ctx, "Laboratory 1",
			uuid.New())

		assert.Nil(t, lab)
		assert.Error(t, err)
	})
}

func TestCreateLaboratory(t *testing.T) {
	ctx := context.Background()
	database := newLaboratoryTestDB(t)
	labRepo := repositories.NewLaboratoryRepo(database)

	t.Run("Success", func(t *testing.T) {
		lab := models.Laboratory{
			Name: "Laboratory 1", Abbreviation: "L1", IsActive: true,
		}

		err := labRepo.CreateLaboratory(ctx, &lab)
		assert.NoError(t, err)

		labs, err := labRepo.GetLaboratories(ctx)

		assert.NoError(t, err)
		assert.Len(t, labs, 1)
		assert.Equal(t, "Laboratory 1", labs[0].Name)
		assert.Equal(t, "L1", labs[0].Abbreviation)
		assert.True(t, labs[0].IsActive)
		assert.NotEqual(t, uuid.Nil, labs[0].ID,
			"ID should be generated by the database")
	})

	t.Run("Error", func(t *testing.T) {
		dropLaboratoriesTable(t, database)

		lab := models.Laboratory{Name: "Laboratory 2", Abbreviation: "L2"}
		err := labRepo.CreateLaboratory(ctx, &lab)

		assert.Error(t, err)
	})
}

func TestUpdateLaboratory(t *testing.T) {
	ctx := context.Background()
	database := newLaboratoryTestDB(t)
	labRepo := repositories.NewLaboratoryRepo(database)

	labID := uuid.New()
	seedLaboratory(t, database, labID, "Laboratory 1", "L1", true)

	t.Run("Success", func(t *testing.T) {
		labToUpdate := models.Laboratory{
			ID:           labID,
			Name:         "Laboratory 1 Updated",
			Abbreviation: "L1U",
			IsActive:     false,
		}

		err := labRepo.UpdateLaboratory(ctx, &labToUpdate)
		assert.NoError(t, err)

		result, err := labRepo.GetLaboratoryByID(ctx, labID)

		assert.NoError(t, err)
		assertLaboratory(t, labToUpdate, *result)
	})

	t.Run("Error", func(t *testing.T) {
		dropLaboratoriesTable(t, database)

		lab := models.Laboratory{ID: labID, Name: "Laboratory 1"}
		err := labRepo.UpdateLaboratory(ctx, &lab)

		assert.Error(t, err)
	})
}

func TestDeleteLaboratory(t *testing.T) {
	ctx := context.Background()
	database := newLaboratoryTestDB(t)
	labRepo := repositories.NewLaboratoryRepo(database)

	labID := uuid.New()
	seedLaboratory(t, database, labID, "Laboratory 1", "L1", true)

	t.Run("Success", func(t *testing.T) {
		lab := models.Laboratory{ID: labID, Name: "Laboratory 1"}

		err := labRepo.DeleteLaboratory(ctx, &lab)
		assert.NoError(t, err)

		result, err := labRepo.GetLaboratoryByID(ctx, labID)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "no rows in result set")
	})

	t.Run("Error", func(t *testing.T) {
		dropLaboratoriesTable(t, database)

		lab := models.Laboratory{ID: labID, Name: "Laboratory 1"}
		err := labRepo.DeleteLaboratory(ctx, &lab)

		assert.Error(t, err)
	})
}
