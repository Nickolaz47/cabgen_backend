package email

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
	gomail "gopkg.in/mail.v2"
)

func TestSetupEmailMessage(t *testing.T) {
	resultMessage := gomail.NewMessage()

	sender := "test@mail.com"
	recipient := "test2@mail.com"
	subject := "Test"
	body := "<p>Test Body</p>"

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "file.txt")
	testutils.WriteMockEnvFile(t, filePath, "random content")

	emailConfig := EmailConfig{
		Sender:    sender,
		Recipient: recipient,
		Subject:   subject,
		Body:      body,
		File:      filePath,
	}
	setupEmailMessage(resultMessage, emailConfig)

	var buf bytes.Buffer
	_, err := resultMessage.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}

	emailRaw := buf.String()

	assert.Equal(t, emailConfig.Sender, resultMessage.GetHeader("From")[0])
	assert.Equal(t, emailConfig.Recipient, resultMessage.GetHeader("To")[0])
	assert.Equal(t, emailConfig.Subject, resultMessage.GetHeader("Subject")[0])

	assert.Contains(t, emailRaw, body, "expected body in email")
	assert.Contains(t, emailRaw, "Content-Disposition: attachment",
		"expected to have an attachment")
	assert.Contains(t, emailRaw, filepath.Base(filePath),
		"expected file name in email")
}

func TestSendEmail(t *testing.T) {
	sender := "test@mail.com"
	recipient := "test2@mail.com"
	subject := "Test"
	body := "<p>Test Body</p>"

	emailConfig := EmailConfig{
		Sender: sender, Recipient: recipient,
		Subject: subject, Body: body,
	}

	t.Run("SendEmail - Success", func(t *testing.T) {
		mockSender := &mocks.MockEmailSender{ShouldFail: false}

		err := SendEmail(emailConfig, mockSender)
		assert.NoError(t, err, "Expected to send email correctly")
	})

	t.Run("SendEmail - Error", func(t *testing.T) {
		mockSender := &mocks.MockEmailSender{ShouldFail: true}
		err := SendEmail(emailConfig, mockSender)

		assert.Error(t, err, "Expected to failed to send email")
	})
}

func TestCreateDefaultSender(t *testing.T) {
	origSenderEmail := config.SenderEmail
	origSenderPassword := config.SenderPassword
	origSMTPHost := config.SMTPHost
	origSMTPPort := config.SMTPPort

	expectedSenderEmail := "cabgen@mail.com"
	expectedSenderPassword := "sender_password"
	expectedSMTPHost := "smtp.mail.com"
	expectedSMTPPort := 587

	config.SenderEmail = expectedSenderEmail
	config.SenderPassword = expectedSenderPassword
	config.SMTPHost = expectedSMTPHost
	config.SMTPPort = expectedSMTPPort
	defer func() {
		config.SenderEmail = origSenderEmail
		config.SenderPassword = origSenderPassword
		config.SMTPHost = origSMTPHost
		config.SMTPPort = origSMTPPort
	}()

	result := CreateDefaultSender()

	assert.NotEmpty(t, result)
	assert.Equal(t, expectedSenderEmail, result.Username)
	assert.Equal(t, expectedSenderPassword, result.Password)
	assert.Equal(t, expectedSMTPHost, result.Host)
	assert.Equal(t, expectedSMTPPort, result.Port)
}
