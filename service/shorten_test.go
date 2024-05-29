//go:build test
// +build test

package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"salbo.ai/short-url/repository"
)



func TestGenerateShortID(t *testing.T) {
	length := 5
	id := generateShortID(length)
	assert.Equal(t, length, len(id), "ID length should be equal to the specified length")
}

func TestFindNewID(t *testing.T) {
	repository.MockSetup()
	defer repository.MockTeardown()

	id, err := findNewID()
	assert.NoError(t, err, "No error should occur when generating a new ID")
	assert.NotEmpty(t, id, "Generated ID should not be empty")
}

func TestShortenURL(t *testing.T) {
	repository.MockSetup()
	defer repository.MockTeardown()
	short_prefix := os.Getenv("SHORT_PREFIX")

	longUrl := "https://example.com"
	shortUrl, err := ShortenURL(longUrl)
	assert.NoError(t, err, "No error should occur when shortening URL")
	assert.Contains(t, shortUrl, short_prefix, "Short URL should contain the base domain")

	fetchedLongUrl, data, err := GetRedirect(shortUrl)
	assert.NoError(t, err, "No error should occur when getting URL")
	assert.Equal(t, 0, len(data), "Expected no data")
	assert.Equal(t, longUrl, fetchedLongUrl, "Expected to fetch long url")
}
