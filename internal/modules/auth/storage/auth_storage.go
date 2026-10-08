package storage

import (
	"project/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	dbOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_operation_duration_seconds",
			Help:    "db request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func observeDBDuration(method string, start time.Time) {
	dbOperationDuration.
		WithLabelValues(method).
		Observe(time.Since(start).Seconds())
}

type Auther interface {
	FindByName(name string) (models.User, error)
	Create(name string, hash []byte) error
}

type AuthStorage struct {
	DB *sqlx.DB
}

func NewAuthStorage(db *sqlx.DB) Auther {
	return &AuthStorage{
		DB: db,
	}
}

func (s *AuthStorage) FindByName(name string) (models.User, error) {
	//Время обращения для метрик
	start := time.Now()
	defer observeDBDuration("FindByName", start)

	query := `
	SELECT id, name, password
	FROM users
	WHERE name = $1
	`
	var user models.User
	err := s.DB.Get(&user, query, name)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *AuthStorage) Create(name string, hash []byte) error {
	//Время обращения для метрик
	start := time.Now()
	defer observeDBDuration("Create", start)

	query := `
	INSERT INTO users 
	(name, password)
	VALUES ($1, $2)
	`
	_, err := s.DB.Exec(query, name, hash)
	if err != nil {
		return err
	}
	return nil
}
