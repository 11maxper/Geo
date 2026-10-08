package cache

import (
	"context"
	"encoding/json"
	"errors"
	"project/internal/models"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
)

var (
	cacheOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cache_operation_duration_seconds",
			Help:    "cache request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func observeCacheDuration(method string, start time.Time) {
	cacheOperationDuration.
		WithLabelValues(method).
		Observe(time.Since(start).Seconds())
}

type Cacher interface {
	SetAddress(ctx context.Context, key string, addresses models.ResponseAddress) error
	GetAddress(ctx context.Context, key string) (models.ResponseAddress, bool, error)
}

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) Cacher {
	return &Cache{
		client: client,
	}
}

func (c *Cache) SetAddress(ctx context.Context, key string, addresses models.ResponseAddress) error {
	//Время исполнения для метрик
	start := time.Now()
	defer observeCacheDuration("SetAddress", start)

	adrrBytes, err := json.Marshal(addresses)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, adrrBytes, 10*time.Minute).Err()
}

func (c *Cache) GetAddress(ctx context.Context, key string) (models.ResponseAddress, bool, error) {
	//Время исполнения для метрик
	start := time.Now()
	defer observeCacheDuration("GetAddress", start)

	addrByte, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return models.ResponseAddress{}, false, nil
		}
		return models.ResponseAddress{}, false, err
	}

	var Addresses models.ResponseAddress

	err = json.Unmarshal(addrByte, &Addresses)
	if err != nil {
		return models.ResponseAddress{}, false, err
	}

	return Addresses, true, nil
}
