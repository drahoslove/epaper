package image

import (
	"github.com/drahoslav7/epaper"
	epd "github.com/drahoslav7/epaper/2in9"
	"testing"
)

func TestMono(t *testing.T) {

	epaper.Use(epd.Module)
	epaper.Init("full")
	m := NewMono(epd.Dimension.WIDTH, epd.Dimension.HEIGHT)
	m.StrokeRect(5, 10, 55, 110)
	m.FillRect(15, 15, 45, 100)

	epaper.SetFrame(m.Bitmap(), 0, 0, m.Width(), m.Height())
	epaper.DisplayFrame()
}
