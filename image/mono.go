package image

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"math"
	"math/bits"
)

var (
	goFont     *truetype.Font
	colorModel color.Model
)

func init() {
	goFont, _ = truetype.Parse(gobold.TTF)
	colorModel = color.ModelFunc(func(c color.Color) color.Color {
		r, g, b, _ := c.RGBA()
		avg := (r + g + b) / 3
		if avg < 1<<15 { // tresholding
			return color.Black
		} else {
			return color.White
		}
	})
}

// Mono is monochromatic image
//
// It implements image.Image and image/draw.Image interface
type Mono []byte

func NewMono(rect image.Rectangle) Mono {
	width, height := rect.Size().X, rect.Size().Y
	h := height
	w := width / 8
	if width%8 > 0 {
		w++
	}
	bitmap := make([]byte, (w*h + 4))
	bitmap[0], bitmap[1] = byte(width>>8), byte(width)
	bitmap[2], bitmap[3] = byte(height>>8), byte(height)
	return bitmap
}

// Bitmap returns byte slice containing actual image data
func (m Mono) Bitmap() []byte {
	return m[4:]
}

// Width returns width of image
func (m Mono) Width() uint {
	return uint(m[0])<<8 | uint(m[1])
}

// Height return height of image
func (m Mono) Height() uint {
	return uint(m[2])<<8 | uint(m[3])
}

// Set sets color on given coordinates.
// Color should be either color.Black or color.White.
//
// Implements image/draw.Image interface.
func (m Mono) Set(x, y int, c color.Color) {
	if x < 0 || y < 0 || x >= int(m.Width()) || y >= int(m.Height()) {
		return
	}
	w := m.Width()
	if w%8 > 0 {
		w += 8 - w%8
	}
	i := uint(x) + uint(y)*w
	Y := colorModel.Convert(c).(color.Gray16).Y
	if Y == 0 {
		m[4+i/8] &^= (1 << (7 - i%8)) // clr
	} else {
		m[4+i/8] |= (1 << (7 - i%8)) // set
	}
}

// At returns color at giver coordinates.
// returned values are either color.Black or color.White.
//
// Implements image.Image iterface.
func (m Mono) At(x, y int) color.Color {
	w := m.Width()
	if w%8 > 0 {
		w += 8 - w%8
	}
	i := uint(x) + uint(y)*w
	bit := m[4+i/8] & (1 << (7 - i%8))
	if bit == 0 {
		return color.Black
	} else {
		return color.White
	}
}

// Bounds returns Rectangle bounding the image.
//
// Implements image.Image interface.
func (m Mono) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(m.Width()), int(m.Height()))
}

// ColorModel return color.Model of the image.
// Color converted to this model results either to color.Black or color.White.
// Basic fixed tresholding method is used.
//
// Implement image.Image interface.
func (m Mono) ColorModel() color.Model {
	return colorModel
}

// Clear sets whole bitmap to given color - color.Black or color.White
func (m *Mono) Clear(c color.Color) {
	r, g, b, _ := c.RGBA()
	Y := byte((r + g + b) / 3)
	for i := 4; i < len(*m); i++ {
		(*m)[i] = Y
	}
}

// DrawHorizontalLine draws horizontal line given by left most point and length
//
// lengtj is distance between centers of first and last dot:
// line of len 0 is dot -> line will consists of length+1 dots
func (m *Mono) DrawHorizontalLine(color color.Color, start image.Point, length int) {
	for x := start.X; x <= start.X+length; x++ {
		m.Set(x, start.Y, color)
	}
}

// DrawHorizontalLine draws vettical line given by top most point and length
//
// length is distance between centers of first and last dot:
// line of len 0 is dot -> line will consists of length+1 dots
func (m *Mono) DrawVerticalLine(color color.Color, start image.Point, length int) {
	for y := start.Y; y <= start.Y+length; y++ {
		m.Set(start.X, y, color)
	}
}

// Draw arbitrary line
//
// line with same start and end is 1 dot
func (m *Mono) DrawLine(color color.Color, start image.Point, end image.Point) {
	drawLineLow := func(start, end image.Point) {
		delta := end.Sub(start)
		yi := +1
		if delta.Y < 0 {
			delta.Y = -delta.Y
			yi = -1
		}
		D := 2*delta.Y - delta.X

		for x, y := start.X, start.Y; x <= start.X+delta.X; x++ {
			m.Set(x, y, color)
			if D > 0 {
				y += yi
				D -= 2 * delta.X
			}
			D += 2 * delta.Y
		}
	}
	drawLineHight := func(start, end image.Point) {
		delta := end.Sub(start)
		xi := +1
		if delta.X < 0 {
			delta.X = -delta.X
			xi = -1
		}
		D := 2*delta.X - delta.Y

		for y, x := start.Y, start.X; y <= start.Y+delta.Y; y++ {
			m.Set(x, y, color)
			if D > 0 {
				x += xi
				D -= 2 * delta.Y
			}
			D += 2 * delta.X
		}
	}

	if math.Abs(float64(end.Y-start.Y)) < math.Abs(float64(end.X-start.X)) {
		if start.X > end.X {
			drawLineLow(end, start)
		} else {
			drawLineLow(start, end)
		}
	} else {
		if start.Y > end.Y {
			drawLineHight(end, start)
		} else {
			drawLineHight(start, end)
		}
	}
}

// StrokeRect draws outline of rectangle
func (m *Mono) StrokeRect(color color.Color, rect image.Rectangle) {
	w, h := rect.Dx(), rect.Dy()
	m.DrawHorizontalLine(color, rect.Min, w)
	m.DrawHorizontalLine(color, rect.Min.Add(image.Pt(0, h)), w)
	m.DrawVerticalLine(color, rect.Min, h)
	m.DrawVerticalLine(color, rect.Min.Add(image.Pt(w, 0)), h)
}

// StrokeRect draws filled rectangle
func (m *Mono) FillRect(color color.Color, rect image.Rectangle) {
	w := rect.Dx()
	down := image.Pt(0, 1)
	for start := rect.Min; start.Y <= rect.Max.Y; start = start.Add(down) {
		m.DrawHorizontalLine(color, start, w)
	}
}

// StrokeCircle draws outline of circle given by center point and raidus.
//
// Center is the coords of pixel in center - circle with radius 3 will be 5 px wide.
func (m *Mono) StrokeCircle(color color.Color, center image.Point, radius int) {
	x := radius - 1
	y := 0
	dx := 1
	dy := 1
	err := dx - (radius << 1)

	for x >= y {
		m.Set(center.X+x, center.Y+y, color)
		m.Set(center.X+y, center.Y+x, color)
		m.Set(center.X-y, center.Y+x, color)
		m.Set(center.X-x, center.Y+y, color)
		m.Set(center.X-x, center.Y-y, color)
		m.Set(center.X-y, center.Y-x, color)
		m.Set(center.X+y, center.Y-x, color)
		m.Set(center.X+x, center.Y-y, color)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (radius << 1)
		}
	}
}

// FillCircle draws filled circle given by center point and radius.
//
// Center is the coords of pixel in center - circle with radius 3 will be 5 px wide.
func (m *Mono) FillCircle(color color.Color, center image.Point, radius int) {
	for x := 0; x < radius; x++ {
		for y := 0; y < radius; y++ {
			if x*x+y*y <= radius*radius {
				m.Set(center.X+x, center.Y+y, color)
				m.Set(center.X-x, center.Y+y, color)
				m.Set(center.X+x, center.Y-y, color)
				m.Set(center.X-x, center.Y-y, color)
			}
		}
	}
}

func (m *Mono) DrawString(color color.Color, text string, size float64, dot image.Point) {
	d := font.Drawer{
		Dst: m,
		Src: image.NewUniform(color),
		Face: truetype.NewFace(goFont, &truetype.Options{
			Size:    size,
			Hinting: font.HintingFull,
		}),
	}
	d.Dot = fixed.P(dot.X, dot.Y)
	d.DrawString(text)
}

// flips order of bits in byte
func flipByte(b byte) byte {
	// b = (b&0xF0)>>4 | (b&0x0F)<<4
	// b = (b&0xCC)>>2 | (b&0x33)<<2
	// b = (b&0xAA)>>1 | (b&0x55)<<1
	return bits.Reverse8(b)
}

// VerticalFlip flips image vertically (along horizontal axe)
// top to bottom and bottom to top
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

// HorizontalFlip flips image horizontally (along vetical axe)
// left to right and right to left
func (m *Mono) HorizontalFlip() {
	data := m.Bitmap()
	w := m.Width() / 8
	h := m.Height()
	for y := uint(0); y < h; y++ {
		for i := uint(0); i < w/2; i++ {
			data[y*w+i], data[y*w+w-1-i] = flipByte(data[y*w+w-1-i]), flipByte(data[y*w+i])
		}
	}
}

// RotateRight will rotate image 90 degrees clockwise.
//
// (center of rotation is center of largest square fitted to top left)
func (m *Mono) RotateRight() {
	w := m.Width()
	h := m.Height()

	n := NewMono(image.Rect(0, 0, int(h), int(w)))

	for y := 0; y < int(h); y++ {
		for x := 0; x < int(w); x++ {
			n.Set(int(h)-1-y, x, m.At(x, y))
		}
	}
	*m = n
}

// RotateLeft will rotate image 90 degrees counterclockwise.
//
// (center of rotation is center of largest square fitted to top left)
func (m *Mono) RotateLeft() {
	w := m.Width()
	h := m.Height()

	n := NewMono(image.Rect(0, 0, int(h), int(w)))

	for y := 0; y < int(h); y++ {
		for x := 0; x < int(w); x++ {
			n.Set(y, int(w)-1-x, m.At(x, y))
		}
	}
	*m = n
}

// Invert inverts colors in image
func (m *Mono) Invert() {
	for i := 4; i < len(*m); i++ {
		(*m)[i] = ^(*m)[i]
	}
}
