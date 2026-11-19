package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/db/sqlc"
	"github.com/glebarez/sqlite"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	MaxConnectionRetries = 10
	ConnectionRetryDelay = 2 * time.Second
)

type Database interface {
	DB() *gorm.DB
	Migrate(models ...any) error
	Close() error
}

type database struct {
	gormDB  *gorm.DB
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewDatabase(ctx context.Context, driver, dsn string) (Database, error) {
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
		db      *gorm.DB
		pool    *pgxpool.Pool
		queries *sqlc.Queries
		err     error
	)

	switch driver {
	case "postgres":
		for attempt := range MaxConnectionRetries {
			pool, err = pgxpool.New(ctx, dsn)
			if err == nil {
				err = pool.Ping(ctx)
			}
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
		if err != nil {
			return nil, err
		}

		sqlDB := stdlib.OpenDBFromPool(pool)

		db, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{Logger: newLogger})
		if err != nil {
			pool.Close()
			return nil, err
		}

		queries = sqlc.New(pool)

	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: newLogger})
	default:
		return nil, fmt.Errorf("unknown driver: %s", driver)
	}

	if err != nil {
		return nil, err
	}

	return &database{
		gormDB:  db,
		pool:    pool,
		queries: queries,
	}, nil
}

func (g *database) DB() *gorm.DB {
	return g.gormDB
}

func (g *database) Queries() *sqlc.Queries {
	return g.queries
}

func (g *database) SQLDB() (*sql.DB, error) {
	return g.gormDB.DB()
}

func (g *database) Migrate(models ...any) error {
	if len(models) == 0 {
		return fmt.Errorf("no models provided for migration")
	}

	return g.gormDB.AutoMigrate(models...)
}

func (g *database) Close() error {
	if g.pool != nil {
		defer g.pool.Close()
	}

	sqlDB, err := g.gormDB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
