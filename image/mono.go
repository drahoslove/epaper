package image

import (
	"bytes"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
)

var (
	goFont *truetype.Font
)

func init() {
	goFont, _ = truetype.Parse(gobold.TTF)
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
	if r == 0 {
		m[4+i/8] &^= (1 << (7 - i%8)) // clr
	} else {
		m[4+i/8] |= (1 << (7 - i%8)) // set
	}
}

func (m Mono) At(x, y int) color.Color {
	i := uint(x) + uint(y)*m.Width()
	r := m[4+i/8] & (1 << (7 - i%8))
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

func (m *Mono) Clear(c color.Color) {
	Y, _, _, _ := c.RGBA()
	*m = append([]byte(*m)[:4], bytes.Repeat([]byte{byte(Y)}, len(*m)-4)...)
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

func (m *Mono) DrawString(text string, size float64, x, y int) {
	d := font.Drawer{
		Dst: m,
		Src: image.Black,
		Face: truetype.NewFace(goFont, &truetype.Options{
			Size:    size,
			Hinting: font.HintingFull,
		}),
	}
	d.Dot = fixed.P(x, y)
	d.DrawString(text)
}

func flip(b byte) byte {
	b = (b&0xF0)>>4 | (b&0x0F)<<4
	b = (b&0xCC)>>2 | (b&0x33)<<2
	b = (b&0xAA)>>1 | (b&0x55)<<1
	return b
}

func (m *Mono) VerticalFlip() {
	data := m.Bitmap()
	w := m.Width() / 8
	h := m.Height()
	for y := uint(0); y < h/2; y++ {
		for i := uint(0); i < w; i++ {
			data[y*w+i], data[(h-1-y)*w+i] = data[(h-1-y)*w+i], data[y*w+i]
		}
	}
}

func (m *Mono) HorizontalFlip() {
	data := m.Bitmap()
	w := m.Width() / 8
	h := m.Height()
	for y := uint(0); y < h; y++ {
		for i := uint(0); i < w/2; i++ {
			data[y*w+i], data[y*w+w-1-i] = flip(data[y*w+w-1-i]), flip(data[y*w+i])
		}
	}
}

func (m *Mono) RotateRight() {
	w := m.Width()
	h := m.Height()

	n := NewMono(h, w)

	for y := 0; y < int(h); y++ {
		for x := 0; x < int(w); x++ {
			n.Set(int(h)-1-y, x, m.At(x, y))
		}
	}
	*m = n
}

func (m *Mono) RotateLeft() {
	w := m.Width()
	h := m.Height()

	n := NewMono(h, w)

	for y := 0; y < int(h); y++ {
		for x := 0; x < int(w); x++ {
			n.Set(y, int(w)-1-x, m.At(x, y))
		}
	}
	*m = n
}
