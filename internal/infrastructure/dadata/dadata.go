package dadata

import (
	"net/url"
	"project/config"

	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
)

func NewDadata(c config.Dadata) *suggest.Api {
	var err error
	endpointUrl, err := url.Parse("https://suggestions.dadata.ru/suggestions/api/4_1/rs/")
	if err != nil {
		return nil
	}

	creds := client.Credentials{
		ApiKeyValue:    c.Key,
		SecretKeyValue: c.SecretKey,
	}

	return &suggest.Api{
		Client: client.NewClient(endpointUrl, client.WithCredentialProvider(&creds)),
	}
}
