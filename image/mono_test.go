// Run by:
// go test -c -o test -v github.com/drahoslove/epaper/image && sudo ./test -test.v
package image

import (
	// "fmt"
	"github.com/drahoslove/epaper"
	epd "github.com/drahoslove/epaper/2in9"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

var (
	black = color.Black
	white = color.White
)

func TestMono(t *testing.T) {

	epaper.Setup()
	defer epaper.Teardown()
	epaper.Init("full")
	defer epaper.Sleep()

	m := NewMono(image.Rect(0, 0, int(epd.Dimension.HEIGHT), int(epd.Dimension.WIDTH)))
	m.Clear(white)
	m.Set(1, 1, black)
	m.Set(1, 2, white)
	m.Set(1, 3, black)
	m.Set(1, 4, black)
	rect := image.Rect(5, 5, 45, 65)
	m.StrokeRect(black, rect)
	m.FillRect(black, rect.Inset(5))

	size := 36
	dot := image.Pt(50, 0)
	for i, s := range []string{"Hello World", "12:59 | 23.5°C"} {
		dot = dot.Add(image.Pt(0, size))
		m.DrawString(black, s, float64(size-6*i), dot)
	}

	// m.HorizontalFlip()
	// m.VerticalFlip()
	m.RotateRight()

	start := image.Pt(0, 50)
	length := 50
	for n := 2; length > 0; n += n / 2 {
		m.DrawHorizontalLine(black, start, length)
		start = start.Add(image.Pt(3, n))
		length -= 6
	}

	m.DrawString(black, "abc", 24, image.Pt(5, 30))
	m.DrawString(white, "abc", 24, image.Pt(70, 30))

	dot = image.Pt(5, 130)
	for i, s := range []string{"26", "24", "22", "20", "18", "16", "14", "12 - Hello", "10 - Hello", "08 - Hello", "06 - Hello"} {
		m.DrawString(black, s, float64(26-i*2), dot)
		dot = dot.Add(image.Pt(0, 25-2*i))
	}

	center := m.Bounds().Inset(11).Size() // 21 px from bottom right corner
	m.FillCircle(black, center, 11)
	m.StrokeCircle(black, center, 16)
	m.StrokeCircle(white, center, 3)
	m.Set(center.X, center.Y, white)

	m.Invert()

	// show bitmap on display
	epaper.Display(m.Bitmap(), 0, 0, m.Width(), m.Height())

	// save bitmap to png file

	f, _ := os.Create("image.png")
	png.Encode(f, m)
	f.Close()
}

func TestLines(t *testing.T) {
	epaper.Setup()
	defer epaper.Teardown()
	epaper.Init("full")
	defer epaper.Sleep()

	irect := image.Rect(0, 0, int(epd.Dimension.HEIGHT), int(epd.Dimension.WIDTH))
	img := NewMono(irect)
	img.Clear(image.White)
	img.FillRect(image.Black, irect)
	img.StrokeRect(image.White, irect.Inset(2))

	center := image.Pt(222, 64)
	img.DrawLine(color.White, center, center.Add(image.Pt(11, 0)))
	img.DrawLine(color.White, center, center.Add(image.Pt(11, 5)))
	img.DrawLine(color.White, center, center.Add(image.Pt(11, 11)))
	img.DrawLine(color.White, center, center.Add(image.Pt(5, 11)))
	img.DrawLine(color.White, center, center.Add(image.Pt(0, 11)))
	img.DrawLine(color.White, center, center.Add(image.Pt(-5, 11)))
	img.DrawLine(color.White, center, center.Add(image.Pt(-11, 11)))
	img.DrawLine(color.White, center, center.Add(image.Pt(-11, 5)))
	img.DrawLine(color.White, center, center.Add(image.Pt(-11, 0)))
	img.DrawLine(color.White, center, center.Add(image.Pt(-11, -5)))
	img.DrawLine(color.White, center, center.Add(image.Pt(-11, -11)))
	img.DrawLine(color.White, center, center.Add(image.Pt(-5, -11)))
	img.DrawLine(color.White, center, center.Add(image.Pt(0, -11)))
	img.DrawLine(color.White, center, center.Add(image.Pt(5, -11)))
	img.DrawLine(color.White, center, center.Add(image.Pt(11, -11)))
	img.DrawLine(color.White, center, center.Add(image.Pt(11, -5)))
	img.StrokeRect(color.White, image.Rectangle{
		center.Sub(image.Pt(12, 12)),
		center.Add(image.Pt(12, 12)),
	})
	center = center.Add(image.Pt(20, 0))
	for i := 0; i < 10; i++ {
		img.DrawLine(color.White, center, center.Add(image.Pt(i, 0)))
		center = center.Add(image.Pt(0, 2))
	}
	for i := 0; i < 5; i++ {
		center = center.Add(image.Pt(0, 6))
		img.FillRect(color.White, image.Rectangle{
			center.Add(image.Pt(0, 0)),
			center.Add(image.Pt(i, i)),
		})
		img.StrokeRect(color.White, image.Rectangle{
			center.Add(image.Pt(6+0, 0)),
			center.Add(image.Pt(6+i, i)),
		})
	}
	img.RotateRight()
	epaper.Display(img.Bitmap(), 0, 0, img.Width(), img.Height())

}
