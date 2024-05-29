package service

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

func blackwhiteImage(img image.Image, threshold uint8) image.Image {
	// Create a new grayscale image
	bounds := img.Bounds()
	bwImg := image.NewGray(bounds)

	// Iterate through each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get the color of the current pixel
			originalColor := img.At(x, y)

			// Convert the color to grayscale
			grayColor := color.GrayModel.Convert(originalColor).(color.Gray)

			// Set the pixel to black or white based on the threshold
			if grayColor.Y > threshold {
				bwImg.Set(x, y, color.White)
			} else {
				bwImg.Set(x, y, color.Black)
			}
		}
	}
	return bwImg
}

func saturatedTempImage(halftone string, saturation int) (string, error) {
	srcImg, err := loadAssetImage(halftone)
	if err != nil {
		return "", err
	}

	// g := gift.New(gift.Saturation(float32(saturation)))
	// dstImg := image.NewNRGBA(g.Bounds(srcImg.Bounds()))
	// g.Draw(dstImg, srcImg)

	dstImg := blackwhiteImage(srcImg, uint8(math.Max(0, math.Min(255.0, -float64(saturation)))))

	tempFile, err := os.CreateTemp("", "tempfile-*.png")
	if err != nil {
		return "", err
	}
	if err = png.Encode(tempFile, dstImg); err != nil {
		return tempFile.Name(), err
	}
	tempFile.Sync()
	return tempFile.Name(), nil
}

func overlay(logo string, qrBytes []byte) ([]byte, error) {
	qrImage, err := png.Decode(bytes.NewReader(qrBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode QR code image: %v", err)
	}

	logoImage, err := loadAssetImage(logo)
	if err != nil {
		return nil, fmt.Errorf("failed to load logo image: %v", err)
	}

	logoImage = scaleLogoImage(logoImage, qrImage.Bounds().Dx(), qrImage.Bounds().Dy(), 0, 5)

	finalImage := drawLogo(logoImage, qrImage)
	finalBuf := new(bytes.Buffer)
	if err = png.Encode(finalBuf, finalImage); err != nil {
		return nil, fmt.Errorf("failed to encode final image: %v", err)
	}

	return finalBuf.Bytes(), nil
}

func drawLogo(logoImage image.Image, qrImage image.Image) image.Image {
	offset := image.Pt((qrImage.Bounds().Dx()-logoImage.Bounds().Dx())/2, (qrImage.Bounds().Dy()-logoImage.Bounds().Dy())/2)
	finalImage := image.NewRGBA(qrImage.Bounds())
	draw.Draw(finalImage, qrImage.Bounds(), qrImage, image.Point{}, draw.Src)
	draw.Draw(finalImage, logoImage.Bounds().Add(offset), logoImage, image.Point{}, draw.Over)
	return finalImage
}
