package models

// RequestAddressSearch запрос для поиска адреса по названию города и улицы
type RequestAddressSearch struct {
	Query string `json:"query" example:"Москва, Тверская 7"`
}

// ResponseAddress ответ на поиск адреса
type ResponseAddress struct {
	Addresses []*Address `json:"addresses"`
}

type Address struct {
	City   string `json:"city" example:"Москва"`
	Street string `json:"street" example:"Тверская"`
	House  string `json:"house" example:"7"`
	Lat    string `json:"lat" example:"55.7579795"`
	Lon    string `json:"lon" example:"37.611263"`
}

// RequestAddressGeocode запрос для поиска по координатам
type RequestAddressGeocode struct {
	Lat string `json:"lat" example:"55.7579795"`
	Lng string `json:"lng" example:"37.611263"`
}
