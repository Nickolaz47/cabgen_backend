package container

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/db"
)

func SetupDatabase(ctx context.Context, driver, dns string,
	modelsToMigrate []any) (db.Database, error) {
	database, err := db.NewDatabase(ctx, driver, dns)
	if err != nil {
		return nil, err
	}

	if err := database.Migrate(modelsToMigrate...); err != nil {
		return nil, err
	}

	return database, nil
}
