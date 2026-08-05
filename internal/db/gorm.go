package db

import (
	"fmt"
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	MaxConnectionRetries = 10
	ConnectionRetryDelay = 2 * time.Second
)

type GormDatabase struct {
	db *gorm.DB
}

func NewGormDatabase(driver, dsn string) (*GormDatabase, error) {
	newLogger := logger.New(
		log.Default(),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var (
		db  *gorm.DB
		err error
	)

	switch driver {
	case "postgres":
		for attempt := range MaxConnectionRetries {
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
			if err == nil {
				break
			}

			if attempt < MaxConnectionRetries-1 {
				log.Printf(
					"failed to connect to database (attempt %d/%d): %v;"+
						" retrying in %s",
					attempt+1, MaxConnectionRetries, err,
					ConnectionRetryDelay)
				time.Sleep(ConnectionRetryDelay)
			}
		}
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: newLogger})
	default:
		return nil, fmt.Errorf("unknown driver: %s", driver)
	}

	if err != nil {
		return nil, err
	}

	return &GormDatabase{db: db}, nil
}

func (g *GormDatabase) DB() *gorm.DB {
	return g.db
}

func (g *GormDatabase) Migrate(models ...any) error {
	if len(models) == 0 {
		return fmt.Errorf("no models provided for migration")
	}

	return g.db.AutoMigrate(models...)
}

func (g *GormDatabase) Close() error {
	sqlDB, err := g.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
