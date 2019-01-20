package epaper

import (
	"bytes"
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"math"
	"math/rand"
	"os"
	"time"
)

const (
	dcPin    = rpio.Pin(25) // OUT 0 = command, 1 = data
	resetPin = rpio.Pin(22) // OUT 0 = reset
	busyPin  = rpio.Pin(24) // IN  0 = busy
)

// setup gpio and SPI interface
func Setup() {
	err := rpio.Open()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	resetPin.Output()
	dcPin.Output()
	busyPin.Input()
	busyPin.PullDown()

	resetPin.High()

	// SETUP SPI:
	err = rpio.SpiBegin(rpio.Spi0)
	// freq 128 divider - default
	// chip select CE0 - default
	// ce enable low - implicit
	// mode 0 - implicit
	// msb first - implicit
	if err != nil {
		fmt.Print(err)
		os.Exit(2)
	}
}

// teardown gpio and SPI interface
func Teardown() {
	rpio.SpiEnd(rpio.Spi0)
	rpio.Close()
}

var (
	lut Lut
	cmd Cmd
	dim Dim
	ink Ink
) // global vars holding used values, fill in with Use

func Use(e Module) {
	lut = e.Lut
	cmd = e.Cmd
	dim = e.Dim
	ink = e.Ink
}

func Init(update string) {
	Reset()
	SendCommand(cmd.DRIVER_OUTPUT_CONTROL)
	SendData(
		byte((dim.HEIGHT-1)&0xFF),
		byte((dim.HEIGHT-1)>>8),
		0x00, // GD = 0; SM = 0; TB = 0;
	)
	SendCommand(cmd.BOOSTER_SOFT_START_CONTROL)
	// SendData(0xD7, 0xD6, 0x9D)
	SendData(0xCF, 0xCE, 0x8D)
	SendCommand(cmd.WRITE_VCOM_REGISTER)
	SendData(0x7c) // VCOM 7C // 8a
	SendCommand(cmd.SET_DUMMY_LINE_PERIOD)
	SendData(0x1A) // 4 dummy lines per gate
	SendCommand(cmd.SET_GATE_TIME)
	SendData(0x08) // 2us per line
	SendCommand(cmd.DATA_ENTRY_MODE_SETTING)
	SendData(0x03) // X increment Y increment
	if update == "partial" {
		SetLut(lut.PARTIAL)
	}
	if update == "full" {
		SetLut(lut.FULL)
	}
}

func SendCommand(cmd byte) {
	dcPin.Low()
	rpio.SpiTransmit(cmd)
}

func SendData(data ...byte) {
	dcPin.High()
	rpio.SpiTransmit(data...)
}

func WaitUntilIdle() {
	for busyPin.Read() == rpio.High { // doc say Low == busy, but it is the oposite
		time.Sleep(time.Millisecond * 50)
	}
}

func Reset() {
	resetPin.Low()
	time.Sleep(time.Millisecond * 100)
	resetPin.High()
	time.Sleep(time.Millisecond * 100)
}

func SetLut(lut []byte) {
	SendCommand(cmd.WRITE_LUT_REGISTER)
	SendData(lut...)
}

func Clear(color byte) {
	h := dim.HEIGHT
	w := dim.WIDTH
	SetMemoryArea(0, 0, w-1, h-1)
	SetMemoryPointer(0, 0)
	SendCommand(cmd.WRITE_RAM)
	/* send the color data */
	var img = bytes.Repeat([]byte{color}, int(inBytes(w)*h))
	SendData(img...)
	SwapFrame()
}

func Randomize() {
	h := dim.HEIGHT
	w := dim.WIDTH
	SetMemoryArea(0, 0, w-1, h-1)
	SetMemoryPointer(0, 0)
	SendCommand(cmd.WRITE_RAM)
	/* send the color data */
	var img = make([]byte, int(inBytes(w)*h))
	for i := range img {
		img[i] = byte(rand.Int())
	}
	SendData(img...)
	SwapFrame()
}

// Will display bitmap
// if image is larger, it will be cropped
func Display(img []byte, x, y int, imgWidth, imgHeight uint) {
	if len(img) < int(imgHeight*inBytes(imgWidth)) {
		fmt.Print("bitmap too small")
		return
	}
	/* x point must be the multiple of 8 or the last 3 bits will be ignored */
	xEnd := uint(math.Min(
		float64(x+int(imgWidth-1)),
		float64(dim.WIDTH-1)),
	)
	yEnd := uint(math.Min(
		float64(y+int(imgHeight-1)),
		float64(dim.HEIGHT-1)),
	)
	xStart := uint(math.Max(float64(x), 0))
	yStart := uint(math.Max(float64(y), 0))

	SetMemoryArea(xStart, yStart, xEnd, yEnd)
	SetMemoryPointer(xStart, yStart)
	SendCommand(cmd.WRITE_RAM)
	/* send the img data, line by line */
	rowsToCrop := (yStart + imgHeight - dim.HEIGHT)
	if y < 0 { // crop top
		img = img[inBytes(imgWidth)*uint(-y):]
		rowsToCrop -= uint(-y)
	}
	for len(img) > 0 && len(img) > int(inBytes(imgWidth)*rowsToCrop) {
		if x >= 0 {
			SendData(img[0:inBytes(xEnd-xStart)]...)
		} else { // crop left part
			SendData(img[inBytes(uint(-x)):inBytes(xEnd+uint(-x))]...)
		}
		img = img[inBytes(imgWidth):] // next line
	}
	SwapFrame()
}

// Will swap back frame with front frame and displays what's on it
func SwapFrame() {
	SendCommand(cmd.DISPLAY_UPDATE_CONTROL_2)
	SendData(0xC4)
	SendCommand(cmd.MASTER_ACTIVATION)
	SendCommand(cmd.TERMINATE_FRAME_READ_WRITE)
	WaitUntilIdle()
}

func SetMemoryArea(x_start, y_start, x_end, y_end uint) {
	SendCommand(cmd.SET_RAM_X_ADDRESS_START_END_POSITION)
	/* x point must be the multiple of 8 or the last 3 bits will be ignored */
	SendData(byte(x_start >> 3))
	SendData(byte(x_end >> 3))
	SendCommand(cmd.SET_RAM_Y_ADDRESS_START_END_POSITION)
	SendData(byte(y_start))
	SendData(byte(y_start >> 8))
	SendData(byte(y_end))
	SendData(byte(y_end >> 8))
	WaitUntilIdle()
}

func SetMemoryPointer(x, y uint) {
	SendCommand(cmd.SET_RAM_X_ADDRESS_COUNTER)
	/* x point must be the multiple of 8 or the last 3 bits will be ignored */
	SendData(byte(x >> 3))
	SendCommand(cmd.SET_RAM_Y_ADDRESS_COUNTER)
	SendData(byte(y))
	SendData(byte(y >> 8))
	WaitUntilIdle()
}

func Sleep() {
	SendCommand(cmd.DEEP_SLEEP_MODE)
	SendData(1)
	// WaitUntilIdle()
}

// division by eight but round up
func inBytes(n uint) uint {
	if n%8 == 0 {
		return n / 8
	} else {
		return n/8 + 1
	}
}
