package repository_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestInitRepositories(t *testing.T) {
	origUserRepo := repository.UserRepo
	origCountryRepo := repository.CountryRepo
	defer func() {
		repository.UserRepo = origUserRepo
		repository.CountryRepo = origCountryRepo
	}()
	
	db := testutils.NewMockDB()
	repository.InitRepositories(db)

	assert.NotEmpty(t, repository.UserRepo)
	assert.NotEmpty(t, repository.CountryRepo)
}
