package services_test

import (
	"errors"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/email"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	gomail "gopkg.in/mail.v2"
)

type MockEmailSender struct {
	ShouldFail bool
}

func (m *MockEmailSender) Send(msg *gomail.Message) error {
	if m.ShouldFail {
		return errors.New("simulated send error")
	}
	return nil
}

func TestSendActivationUserEmail(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := testutils.SetupTestRepos()
		repository.InitRepositories(db)

		mockAdminUser := testmodels.NewAdminLoginUser()
		mockAdminUser.Email = "yt4bdzmzze@bwmyga.com"
		db.Create(&mockAdminUser)

		userToActivate := "johndoe"
		sender := MockEmailSender{}

		err := services.SendActivationUserEmail(userToActivate, &sender)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		repository.InitRepositories(mockDB)

		userToActivate := "johndoe"
		sender := email.SMTPEmailSender{
			Username: config.SenderEmail,
			Password: config.SenderPassword,
			Host:     config.SMTPHost,
			Port:     config.SMTPPort,
		}

		err = services.SendActivationUserEmail(userToActivate, &sender)
		assert.Error(t, err)
	})
}
