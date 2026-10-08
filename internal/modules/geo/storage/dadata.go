package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"project/config"
	"project/internal/models"

	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	apiOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_operation_duration_seconds",
			Help:    "api request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func observeAPIDuration(method string, start time.Time) {
	apiOperationDuration.
		WithLabelValues(method).
		Observe(time.Since(start).Seconds())
}

type GeoProvider interface {
	AddressSearch(input models.RequestAddressSearch) (models.ResponseAddress, error)
	GeoCode(input models.RequestAddressGeocode) (models.ResponseAddress, error)
}

type GeoStorage struct {
	api    *suggest.Api
	config *config.Dadata
}

func NewStorage(api *suggest.Api, cfg *config.Dadata) GeoProvider {
	return &GeoStorage{
		api:    api,
		config: cfg,
	}
}

func (s *GeoStorage) AddressSearch(input models.RequestAddressSearch) (models.ResponseAddress, error) {
	// Время обращения для метрик
	start := time.Now()
	defer observeAPIDuration("AddressSearch", start)

	var res models.ResponseAddress
	rawRes, err := s.api.Address(context.Background(), &suggest.RequestParams{Query: input.Query})
	if err != nil {
		return models.ResponseAddress{}, err
	}

	for _, r := range rawRes {
		if r.Data.City == "" || r.Data.Street == "" {
			continue
		}
		res.Addresses = append(res.Addresses, &models.Address{City: r.Data.City, Street: r.Data.Street, House: r.Data.House, Lat: r.Data.GeoLat, Lon: r.Data.GeoLon})
	}

	return res, nil
}

func (s *GeoStorage) GeoCode(input models.RequestAddressGeocode) (models.ResponseAddress, error) {
	// Время обращения для метрик
	start := time.Now()
	defer observeAPIDuration("GeoCode", start)

	httpClient := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`{"lat": %s, "lon": %s}`, input.Lat, input.Lng))
	req, err := http.NewRequest("POST", "https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address", data)
	if err != nil {
		return models.ResponseAddress{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", s.config.Key))
	resp, err := httpClient.Do(req)
	if err != nil {
		return models.ResponseAddress{}, err
	}
	var geoCode GeoCode

	err = json.NewDecoder(resp.Body).Decode(&geoCode)
	if err != nil {
		return models.ResponseAddress{}, err
	}
	var res models.ResponseAddress
	for _, r := range geoCode.Suggestions {
		var address models.Address
		address.City = string(r.Data.City)
		address.Street = string(r.Data.Street)
		address.House = r.Data.House
		address.Lat = r.Data.GeoLat
		address.Lon = r.Data.GeoLon

		res.Addresses = append(res.Addresses, &address)
	}

	return res, nil
}
