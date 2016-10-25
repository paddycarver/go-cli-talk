package cmd

import (
	"errors"
	"net/http"
	"os"

	"github.com/paddyforan/go-cli-talk/proverbs"
)

func getProverb(id string) (proverbs.Quote, error) {
	baseURL := os.Getenv("PROVERBS_URL")
	if baseURL == "" {
		return proverbs.Quote{}, errors.New("PROVERBS_URL must be set to the API endpoint to retrieve proverbs from.")
	}

	headers := http.Header{}
	errMode := os.Getenv("ERROR")
	switch errMode {
	case "400":
		headers.Set("Return-Error", "bad-request")
	case "500":
		headers.Set("Return-Error", "internal")
	}

	return proverbs.GetProverb(baseURL, id, headers)
}
