package service

import (
	"context"
	"fmt"
	"project/internal/infrastructure/cache"
	"project/internal/models"
	"project/internal/modules/geo/storage"

	"go.uber.org/zap"
)

type GeoProvider interface {
	AddressSearch(ctx context.Context, input models.RequestAddressSearch) (models.ResponseAddress, error)
	GeocodeSearch(ctx context.Context, input models.RequestAddressGeocode) (models.ResponseAddress, error)
}

type GeoService struct {
	storage storage.GeoProvider
	cache   cache.Cacher
	logger  *zap.Logger
}

func NewGeoservice(storage storage.GeoProvider, cache cache.Cacher, logger *zap.Logger) GeoProvider {
	return &GeoService{
		storage: storage,
		cache:   cache,
		logger:  logger,
	}
}

// AddressSearch поиск адреса по названию города и улицы
func (s *GeoService) AddressSearch(ctx context.Context, input models.RequestAddressSearch) (models.ResponseAddress, error) {
	key := fmt.Sprintf("address:%s", input.Query)

	// пытаемся получить данные из кеша
	if cached, exists, err := s.cache.GetAddress(ctx, key); err == nil && exists {
		s.logger.Info("в ответе использовался кеш")
		return cached, nil
	} else if err != nil {
		s.logger.Info("cache.GetAddress()", zap.Error(err))
	}

	// если в кеше нет данных, получаем из api
	addresses, err := s.storage.AddressSearch(input)
	if err != nil {
		s.logger.Info("storage.AddressSearch()", zap.Error(err))
		return models.ResponseAddress{}, err
	}

	// сохраняем данные в кеше
	if err := s.cache.SetAddress(ctx, key, addresses); err != nil {
		s.logger.Info("cache.SetAddress()", zap.Error(err))
	}

	return addresses, nil
}

// GeocodeSearch поиск адреса по координатам
func (s *GeoService) GeocodeSearch(ctx context.Context, input models.RequestAddressGeocode) (models.ResponseAddress, error) {
	key := fmt.Sprintf("geocode:%s:%s", input.Lat, input.Lng)

	// пытаемся получить данные из кеша
	if cached, exists, err := s.cache.GetAddress(ctx, key); err == nil && exists {
		s.logger.Info("в ответе использовался кеш")
		return cached, nil
	} else if err != nil {
		s.logger.Info("cache.GetAddress()", zap.Error(err))
	}

	// если в кеше нет данных, получаем из api
	addresses, err := s.storage.GeoCode(input)
	if err != nil {
		s.logger.Info("storage.AddressSearch()", zap.Error(err))
		return models.ResponseAddress{}, err
	}

	// сохраняем данные в кеше
	if err := s.cache.SetAddress(ctx, key, addresses); err != nil {
		s.logger.Info("cache.SetAddress()", zap.Error(err))
	}

	return addresses, nil
}
