package services_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
)

func TestSendActivationUserEmail(t *testing.T) {
	ctx := context.Background()
	userToActivate := "johndoe"
	adminUser := testmodels.NewAdminLoginUser()

	t.Run("Success", func(t *testing.T) {
		userRepo := mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter) ([]models.User, error) {
				return []models.User{adminUser}, nil
			},
		}
		sender := mocks.MockEmailSender{}

		svc := services.NewEmailService(&userRepo, &sender)
		err := svc.SendActivationUserEmail(ctx, userToActivate)

		assert.NoError(t, err)
	})

	t.Run("Error - Get Admins", func(t *testing.T) {
		userRepo := mocks.MockUserRepository{
			GetUsersFunc: func(ctx context.Context,
				filter models.AdminUserFilter) ([]models.User, error) {
				return nil, gorm.ErrInvalidTransaction
			},
		}
		sender := mocks.MockEmailSender{}

		svc := services.NewEmailService(&userRepo, &sender)
		err := svc.SendActivationUserEmail(ctx, userToActivate)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Failed to get admins:")
	})
}
