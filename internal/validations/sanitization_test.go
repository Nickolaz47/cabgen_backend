package validations

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string { return &s }

func TestSanitizeInput_UserRegisterInput(t *testing.T) {
	input := models.UserRegisterInput{
		Name:            "  John  Doe  ",
		Username:        "  JohnDoe  ",
		Email:           "  JOHN@EXAMPLE.COM  ",
		ConfirmEmail:    "  JOHN@EXAMPLE.COM  ",
		Password:        "  secret123  ",
		ConfirmPassword: "  secret123  ",
		CountryCode:     "  BR  ",
		Interest:        strPtr("  engineering  "),
	}
	SanitizeInput(&input)
	assert.Equal(t, "John  Doe", input.Name)
	assert.Equal(t, "johndoe", input.Username)
	assert.Equal(t, "john@example.com", input.Email)
	assert.Equal(t, "john@example.com", input.ConfirmEmail)
	assert.Equal(t, "  secret123  ", input.Password)
	assert.Equal(t, "  secret123  ", input.ConfirmPassword)
	assert.Equal(t, "BR", input.CountryCode)
	assert.Equal(t, "engineering", *input.Interest)
}

func TestSanitizeInput_LoginInput(t *testing.T) {
	input := models.LoginInput{
		Username: "  JohnDoe  ",
		Password: "  secret123  ",
	}
	SanitizeInput(&input)
	assert.Equal(t, "johndoe", input.Username)
	assert.Equal(t, "  secret123  ", input.Password)
}

func TestSanitizeInput_UserUpdateInput_NilFields(t *testing.T) {
	input := models.UserUpdateInput{}
	SanitizeInput(&input)
	assert.Nil(t, input.Name)
	assert.Nil(t, input.Username)
}

func TestSanitizeInput_AdminUserCreateInput(t *testing.T) {
	input := models.AdminUserCreateInput{
		Name:        "  Jane  Smith  ",
		Username:    "  JaneSmith  ",
		Email:       "  JANE@EXAMPLE.COM  ",
		Password:    "  secret123  ",
		CountryCode: "  US  ",
	}
	SanitizeInput(&input)
	assert.Equal(t, "Jane  Smith", input.Name)
	assert.Equal(t, "janesmith", input.Username)
	assert.Equal(t, "jane@example.com", input.Email)
	assert.Equal(t, "  secret123  ", input.Password)
	assert.Equal(t, "US", input.CountryCode)
}

func TestSanitizeInput_CreateTicketInput(t *testing.T) {
	input := models.CreateTicketInput{
		Name:        "  John  Doe  ",
		Email:       "  JOHN@EXAMPLE.COM  ",
		Institution: "  University  ",
		Subject:     "  Test  Subject  ",
		Message:     "  Test  Message  ",
	}
	SanitizeInput(&input)
	assert.Equal(t, "John  Doe", input.Name)
	assert.Equal(t, "john@example.com", input.Email)
	assert.Equal(t, "University", input.Institution)
	assert.Equal(t, "Test  Subject", input.Subject)
	assert.Equal(t, "Test  Message", input.Message)
}

func TestSanitizeInput_CreateTicketInput_XSS(t *testing.T) {
	input := models.CreateTicketInput{
		Name:        "John Doe",
		Email:       "john@example.com",
		Institution: "University",
		Subject:     `<script>alert("xss")</script>Subject`,
		Message:     `<img onerror="alert(1)" src="x">Message`,
	}
	SanitizeInput(&input)
	assert.Equal(t, "Subject", input.Subject)
	assert.Equal(t, "Message", input.Message)
}

func TestSanitizeInput_AdminSampleCreateInput(t *testing.T) {
	input := models.AdminSampleCreateInput{
		OriginCode:  "  ABC  ",
		RunNumber:   "  123  ",
		CountryCode: "  BR  ",
		City:        strPtr("  Sao Paulo  "),
	}
	SanitizeInput(&input)
	assert.Equal(t, "ABC", input.OriginCode)
	assert.Equal(t, "123", input.RunNumber)
	assert.Equal(t, "BR", input.CountryCode)
	assert.Equal(t, "Sao Paulo", *input.City)
}

func TestSanitizeInput_SampleUpdateInput_NilCity(t *testing.T) {
	input := models.SampleUpdateInput{
		OriginCode: strPtr("  ABC  "),
		City:       nil,
	}
	SanitizeInput(&input)
	assert.Equal(t, "ABC", *input.OriginCode)
	assert.Nil(t, input.City)
}

func TestSanitizeInput_SequencerCreateInput(t *testing.T) {
	input := models.SequencerCreateInput{
		Model: "  MiSeq  ",
		Brand: "  Illumina  ",
	}
	SanitizeInput(&input)
	assert.Equal(t, "MiSeq", input.Model)
	assert.Equal(t, "Illumina", input.Brand)
}

func TestSanitizeInput_LaboratoryCreateInput(t *testing.T) {
	input := models.LaboratoryCreateInput{
		Name:         "  Lab  One  ",
		Abbreviation: "  L1  ",
	}
	SanitizeInput(&input)
	assert.Equal(t, "Lab  One", input.Name)
	assert.Equal(t, "L1", input.Abbreviation)
}

func TestSanitizeInput_MicroorganismCreateInput(t *testing.T) {
	input := models.MicroorganismCreateInput{
		Species: "  E.  coli  ",
		Variety: map[string]string{"en": "  E.  coli  ", "pt": "  E.  coli  "},
	}
	SanitizeInput(&input)
	assert.Equal(t, "E.  coli", input.Species)
	assert.Equal(t, "E.  coli", input.Variety["en"])
}

func TestSanitizeInput_HealthServiceCreateInput(t *testing.T) {
	input := models.HealthServiceCreateInput{
		Name:         "  Hospital  ",
		Type:         models.HealthServiceType("  Public  "),
		CountryCode:  "  BR  ",
		City:         strPtr("  Rio  "),
		Contactant:   strPtr("  Dr.  Silva  "),
		ContactEmail: strPtr("  DR@HOSPITAL.COM  "),
		ContactPhone: strPtr("  +5511999999999  "),
	}
	SanitizeInput(&input)
	assert.Equal(t, "Hospital", input.Name)
	assert.Equal(t, models.HealthServiceType("Public"), input.Type)
	assert.Equal(t, "BR", input.CountryCode)
	assert.Equal(t, "Rio", *input.City)
	assert.Equal(t, "Dr.  Silva", *input.Contactant)
	assert.Equal(t, "dr@hospital.com", *input.ContactEmail)
	assert.Equal(t, "+5511999999999", *input.ContactPhone)
}

func TestSanitizeInput_CountryCreateInput(t *testing.T) {
	input := models.CountryCreateInput{
		Code:  "  BR  ",
		Names: map[string]string{"en": "  Brazil  ", "pt": "  Brasil  "},
	}
	SanitizeInput(&input)
	assert.Equal(t, "BR", input.Code)
	assert.Equal(t, "Brazil", input.Names["en"])
	assert.Equal(t, "Brasil", input.Names["pt"])
}

func TestSanitizeInput_OriginCreateInput(t *testing.T) {
	input := models.OriginCreateInput{
		Names: map[string]string{"en": "  Hospital  ", "pt": "  Hospital  "},
	}
	SanitizeInput(&input)
	assert.Equal(t, "Hospital", input.Names["en"])
	assert.Equal(t, "Hospital", input.Names["pt"])
}

func TestSanitizeInput_SampleSourceCreateInput(t *testing.T) {
	input := models.SampleSourceCreateInput{
		Names:  map[string]string{"en": "  Blood  "},
		Groups: map[string]string{"en": "  Clinical  "},
	}
	SanitizeInput(&input)
	assert.Equal(t, "Blood", input.Names["en"])
	assert.Equal(t, "Clinical", input.Groups["en"])
}

func TestSanitizeInput_ForgotPasswordInput(t *testing.T) {
	input := models.ForgotPasswordInput{
		Email: "  JOHN@EXAMPLE.COM  ",
	}
	SanitizeInput(&input)
	assert.Equal(t, "john@example.com", input.Email)
}

func TestSanitizeInput_ResetPasswordInput_PasswordUntouched(t *testing.T) {
	input := models.ResetPasswordInput{
		Token:           "  abc123  ",
		NewPassword:     "  secret123  ",
		ConfirmPassword: "  secret123  ",
	}
	SanitizeInput(&input)
	assert.Equal(t, "abc123", input.Token)
	assert.Equal(t, "  secret123  ", input.NewPassword)
	assert.Equal(t, "  secret123  ", input.ConfirmPassword)
}

func TestSanitizeInput_RequestEmailUpdateInput(t *testing.T) {
	input := models.RequestEmailUpdateInput{
		NewEmail:        "  NEW@EXAMPLE.COM  ",
		ConfirmNewEmail: "  NEW@EXAMPLE.COM  ",
	}
	SanitizeInput(&input)
	assert.Equal(t, "new@example.com", input.NewEmail)
	assert.Equal(t, "new@example.com", input.ConfirmNewEmail)
}

func TestSanitizeInput_ConfirmEmailUpdateInput(t *testing.T) {
	input := models.ConfirmEmailUpdateInput{
		Token: "  abc123  ",
	}
	SanitizeInput(&input)
	assert.Equal(t, "abc123", input.Token)
}

func TestSanitizeInput_AdminAnalysisUpdateInput(t *testing.T) {
	input := models.AdminAnalysisUpdateInput{
		FastQC1:        strPtr("  qc1.html  "),
		FastQC2:        strPtr("  qc2.html  "),
		ResultsZipPath: strPtr("  results.zip  "),
		ErrorMessage:   strPtr("  error  "),
	}
	SanitizeInput(&input)
	assert.Equal(t, "qc1.html", *input.FastQC1)
	assert.Equal(t, "qc2.html", *input.FastQC2)
	assert.Equal(t, "results.zip", *input.ResultsZipPath)
	assert.Equal(t, "error", *input.ErrorMessage)
}

func TestSanitizeInput_AnalysisCreateInput_NoStringFields(t *testing.T) {
	input := models.AnalysisCreateInput{}
	SanitizeInput(&input)
}

func TestSanitizeInput_AdminAnalysisCreateInput_NoStringFields(t *testing.T) {
	input := models.AdminAnalysisCreateInput{}
	SanitizeInput(&input)
}

func TestSanitizeInput_AnalysisTSVDownloadInput_NoStringFields(t *testing.T) {
	input := models.AnalysisTSVDownloadInput{}
	SanitizeInput(&input)
}

func TestSanitizeInput_UpdatePasswordInput_PasswordsUntouched(t *testing.T) {
	input := models.UpdatePasswordInput{
		CurrentPassword: "  oldpass  ",
		NewPassword:     "  newpass  ",
		ConfirmPassword: "  newpass  ",
	}
	SanitizeInput(&input)
	assert.Equal(t, "  oldpass  ", input.CurrentPassword)
	assert.Equal(t, "  newpass  ", input.NewPassword)
	assert.Equal(t, "  newpass  ", input.ConfirmPassword)
}

func TestSanitizeInput_AdminUserUpdateInput_NilPassword(t *testing.T) {
	input := models.AdminUserUpdateInput{
		Name:     strPtr("  Jane  "),
		Username: strPtr("  JaneDoe  "),
		Email:    strPtr("  JANE@EXAMPLE.COM  "),
		Password: nil,
	}
	SanitizeInput(&input)
	assert.Equal(t, "Jane", *input.Name)
	assert.Equal(t, "janedoe", *input.Username)
	assert.Equal(t, "jane@example.com", *input.Email)
	assert.Nil(t, input.Password)
}

func TestSanitizeInput_SampleAttachmentInput(t *testing.T) {
	input := models.SampleAttachmentInput{
		Fastq1: strPtr("  /path/to/file1.fq  "),
		Fastq2: strPtr("  /path/to/file2.fq  "),
		Fasta:  strPtr("  /path/to/file.fasta  "),
	}
	SanitizeInput(&input)
	assert.Equal(t, "/path/to/file1.fq", *input.Fastq1)
	assert.Equal(t, "/path/to/file2.fq", *input.Fastq2)
	assert.Equal(t, "/path/to/file.fasta", *input.Fasta)
}
