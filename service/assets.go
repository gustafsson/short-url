package service

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	assetsDir    = os.Getenv("ASSETS_DIR")
)

func openAssetFile(logodir, logo string) (io.ReadCloser, error) {
	if strings.HasPrefix(logo, "http") {
		resp, err := http.Get(logo)
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	} else {
		logoFilename := fmt.Sprintf("%s/%s", logodir, logo)

		// cwd, err := os.Getwd()
		// if err != nil {
		// 	return nil, err
		// }

		// log.Printf("Opening %s with cwd=%s", logoFilename, cwd)
		return os.Open(logoFilename)
	}
}

func loadAssetImage(logo string) (image.Image, error) {
	assetFile, err := openAssetFile(assetsDir, logo)
	if err != nil {
		return nil, fmt.Errorf("failed to open logo file: %v", err)
	}
	defer assetFile.Close()

	image, _, err := image.Decode(assetFile)
	if false {
		jpeg.Decode(assetFile)
		png.Decode(assetFile)
	}
	return image, err
}

func checkFileType(body io.Reader) (string, error) {
	// Read the first 8 bytes
	buf := make([]byte, 8)
	_, err := io.ReadFull(body, buf)
	if err != nil {
		return "", fmt.Errorf("error reading file header: %v", err)
	}

	// Check for JPEG magic number
	if buf[0] == 0xFF && buf[1] == 0xD8 && buf[2] == 0xFF {
		return "JPEG", nil
	}

	// Check for PNG magic number
	if buf[0] == 0x89 && buf[1] == 0x50 && buf[2] == 0x4E && buf[3] == 0x47 {
		return "PNG", nil
	}

	return "", fmt.Errorf("did not recognise file header as PNG or JPEG")
}
