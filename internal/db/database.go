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

	Queries() *sqlc.Queries
	SQLDB() *sql.DB
	WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error
}

type database struct {
	gormDB  *gorm.DB
	sqlDB   *sql.DB
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
		sqlDB   *sql.DB
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

		sqlDB = stdlib.OpenDBFromPool(pool)

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
		sqlDB:   sqlDB,
		gormDB:  db,
		pool:    pool,
		queries: queries,
	}, nil
}

func (d *database) DB() *gorm.DB {
	return d.gormDB
}

func (d *database) SQLDB() *sql.DB {
	return d.sqlDB
}

func (d *database) Queries() *sqlc.Queries {
	return d.queries
}

func (d *database) Migrate(models ...any) error {
	if len(models) == 0 {
		return fmt.Errorf("no models provided for migration")
	}

	return d.gormDB.AutoMigrate(models...)
}

func (d *database) Close() error {
	if d.pool != nil {
		defer d.pool.Close()
	}

	sqlDB, err := d.gormDB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func (d *database) WithTx(ctx context.Context,
	fn func(q *sqlc.Queries) error) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(d.queries.WithTx(tx)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
