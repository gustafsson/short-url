package service

import (
	"fmt"
	"image"
	"math"
	"os"

	"github.com/nfnt/resize"
	qrcode "github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type QRCodeOptions struct {
	Logo               string  `json:"logo"`
	Halftone           string  `json:"halftone"`
	Saturation         int     `json:"saturation"`
	Padding            float64 `json:"padding"`
	Border             float64 `json:"border"`
	BorderColor        string  `json:"border_color"`
	BorderEdgeColor    string  `json:"border_edge_color"`
	Radius             float64 `json:"radius"`
	LogoSizeMultiplier float64 `json:"logo_size_multiplier"`
	Transparent        bool    `json:"transparent"`
	QrWidth            int     `json:"qrwidth"`
	CorrectionLevel    string  `json:"correctionlevel"`
}

// DefaultQRCodeOptions returns a QRCodeOptions struct with default values
func DefaultQRCodeOptions() QRCodeOptions {
	return QRCodeOptions{
		Logo:               "",
		Halftone:           "",
		Saturation:         -90,
		Padding:            0.03,
		Border:             0.01,
		BorderColor:        "black",
		BorderEdgeColor:    "white",
		Radius:             1.0 / 15.0,
		LogoSizeMultiplier: 4,
		Transparent:        false,
		QrWidth:            41,
		CorrectionLevel:    "medium",
	}
}

func scaleLogoImage(logoImage image.Image, width, height, padding int, logoSizeMultiplier float64) image.Image {
	// Rescale the logo to fit within 1/5 of the QR code's dimensions
	logoW, logoH := float64(logoImage.Bounds().Dx()), float64(logoImage.Bounds().Dy())
	m := logoSizeMultiplier
	scale := math.Min(
		(float64(width)-float64(padding)*m*2)/float64(m)/logoW,
		(float64(height)-float64(padding)*m*2)/float64(m)/logoH)
	return resize.Resize(uint(scale*logoW), uint(scale*logoH), logoImage, resize.Lanczos3)
}

func GenerateQRCode(text string, opt QRCodeOptions) ([]byte, error) {
	ecLevel := qrcode.ErrorCorrectionMedium
	switch opt.CorrectionLevel {
	case "low":
		ecLevel = qrcode.ErrorCorrectionLow
	case "medium":
		ecLevel = qrcode.ErrorCorrectionMedium
	case "high":
		ecLevel = qrcode.ErrorCorrectionHighest
	case "quartile":
		ecLevel = qrcode.ErrorCorrectionQuart
	default:
		return nil, fmt.Errorf("expected one of 'low', 'medium', 'high', 'quartile', but got '%s'", opt.CorrectionLevel)
	}

	qrc, err := qrcode.NewWith(text, qrcode.WithErrorCorrectionLevel(ecLevel))
	if err != nil {
		return nil, err
	}

	options := []standard.ImageOption{
		standard.WithQRWidth(uint8(opt.QrWidth)),
		standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
	}

	if opt.Transparent {
		options = append(options,
			standard.WithBgTransparent(),
			standard.WithBgColorRGBHex("#000000"),
		)
	}

	if opt.Logo != "" {
		logoImage, err := loadAssetImage(opt.Logo)
		if err != nil {
			return nil, err
		}

		options = append(
			options, standard.WithLogoSizeMultiplier(int(opt.LogoSizeMultiplier)))

		expectedWidth := (qrc.Dimension() + 1) * opt.QrWidth
		padding_px := 0
		if opt.Border > 0 {
			padding_px = int(float64(expectedWidth) * opt.Padding)
		}
		logoImage = scaleLogoImage(logoImage, expectedWidth, expectedWidth, padding_px, opt.LogoSizeMultiplier)
		if opt.Border > 0 || opt.Radius != 0 {
			if padding_px != 0 {
				logoImage = AddPadding(logoImage, padding_px)
			}
			borderColor, err := parseColor(opt.BorderColor)
			if err != nil {
				return nil, err
			}
			borderEdgeColor, err := parseColor(opt.BorderEdgeColor)
			if err != nil {
				return nil, err
			}

			logoImage = MaskRoundCorners(logoImage,
				float64(expectedWidth)*opt.Radius,
				float64(expectedWidth)*opt.Border,
				borderColor, borderEdgeColor,
			)
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
