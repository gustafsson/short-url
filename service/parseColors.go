package service

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"
)

func parseColor(s string) (color.NRGBA, error) {
    // Trim spaces and convert to lowercase
    s = strings.TrimSpace(strings.ToLower(s))

    // Try to parse as a named color
    if c, ok := colornames.Map[s]; ok {
        return color.NRGBAModel.Convert(c).(color.NRGBA), nil
    }

    // Try to parse as a hex color
    if c, err := colorful.Hex("#" + s); err == nil {
        return colorfulToNRGBA(c), nil
    }

    // Try to parse as an RGB color
    var r, g, b int
    if _, err := fmt.Sscanf(s, "rgb(%d,%d,%d)", &r, &g, &b); err == nil {
        return color.NRGBA{uint8(r), uint8(g), uint8(b), 255}, nil
    }

    // Try to parse as an RGBA color
    var a int
    if _, err := fmt.Sscanf(s, "rgba(%d,%d,%d,%d)", &r, &g, &b, &a); err == nil {
        return color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}, nil
    }

    // If all parsing attempts fail
    return color.NRGBA{}, fmt.Errorf("unable to parse color: %s", s)
}

// Convert colorful.Color to color.NRGBA
func colorfulToNRGBA(c colorful.Color) color.NRGBA {
    r, g, b := c.RGB255()
    return color.NRGBA{R: r, G: g, B: b, A: 255}
}
