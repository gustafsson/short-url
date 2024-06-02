//go:build test
// +build test

package repository

import (
	"sync"
)

var (
	mockDB    = make(map[string]string)
	mockMutex sync.Mutex
)

func MockSetup() {
	mockMutex.Lock()
	mockDB = make(map[string]string)
	mockMutex.Unlock()
}

func MockTeardown() {
	mockMutex.Lock()
	mockDB = nil
	mockMutex.Unlock()
}

func CheckIfExists(id string) (bool, error) {
	mockMutex.Lock()
	defer mockMutex.Unlock()
	_, exists := mockDB[id]
	return exists, nil
}

func SaveURL(id, longURL string) error {
	mockMutex.Lock()
	defer mockMutex.Unlock()
	mockDB[id] = longURL
	return nil
}

func SaveRequest(data map[string]interface{}) error {
	return nil
}

func GetRedirect(id string) (string, []byte, error) {
	mockMutex.Lock()
	defer mockMutex.Unlock()
	url, exists := mockDB[id]
	if !exists {
		return "", nil, nil
	}
	return url, nil, nil
}
