package epaper

// sudo GOPATH=/home/pi/go /usr/local/go/bin/go test -v

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	model "github.com/drahoslav7/epaper/2in9"
)

func TestReset(t *testing.T) {
	displayBitmap := func(bitmap []byte) {
		width := int(bitmap[0])<<8 + int(bitmap[1])
		height := int(bitmap[2])<<8 + int(bitmap[3])

		SetFrame(bitmap[4:], 0, 0, width, height)
		// SetFrame(bitmap, 0, 0, model.Spec.Res.WIDTH, model.Spec.Res.HEIGHT)
		DisplayFrame()
	}

	Use(model.Spec)

	filename := os.Getenv("FILE")
	mode := os.Getenv("MODE")
	port := os.Getenv("SERVE")

	Init(mode)

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
			Init("full")
			bodyContent, err := ioutil.ReadAll(r.Body)
			if err != nil {
				println(err)
			}
			displayBitmap(bodyContent)
			w.Header().Add("Access-Control-Allow-Origin", "*")
			Init("partial")
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
			RandomizeFrame()
			DisplayFrame()

			time.Sleep(time.Second)
		}
	}

	if os.Getenv("BLINK") != "" {
		tick := time.Tick(time.Millisecond * 300)
		for {
			ClearFrame(ink.UNCOLORED)
			DisplayFrame()
			fmt.Println(<-tick)

			ClearFrame(ink.COLORED)
			DisplayFrame()
			fmt.Println(<-tick)
		}
	}
}
