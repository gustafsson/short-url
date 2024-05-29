//go:build test
// +build test

package service

import (
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundPadding(t *testing.T) {
	inImg, err := loadAssetImage("example.png")
	assert.NoError(t, err, "No error should occur when loading an example image")
	outImg := AddPadding(inImg, 50)
	borderColor := color.NRGBA{R: 127, G: 255, B: 0, A: 255,}
	borderEdgeColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255,}
	outImg = MaskRoundCorners(outImg, 50, 5.0, borderColor, borderEdgeColor)

	// Save the new image to file
	outFile, err := os.Create("roundPadding_test.png")
	if err != nil {
		assert.NoError(t, err, "No error should occur when padding")
	}
	defer outFile.Close()

	png.Encode(outFile, outImg)
}
