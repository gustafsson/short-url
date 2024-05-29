//go:build test
// +build test

package service

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaturatedTempImage(t *testing.T) {
	fn, err := saturatedTempImage("example.png", -90)
	if fn != "" {
		defer os.Remove(fn)
	}
	assert.NoError(t, err, "Producing a saturated image should not produce an error")

	file, err := os.Open(fn)
	assert.NoError(t, err, "Opening a tmp file should not produce an error")
	bytes, err := io.ReadAll(file)
	assert.NoError(t, err, "Reading a tmp file should not produce an error")

	err = saveBytes("saturation_test.png", bytes)
	assert.NoError(t, err, "Saving a saturated image should not produce an error")
}
