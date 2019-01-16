package main

// sudo GOPATH=/home/pi/go /usr/local/go/bin/go test -v

import (
	"flag"
	"io/ioutil"
	"net/http"

	"github.com/drahoslove/epaper"
	epd "github.com/drahoslove/epaper/2in9"
)

func main() {
	displayBitmap := func(bitmap []byte) {
		width := uint(bitmap[0])<<8 + uint(bitmap[1])
		height := uint(bitmap[2])<<8 + uint(bitmap[3])

		epaper.SetFrame(bitmap[4:], 0, 0, width, height)
		epaper.DisplayFrame()
	}

	epaper.Use(epd.Module)

	filename := flag.String("file", "", "bitmap file to show")
	mode := flag.String("mode", "full", "refresh mode 'full' or 'partial'")
	port := flag.String("port", "", "port on which to listen for incomming bitmaps, eg. '6969'")
	clr := flag.Bool("clr", false, "clears display")

	flag.Parse()

	epaper.Init(*mode)

	println("FILE", *filename)
	println("MODE", *mode)
	println("SERVE", *port)

	if *clr {
		epaper.ClearFrame(epd.Ink.UNCOLORED)
		epaper.DisplayFrame()
	}

	if *filename != "" {
		fileContent, err := ioutil.ReadFile(*filename)
		if err != nil {
			panic(err)
		}
		displayBitmap(fileContent)
	}

	if *port != "" {
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
		http.ListenAndServe(":"+*port, nil)
	}
}
