package main

import (
	"net/url"
)

type Config struct {
	Port int `env:"PORT" envDefault:"8080"`

	GristEndpoint   *url.URL `env:"GRIST_ENDPOINT"`
	GristAPIKey     string   `env:"GRIST_API_KEY"`
	GristDocumentID string   `env:"GRIST_DOCUMENT_ID"`

	FrankfurterEndpoint *url.URL `env:"FRANKFURTER_ENDPOINT"`

	CurrencyTableID string `env:"BACKEND_CURRENCY_TABLE_ID"`
}
