package DB

import (
	"fmt"

	"project/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func NewDB(cfg config.DB, logger *zap.Logger) *sqlx.DB {
	dsn := fmt.Sprintf("password=%s user=%s dbname=%s port=%s host=%s sslmode=disable",
		cfg.Password, cfg.User, cfg.Name, cfg.Port, cfg.Host)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Fatal("Failed to connect to postgres", zap.String("error", err.Error()))
	}

	return db
}
