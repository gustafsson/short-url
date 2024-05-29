package service

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
)

// createCircleMask creates a circular mask with the given radius and center
func createCircleMask(imgWidth, imgHeight int, radius int) *image.Alpha {
	mask := image.NewAlpha(image.Rect(0, 0, imgWidth, imgHeight))
	for y := 0; y < radius; y++ {
		for x := 0; x < radius; x++ {
			dx := float64(x - radius)
			dy := float64(y - radius)
			distance := math.Sqrt(dx*dx + dy*dy)
			if distance > float64(radius) {
				mask.SetAlpha(x, y, color.Alpha{A: 255})
				mask.SetAlpha(imgWidth-x, imgHeight-y, color.Alpha{A: 255})
				mask.SetAlpha(imgWidth-x, y, color.Alpha{A: 255})
				mask.SetAlpha(x, imgHeight-y, color.Alpha{A: 255})
			}
		}
	}
	return mask
}

func blend8(x, y, a uint8) uint8 {
	z := uint32(x)*uint32(255-a) + uint32(y)*uint32(a)
	if z/256 >= 256 {
		log.Printf("z %v\n", z)
	}
	return uint8(z / 256)
}

func mul8(x, y uint8) uint8 {
	return uint8((uint16(x)*uint16(y) + 127) / 255)
}

func avg8(x, y, ax, ay uint8) uint8 {
	totalWeight := uint32(ax) + uint32(ay)
	if totalWeight == 0 {
		return uint8((uint32(x) + uint32(y) + 1) / 2)
	}
	weightedSum := uint32(x)*uint32(ax) + uint32(y)*uint32(ay)

	// Compute the weighted average
	return uint8((weightedSum + totalWeight/2) / totalWeight)
}

func unitary2uint8(x float64) uint8 {
	return uint8(255 * x)
}

func MaskRoundCorners(img image.Image, radius, border float64, borderColor, borderEdgeColor color.NRGBA) image.Image {
	// Create a circular mask with rounded corners
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	// mask := createCircleMask(w, h, radius)
	fw, fh := float64(w), float64(h)
	newImg := image.NewNRGBA(image.Rect(0, 0, w, h))
	draw.Draw(newImg, img.Bounds(), img, image.Point{}, draw.Over)

	// Apply the mask to the new image
	// for y := 0; y < h; y++ {
	// 	for x := 0; x < w; x++ {
	// 		alpha := mask.AlphaAt(x, y).A
	// 		if alpha == 255 {
	// 			newImg.SetNRGBA(x, y, color.NRGBA{255, 255, 255, 0})
	// 		}
	// 	}
	// }

	D := func(x, y float64) float64 {
		dx := math.Min(0, math.Min(x-radius, (fw-x)-radius))
		dy := math.Min(0, math.Min(y-radius, (fh-y)-radius))
		distance := math.Sqrt(dx*dx + dy*dy)
		return distance
	}

	C := func(distance float64) (float64, float64, float64) {
		v, as, ab := 0.0, 0.0, 1.0 // show src over background
		d := distance - (radius - border)
		if d > border {
			ab = 0.0 // show background
		} else if d >= -border {
			f := d / border
			v = math.Abs(f)
			as = 1 - math.Max(0, -f*f*f)
			ab = 1 - math.Max(0, f*f*f)
		}
		return v, as, ab
	}

	C4 := func(x, y float64) (color.NRGBA, uint8) {
		c1, a1, b1 := C(D(x+0.25, y+0.25))
		c2, a2, b2 := C(D(x+0.75, y+0.25))
		c3, a3, b3 := C(D(x+0.25, y+0.75))
		c4, a4, b4 := C(D(x+0.75, y+0.75))

		c := (c1 + c2 + c3 + c4) / 4
		a := (a1 + a2 + a3 + a4) / 4
		b := (b1 + b2 + b3 + b4) / 4

		nrgba := color.NRGBA{
			uint8(float64(borderColor.R)*(1.0-c) + float64(borderEdgeColor.R)*c),
			uint8(float64(borderColor.G)*(1.0-c) + float64(borderEdgeColor.R)*c),
			uint8(float64(borderColor.B)*(1.0-c) + float64(borderEdgeColor.R)*c),
			// uint8(float64(borderColor.A) * (1.0 - a) + float64(borderEdgeColor.A) * a),
			unitary2uint8(a),
		}

		return nrgba, unitary2uint8(b)
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			t, b := C4(float64(x), float64(y))
			if t.A > 0 || b < 255 || true {
				s := newImg.NRGBAAt(x, y)
				a := uint8((uint32(s.A)*uint32(b) + 127) / 256)
				c2 := color.NRGBA{
					blend8(s.R, t.R, t.A),
					blend8(s.G, t.G, t.A),
					blend8(s.B, t.B, t.A),
					// mul8(blend8(s.R, t.R, t.A), a),
					// mul8(blend8(s.G, t.G, t.A), a),
					// mul8(blend8(s.B, t.B, t.A), a),
					a,
				}
				// c2 = color.NRGBA{
				// 	b.A, b.A, b.A, 255,
				// }
				newImg.Set(x, y, c2)
				// newImg.SetNRGBA(x, y, c2)
			}
		}
	}

	return newImg
}

func HasAlpha(img image.Image) bool {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a != 0xffff { // 0xffff means fully opaque in 16-bit color depth
				return true
			}
		}
	}
	return false
}

func AddPadding(img image.Image, padding int) image.Image {
	// log.Printf("Padding %v\n", padding)

	// Calculate new dimensions
	newWidth := img.Bounds().Dx() + 2*padding
	newHeight := img.Bounds().Dy() + 2*padding

	// Create a new RGBA image with white background
	newImg := image.NewNRGBA(image.Rect(0, 0, newWidth, newHeight))
	white := color.NRGBA{255, 255, 255, 255}
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	// Draw the original image onto the new image with padding
	draw.Draw(newImg, img.Bounds().Add(image.Point{X: padding, Y: padding}), img, image.Point{}, draw.Over)

	return newImg
}
