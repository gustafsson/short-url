package service

import (
	"fmt"
	"math/rand"

	"salbo.ai/short-url/repository"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func init() {
	// rand.Seed(time.Now().UnixNano())
}

func ShortenURL(longURL string) (string, error) {
	id, err := findNewID()
	if err != nil {
		return "", err
	}

	if err := repository.SaveURL(id, longURL); err != nil {
		return "", err
	}

	return ShortUrl(id), nil
}

func findNewID() (string, error) {
	maxAttempts := 10  // Maximum number of attempts to find a unique ID
	maxLength := 10    // Maximum length of the ID
	initialLength := 1 // Start with an ID of length 1

	for length := initialLength; length <= maxLength; length++ {
		for attempt := 1; attempt <= maxAttempts; attempt++ {
			id := generateShortID(length)
			exists, err := repository.CheckIfExists(id)
			if err != nil {
				return "", err
			}
			if !exists {
				return id, nil
			}
		}
	}
	return "", fmt.Errorf("could not generate a unique ID after %d attempts", maxAttempts*(maxLength-initialLength+1))
}

func generateShortID(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
