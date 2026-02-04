package container

import (
	"github.com/CABGenOrg/cabgen_backend/internal/db"
)

func SetupDatabase(driver, dns string,
	modelsToMigrate []any) (*db.GormDatabase, error) {
	database, err := db.NewGormDatabase(driver, dns)
	if err != nil {
		return nil, err
	}

	if err := database.Migrate(modelsToMigrate...); err != nil {
		return nil, err
	}

	return database, nil
}
