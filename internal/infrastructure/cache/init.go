package cache

import (
	"context"
	"project/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewCacheClient(ctx context.Context, cfg config.Cache, logger *zap.Logger) *redis.Client {
	opt := &redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	}

	client := redis.NewClient(opt)

	//проверка соединения
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Fatal("Ошибка соединения с Redis:", zap.Error(err))
	}
	logger.Info("Соединение с Redis успешно:", zap.String("Ping", pong))

	return client
}
