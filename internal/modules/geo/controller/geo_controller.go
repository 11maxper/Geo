package controller

import (
	"encoding/json"
	"net/http"

	"project/internal/models"
	"project/internal/modules/geo/service"
)

type GeoProvider interface {
	AddressSearch(w http.ResponseWriter, r *http.Request)
	GeocodeSearch(w http.ResponseWriter, r *http.Request)
}

type GeoController struct {
	service service.GeoProvider
}

func NewGeoController(servise service.GeoProvider) GeoProvider {
	return &GeoController{
		service: servise,
	}
}

// AdressesSearch выполняет поиск адресов по названию города и улицы.
//
// @Summary Поиск адресов
// @Description Выполняет поиск адресов по введённой строке, возвращает список найденных адресов с координатами.
// @Tags address
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RequestAddressSearch true "Запрос с текстом для поиска адреса"
// @Success 200 {object} models.ResponseAddress "Успешный ответ с найденными адресами"
// @Failure 400 "Некорректный запрос"
// @Failure 403 "Недопустимый токен"
// @Failure 500 "Ошибка сервера"
// @Router /address/search [post]
func (c *GeoController) AddressSearch(w http.ResponseWriter, r *http.Request) {
	var req models.RequestAddressSearch
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	addresses, err := c.service.AddressSearch(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content type", "application/json")
	json.NewEncoder(w).Encode(addresses.Addresses)
}

// GeocodeSearch выполняет геокодирование по координатам.
//
// @Summary Геокодирование координат
// @Description Принимает координаты (широту и долготу) и возвращает информацию о соответствующем адресе.
// @Tags address
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RequestAddressGeocode true "Координаты для геокодирования"
// @Success 200 {object} models.ResponseAddress "Успешный ответ с данными об адресе"
// @Failure 400 "Некорректный запрос"
// @Failure 403 "Недопустимый токен"
// @Failure 500 "Ошибка сервера"
// @Router /address/geocode [post]
func (c *GeoController) GeocodeSearch(w http.ResponseWriter, r *http.Request) {
	var req models.RequestAddressGeocode
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	addresses, err := c.service.GeocodeSearch(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content type", "application/json")
	json.NewEncoder(w).Encode(addresses.Addresses)
}
