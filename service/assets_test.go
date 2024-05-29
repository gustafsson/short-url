package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAssetImage(t *testing.T) {
	exampleUrl := "https://upload.wikimedia.org/wikipedia/commons/7/70/Example.png"
	img, err := loadAssetImage(exampleUrl)
	assert.NoError(t, err, "Loading exampleUrl should not produce an error")
	assert.NotNil(t, img, "Image should not be nil")

	exampleName := "example.png"
	img, err = loadAssetImage(exampleName)
	assert.NoError(t, err, "Loading exampleName should not produce an error")
	assert.NotNil(t, img, "Image should not be nil")
}
