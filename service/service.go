package service

import (
	"os"

	"salbo.ai/short-url/repository"
)

var short_prefix = os.Getenv("SHORT_PREFIX")

func ShortUrl(id string) string {
	return short_prefix + id
}

func GetRedirect(id string) (string, []byte, error) {
	return repository.GetRedirect(id)
}

func SaveRequest(data map[string]interface{}) error {
	return repository.SaveRequest(data)
}
