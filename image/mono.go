package image

import (
	"bytes"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
)

var (
	goFont *truetype.Font
)

func init() {
	goFont, _ = truetype.Parse(gomono.TTF)
}

// Mono implements image.Image and image/draw.Image
type Mono []byte

func NewMono(width, height uint) Mono {
	bitmap := make([]byte, width*height/8+4)
	bitmap[0], bitmap[1] = byte(width>>8), byte(width)
	bitmap[2], bitmap[3] = byte(height>>8), byte(height)
	return bitmap
}

func (m Mono) Bitmap() []byte {
	return m[4:]
}

func (m Mono) Width() uint {
	return uint(m[0])<<8 + uint(m[1])
}

func (m Mono) Height() uint {
	return uint(m[2])<<8 + uint(m[3])
}

func (m Mono) Set(x, y int, c color.Color) {
	if x < 0 || y < 0 || x >= int(m.Width()) || y >= int(m.Height()) {
		return
	}
	i := uint(x) + uint(y)*m.Width()
	r, _, _, _ := c.RGBA()
	m[4:][i/8] |= byte(r) & (1 << (7 - i%8)) // i >> 3 means Math.floor(i/(8))
}

func (m Mono) At(x, y int) color.Color {
	i := uint(x) + uint(y)*m.Width()
	r := m[4:][i/8] & (1 << (7 - i%8))
	if r == 0 {
		return color.Black
	} else {
		return color.White
	}
}

func (m *Mono) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(m.Width()), int(m.Height()))
}

func (m *Mono) ColorModel() color.Model {
	return color.ModelFunc(func(c color.Color) color.Color {
		r, g, b, _ := c.RGBA()
		avg := (r + g + b) / 3
		if avg > 1<<15 { // basic tresholding
			return color.Black
		} else {
			return color.White
		}
	})
}

func (m *Mono) Clear() {
	*m = append([]byte(*m)[4:], bytes.Repeat([]byte{0}, len(*m)-4)...)
}

func (m *Mono) DrawHorizontalLine(y, x_start, x_end int) {
	for x := x_start; x < x_end; x++ {
		m.Set(x, y, color.Black)
	}
}

func (m *Mono) DrawVerticalLine(x, y_start, y_end int) {
	for y := y_start; y < y_end; y++ {
		m.Set(x, y, color.Black)
	}
}

func (m *Mono) StrokeRect(x_start, y_start, x_end, y_end int) {
	m.DrawHorizontalLine(y_start, x_start, x_end)
	m.DrawHorizontalLine(y_end, x_start, x_end)
	m.DrawVerticalLine(x_start, y_start, y_end)
	m.DrawVerticalLine(x_end, y_start, y_end)
}

func (m *Mono) FillRect(x_start, y_start, x_end, y_end int) {
	for y := y_start; y < y_end; y++ {
		m.DrawHorizontalLine(y, x_start, x_end)
	}
}

func (m *Mono) DrawString(s string, size float64, x, y int) {
	d := font.Drawer{
		Dst: m,
		Src: image.Black,
		Face: truetype.NewFace(goFont, &truetype.Options{
			Size:    size,
			DPI:     72,
			Hinting: font.HintingNone,
		}),
	}
	d.Dot = fixed.P(x, y)
	d.DrawString(s)
}
