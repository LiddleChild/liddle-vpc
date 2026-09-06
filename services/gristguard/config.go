package main

import (
	"net/url"
)

type Config struct {
	Port int `env:"PORT"`

	GristEndpoint   *url.URL `env:"GRIST_ENDPOINT"`
	GristAPIKey     string   `env:"GRIST_API_KEY"`
	GristDocumentID string   `env:"GRIST_DOCUMENT_ID"`
}
