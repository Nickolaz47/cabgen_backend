package handlers

import (
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/country"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"gorm.io/gorm"
)

func InitRepositories(db *gorm.DB) {
	public.UserRepo = repository.NewUserRepo(db)
	country.CountryRepo = repository.NewCountryRepo(db)
}
