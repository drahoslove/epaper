package image

import (
	"testing"
	"bytes"
	"image"
	"image/color"
	"image/draw"
)

func clearA(m *Mono, c color.Color) {
	draw.Draw(m, m.Bounds(), image.NewUniform(c), image.ZP, draw.Src)
}
func clearB(m *Mono, c color.Color) {
	r, g, b, _ := c.RGBA()
	Y := byte((r + g + b) / 3)
	for i := 4; i < len(*m); i++ {
		(*m)[i] = Y
	}
}
func clearC(m *Mono, c color.Color) {
	*m = append((*m)[:4], bytes.Repeat([]byte{byte(c.(color.Gray16).Y)}, len(*m)-4)...)
}

func fillRectA(m *Mono, rect image.Rectangle, c color.Color) {
	draw.Draw(m, rect, image.NewUniform(c), image.ZP, draw.Src)
}
func fillRectB(m *Mono, rect image.Rectangle, c color.Color) {
	w := rect.Dx()
	down := image.Pt(0, 1)
	for start := rect.Min; start.Y < rect.Max.Y; start = start.Add(down) {
		m.DrawHorizontalLine(c, start, w)
	}
}
func fillCircleA(m *Mono, color color.Color, center image.Point, radius int) {
	for x := -radius+1; x < radius; x++ {
		for y := -radius+1; y < radius; y++ {
			if x*x+y*y <= radius*radius {
				m.Set(center.X+x, center.Y+y, color)
			}
		}
	}
}
func fillCircleB(m *Mono, color color.Color, center image.Point, radius int) {
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

func BenchmarkClearA(b *testing.B) {
	m := NewMono(image.Rect(0, 0, 1024, 1024))
	for n := 0; n < b.N; n++ {
		clearA(&m, color.White)
		clearA(&m, color.Black)
	}
}
func BenchmarkClearB(b *testing.B) {
	m := NewMono(image.Rect(0, 0, 1024, 1024))
	for n := 0; n < b.N; n++ {
		clearB(&m, color.White)
		clearB(&m, color.Black)
	}
}
func BenchmarkClearC(b *testing.B) {
	m := NewMono(image.Rect(0, 0, 1024, 1024))
	for n := 0; n < b.N; n++ {
		clearC(&m, color.White)
		clearC(&m, color.Black)
	}
}
func BenchmarkFillRectA(b *testing.B) {
	m := NewMono(image.Rect(0, 0, 1024, 1024))
	for n := 0; n < b.N; n++ {
		fillRectA(&m, image.Rect(19, 23, 1001, 1017), color.White)
		fillRectA(&m, image.Rect(19, 23, 1001, 1017), color.Black)
	}
}
func BenchmarkFillRectB(b *testing.B) {
	m := NewMono(image.Rect(0, 0, 1024, 1024))
	for n := 0; n < b.N; n++ {
		fillRectB(&m, image.Rect(19, 23, 1001, 1017), color.White)
		fillRectB(&m, image.Rect(19, 23, 1001, 1017), color.Black)
	}
}


func BenchmarkFillCircleA(b *testing.B) {
	m := NewMono(image.Rect(0, 0, 1024, 1024))
	for n := 0; n < b.N; n++ {
		fillCircleA(&m, color.White, image.Pt(512, 512), 421)
		fillCircleA(&m, color.Black, image.Pt(512, 512), 421)
	}
}
func BenchmarkFillCircleB(b *testing.B) {
	m := NewMono(image.Rect(0, 0, 1024, 1024))
	for n := 0; n < b.N; n++ {
		fillCircleB(&m, color.White, image.Pt(512, 512), 421)
		fillCircleB(&m, color.Black, image.Pt(512, 512), 421)
	}
}
