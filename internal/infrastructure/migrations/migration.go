package migrations

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Migrator interface {
	Run()
}

type Migration struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewMigration(db *sqlx.DB, logger *zap.Logger) Migrator {
	return &Migration{
		db:     db,
		logger: logger,
	}
}

func (m *Migration) Run() {
	driver, err := postgres.WithInstance(m.db.DB, &postgres.Config{})
	if err != nil {
		m.logger.Error("Failed to create postgres driver", zap.String("error", err.Error()))
		return
	}

	migrateInstance, err := migrate.NewWithDatabaseInstance(
		"file:///app/internal/infrastructure/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		m.logger.Error("Failed to create migrate instance", zap.String("error", err.Error()))
		return
	}

	if err := migrateInstance.Up(); err != nil && err != migrate.ErrNoChange {
		m.logger.Error("Failed to apply migrations", zap.String("error", err.Error()))
		return
	}

	m.logger.Info("Migrations applied successfully")
}
