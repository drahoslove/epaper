package image

import (
	// "fmt"
	"github.com/drahoslav7/epaper"
	epd "github.com/drahoslav7/epaper/2in9"
	"image/color"
	"testing"
)

func TestMono(t *testing.T) {

	epaper.Use(epd.Module)
	epaper.Init("full")
	m := NewMono(epd.Dimension.HEIGHT, epd.Dimension.WIDTH)
	m.Clear(color.White)
	m.Set(1, 1, color.Black)
	m.Set(1, 2, color.White)
	m.Set(1, 3, color.Black)
	m.Set(1, 4, color.Black)
	m.StrokeRect(5, 5, 45, 65)
	m.FillRect(10, 10, 40, 60)

	for i, s := range []string{"Hello World", "12:59 | 23.5°C"} {
		size := 38
		m.DrawString(s, float64(size), 50, size*(i+1))
	}

	// m.HorizontalFlip()
	// m.VerticalFlip()
	m.RotateRight()

	for i, n := 0, 2; i < 10; i, n = i+1, n+n/2 {
		m.DrawHorizontalLine(n+50, 3*i, 50-3*i)
	}

	m.DrawString("abc", 24, 5, 30)

	epaper.SetFrame(m.Bitmap(), 0, 0, m.Width(), m.Height())
	epaper.DisplayFrame()
}
