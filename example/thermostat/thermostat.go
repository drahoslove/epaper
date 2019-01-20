package main

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/drahoslove/epaper"
	epd "github.com/drahoslove/epaper/2in9"
	eimage "github.com/drahoslove/epaper/image"
)

const zeroK = -273.15

func main() {
	temps := [][4]float32{}
	tempIndex := 0

	devices, err := ioutil.ReadDir("/sys/bus/w1/devices/")
	if err != nil {
		log.Fatal(err)
	}
	names := []string{}
	shortNames := []string{}

	for _, device := range devices {
		name := device.Name()
		if name[:3] == "28-" { // filter thermostats
			shortName := name[5:7] + "/" + name[13:15]
			names = append(names, name)
			shortNames = append(shortNames, shortName)
			temps = append(temps, [4]float32{zeroK, zeroK, zeroK, zeroK})
		}
	}

	go func() {
		epaper.Setup()
		defer epaper.Teardown()
		epaper.Init("full")
		epaper.Clear(255)
		epaper.Clear(255)
		epaper.Init("partial")
		for t := range time.Tick(time.Second * 1) {
			render(shortNames, temps, t)
		}
	}()

	for _ = range time.Tick(time.Second * 5) {
		for i, name := range names {
			temp := readTemp(name)
			// fmt.Println(time.Now().String()[11:22], shortNames[i], temp)
			temps[i][tempIndex] = temp
		}
		tempIndex++
		tempIndex %= 4
	}
}

func readTemp(devName string) float32 {
	data, err := ioutil.ReadFile("/sys/bus/w1/devices/" + devName + "/w1_slave")
	if err != nil {
		return zeroK
	}
	hb := data[3:5]
	lb := data[0:2]
	val, err := hex.DecodeString(string(hb) + string(lb))
	if err != nil {
		return zeroK
	}
	temp := float32(int(val[0])<<8+int(uint(val[1]))) / 16
	return temp
}

func render(names []string, temps [][4]float32, t time.Time) {
	if t.Second() == 0 {
		epaper.Init("full")
		defer epaper.Init("partial")
	}
	irect := image.Rect(0, 0, int(epd.Dimension.HEIGHT), int(epd.Dimension.WIDTH))
	img := eimage.NewMono(irect)
	img.Clear(image.White)
	img.FillRect(image.Black, irect)

	for i, name := range names {
		avgTemp := float64((temps[i][0] + temps[i][1] + temps[i][2] + temps[i][3]) / 4)
		tempStr := strconv.FormatFloat(avgTemp, 'f', 2, 32)

		fmt.Println(name, tempStr)

		img.DrawString(color.White, name+":", 26, image.Pt(5, 30*(i+1)))
		img.DrawString(color.White, tempStr+"°C", 26, image.Pt(105, 30*(i+1)))
	}

	img.RotateRight()
	img.DrawString(color.White, t.String()[11:19], 28, image.Pt(5, 280))
	epaper.Display(img.Bitmap(), 0, 0, img.Width(), img.Height())
	// epaper.Sleep()
}
