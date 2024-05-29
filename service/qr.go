package service

import (
	"image"
	"math"
	"os"

	"github.com/nfnt/resize"
	qrcode "github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type QRCodeOptions struct {
	Logo              string  `json:"logo"`
	Halftone          string  `json:"halftone"`
	Saturation        int     `json:"saturation"`
	Padding           float64 `json:"padding"`
	Border            float64 `json:"border"`
	Radius            float64 `json:"radius"`
	LogoSizeMultiplier int    `json:"logo_size_multiplier"`
}

// DefaultQRCodeOptions returns a QRCodeOptions struct with default values
func DefaultQRCodeOptions() QRCodeOptions {
	return QRCodeOptions{
		Logo:              "",
		Halftone:          "",
		Saturation:        -90,
		Padding:           0.03,
		Border:            0.01,
		Radius:            1.0/15.0,
		LogoSizeMultiplier: 4,
	}
}

func scaleLogoImage(logoImage image.Image, width, height, padding, logoSizeMultiplier int) image.Image {
	// Rescale the logo to fit within 1/5 of the QR code's dimensions
	logoW, logoH := float64(logoImage.Bounds().Dx()), float64(logoImage.Bounds().Dy())
	m := logoSizeMultiplier
	scale := math.Min(
		float64(width-padding*m*2)/float64(m)/logoW,
		float64(height-padding*m*2)/float64(m)/logoH)
	return resize.Resize(uint(scale*logoW), uint(scale*logoH), logoImage, resize.Lanczos3)
}

func GenerateQRCode(text string, opt QRCodeOptions) ([]byte, error) {
	qrc, err := qrcode.New(text)
	if err != nil {
		return nil, err
	}

	qrwidth := 41
	options := []standard.ImageOption{
		standard.WithQRWidth(uint8(qrwidth)),
		standard.WithBgTransparent(), standard.WithBgColorRGBHex("#000000"),
		standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
	}

	if opt.Logo != "" {
		logoImage, err := loadAssetImage(opt.Logo)
		if err != nil {
			return nil, err
		}

		options = append(
			options, standard.WithLogoSizeMultiplier(opt.LogoSizeMultiplier))

		expectedWidth := (qrc.Dimension() + 2) * qrwidth
		padding_px := 0
		if !HasAlpha(logoImage) {
			padding_px = int(float64(expectedWidth) * opt.Padding)
		}
		logoImage = scaleLogoImage(logoImage, expectedWidth, expectedWidth, padding_px, opt.LogoSizeMultiplier)
		if !HasAlpha(logoImage) {
			if padding_px != 0 {
				logoImage = AddPadding(logoImage, padding_px)
			}
			logoImage = MaskRoundCorners(logoImage, float64(expectedWidth)*opt.Radius, float64(expectedWidth)*opt.Border)
		}

		options = append(
			options, standard.WithLogoImage(logoImage))
	}

	if opt.Halftone != "" {
		halftone_fn, err := saturatedTempImage(opt.Halftone, opt.Saturation)
		if halftone_fn != "" {
			defer os.Remove(halftone_fn)
		}
		if err != nil {
			return nil, err
		}

		options = append(
			options, standard.WithHalftone(halftone_fn))
	}

	buf := newBufferWriteCloser()
	pngWriter := standard.NewWithWriter(buf, options...)
	if err = qrc.Save(pngWriter); err != nil {
		return nil, err
	}

	qrBytes := buf.(*bufferWriteCloser).Bytes()

	return qrBytes, nil
}