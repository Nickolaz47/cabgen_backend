package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewAuditRepo(t *testing.T) {
	db := testutils.NewMockDB()
	auditRepo := repositories.NewAuditRepository(db)

	assert.NotEmpty(t, auditRepo)
}

func TestGetAuditLogs(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	auditRepo := repositories.NewAuditRepository(db)

	loginAudit := testmodels.NewAudit(models.AuditEventLogin, "10.0.0.1",
		"{}", 200)
	loginAudit.CreatedAt = time.Date(2026, 1, 2, 10, 30, 0, 0, time.UTC)

	failedAudit := testmodels.NewAudit(models.AuditEventLoginFailed,
		"10.0.0.99", "{}", 401)
	failedAudit.CreatedAt = time.Date(2026, 1, 2, 0, 0, 30, 0, time.UTC)
	failedAudit.UserID = nil
	failedAudit.User = nil

	deleteAudit := testmodels.NewAudit(models.AuditEventAdminUsersDelete,
		"192.168.0.1", "{}", 200)
	deleteAudit.CreatedAt = time.Date(2026, 1, 3, 8, 0, 0, 0, time.UTC)

	db.Create(&loginAudit)
	db.Create(&failedAudit)
	db.Create(&deleteAudit)

	createdAt := func(audit models.Audit) time.Time {
		return audit.CreatedAt
	}

	t.Run("Success - Without Filter", func(t *testing.T) {
		auditLogs, err := auditRepo.GetAuditLogs(ctx, models.AuditFilter{})

		assert.NoError(t, err)
		assert.Len(t, auditLogs, 3)
		assert.Equal(t, deleteAudit.ID, auditLogs[0].ID)
		assert.Equal(t, loginAudit.ID, auditLogs[1].ID)
		assert.Equal(t, failedAudit.ID, auditLogs[2].ID)
		assert.NotZero(t, createdAt(auditLogs[0]))
	})

	t.Run("Success - Filter by Event", func(t *testing.T) {
		auditLogs, err := auditRepo.GetAuditLogs(ctx, models.AuditFilter{
			Event: models.AuditEventLogin,
		})

		assert.NoError(t, err)
		assert.Len(t, auditLogs, 1)
		assert.Equal(t, loginAudit.ID, auditLogs[0].ID)
		assert.Equal(t, models.AuditEventLogin, auditLogs[0].Event)
	})

	t.Run("Success - Filter by Source Partial IP", func(t *testing.T) {
		auditLogs, err := auditRepo.GetAuditLogs(ctx, models.AuditFilter{
			Source: "10.0.",
		})

		assert.NoError(t, err)
		assert.Len(t, auditLogs, 2)

		ids := []string{auditLogs[0].ID.String(), auditLogs[1].ID.String()}
		assert.Contains(t, ids, loginAudit.ID.String())
		assert.Contains(t, ids, failedAudit.ID.String())
	})

	t.Run("Success - Filter by Status", func(t *testing.T) {
		auditLogs, err := auditRepo.GetAuditLogs(ctx, models.AuditFilter{
			Status: 401,
		})

		assert.NoError(t, err)
		assert.Len(t, auditLogs, 1)
		assert.Equal(t, failedAudit.ID, auditLogs[0].ID)
		assert.Equal(t, 401, auditLogs[0].Status)
	})

	t.Run("Success - Filter by Date", func(t *testing.T) {
		date := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
		auditLogs, err := auditRepo.GetAuditLogs(ctx, models.AuditFilter{
			Date: &date,
		})

		assert.NoError(t, err)
		assert.Len(t, auditLogs, 2)

		ids := []string{auditLogs[0].ID.String(), auditLogs[1].ID.String()}
		assert.Contains(t, ids, loginAudit.ID.String())
		assert.Contains(t, ids, failedAudit.ID.String())
		assert.NotContains(t, ids, deleteAudit.ID.String())
	})

	t.Run("Success - Filter by User ID", func(t *testing.T) {
		admin := testmodels.NewAdminLoginUser()
		auditLogs, err := auditRepo.GetAuditLogs(ctx, models.AuditFilter{
			UserID: &admin.ID,
		})

		assert.NoError(t, err)
		assert.Len(t, auditLogs, 2)

		ids := []string{auditLogs[0].ID.String(), auditLogs[1].ID.String()}
		assert.Contains(t, ids, loginAudit.ID.String())
		assert.Contains(t, ids, deleteAudit.ID.String())
	})

	t.Run("Success - Combined Filters", func(t *testing.T) {
		auditLogs, err := auditRepo.GetAuditLogs(ctx, models.AuditFilter{
			Event:  models.AuditEventLogin,
			Status: 200,
		})

		assert.NoError(t, err)
		assert.Len(t, auditLogs, 1)
		assert.Equal(t, loginAudit.ID, auditLogs[0].ID)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockAuditRepo := repositories.NewAuditRepository(mockDB)
		auditLogs, err := mockAuditRepo.GetAuditLogs(ctx, models.AuditFilter{})

		assert.Empty(t, auditLogs)
		assert.Error(t, err)
	})
}

func TestCreateAudit(t *testing.T) {
	ctx := context.Background()

	db := testutils.NewMockDB()
	auditRepo := repositories.NewAuditRepository(db)

	audit := testmodels.NewAudit(models.AuditEventLogin, "10.0.0.1", "{}", 200)
	auditWithoutUser := testmodels.NewAudit(
		models.AuditEventLoginFailed, "10.0.0.1", "{}", 401)
	auditWithoutUser.UserID = nil
	auditWithoutUser.User = nil

	t.Run("Success", func(t *testing.T) {
		err := auditRepo.CreateAudit(ctx, &audit)
		assert.NoError(t, err)

		var result models.Audit
		err = db.Where("id = ?", audit.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, audit.ID, result.ID)
		assert.Equal(t, audit.Event, result.Event)
		assert.Equal(t, audit.Source, result.Source)
		assert.Equal(t, audit.Status, result.Status)
		assert.Equal(t, audit.Metadata, result.Metadata)
		assert.Equal(t, audit.UserID, result.UserID)
		assert.NotZero(t, result.CreatedAt)
	})

	t.Run("Success - Without User", func(t *testing.T) {
		err := auditRepo.CreateAudit(ctx, &auditWithoutUser)
		assert.NoError(t, err)

		var result models.Audit
		err = db.Where("id = ?", auditWithoutUser.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Nil(t, result.UserID)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockAuditRepo := repositories.NewAuditRepository(mockDB)
		err = mockAuditRepo.CreateAudit(ctx, &models.Audit{})

		assert.Error(t, err)
	})
}

func TestAuditUserDeletionKeepsAuditLog(t *testing.T) {
	ctx := context.Background()
	db := testutils.NewMockDB()
	auditRepo := repositories.NewAuditRepository(db)

	admin := testmodels.NewAdminLoginUser()
	audit := testmodels.NewAudit(models.AuditEventLogin, "10.0.0.1", "{}", 200)

	db.Create(&audit)

	db.Delete(&admin)

	auditLogs, err := auditRepo.GetAuditLogs(ctx,
		models.AuditFilter{})

	assert.NoError(t, err)
	assert.Len(t, auditLogs, 1)
	assert.Nil(t, auditLogs[0].UserID)
	assert.Nil(t, auditLogs[0].User)
}
