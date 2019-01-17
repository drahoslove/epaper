package epaper_test

// sudo GOPATH=/home/pi/go /usr/local/go/bin/go test -v

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/drahoslove/epaper"
	_ "github.com/drahoslove/epaper/2in9" // will be used automatically
)

func TestReset(t *testing.T) {
	displayBitmap := func(bitmap []byte) {
		width := uint(bitmap[0])<<8 + uint(bitmap[1])
		height := uint(bitmap[2])<<8 + uint(bitmap[3])

		epaper.SetFrame(bitmap[4:], 0, 0, width, height)
		epaper.DisplayFrame()

		// SetFrame(bitmap, 0, 0, model.Spec.Res.WIDTH, model.Spec.Res.HEIGHT)
	}

	epaper.Setup()
	defer epaper.Teardown()

	filename := os.Getenv("FILE")
	mode := os.Getenv("MODE")
	port := os.Getenv("SERVE")

	epaper.Init(mode)

	println("FILE", filename)
	println("MODE", mode)
	println("SERVE", port)

	if filename != "" {
		fileContent, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		displayBitmap(fileContent)
	}

	if port != "" {
		lock := make(chan bool, 1)
		http.HandleFunc("/epd/full", func(w http.ResponseWriter, r *http.Request) {
			lock <- true
			epaper.Init("full")
			bodyContent, err := ioutil.ReadAll(r.Body)
			if err != nil {
				println(err)
			}
			displayBitmap(bodyContent)
			w.Header().Add("Access-Control-Allow-Origin", "*")
			epaper.Init("partial")
			<-lock
		})
		http.HandleFunc("/epd/partial", func(w http.ResponseWriter, r *http.Request) {
			lock <- true
			bodyContent, err := ioutil.ReadAll(r.Body)
			if err != nil {
				println(err)
			}
			displayBitmap(bodyContent)
			w.Header().Add("Access-Control-Allow-Origin", "*")
			<-lock
		})
		http.ListenAndServe(":"+port, nil)
	}

	if os.Getenv("NOISE") != "" {
		time.Sleep(time.Second)

		for {
			epaper.RandomizeFrame()
			epaper.DisplayFrame()

			time.Sleep(time.Second)
		}
	}

	if os.Getenv("BLINK") != "" {
		tick := time.Tick(time.Millisecond * 300)
		for {
			epaper.ClearFrame(0xFF)
			epaper.DisplayFrame()
			fmt.Println(<-tick)

			epaper.ClearFrame(0x00)
			epaper.DisplayFrame()
			fmt.Println(<-tick)
		}
	}
}
