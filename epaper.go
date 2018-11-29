package epaper

import (
	"bytes"
	"fmt"
	"github.com/drahoslav7/epaper/model"
	"github.com/stianeikeland/go-rpio"
	"math/rand"
	"os"
	"time"
)

const (
	dcPin    = rpio.Pin(25) // OUT 0 = command, 1 = data
	resetPin = rpio.Pin(22) // OUT 0 = reset
	busyPin  = rpio.Pin(24) // IN  0 = busy
)

// init interface
func init() {
	err := rpio.Open()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	resetPin.Output()
	dcPin.Output()
	busyPin.Input()

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

var (
	lut model.Lut
	cmd model.Cmd
	res model.Res
	ink model.Ink
) // global vars holding used values, fill in with Use

func Use(e model.Model) {
	lut = e.Lut
	cmd = e.Cmd
	res = e.Res
	ink = e.Ink
}

func Init(update string) {
	Reset()
	SendCommand(cmd.DRIVER_OUTPUT_CONTROL)
	SendData(
		byte((res.HEIGHT-1)&0xFF),
		byte((res.HEIGHT-1)>>8),
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

func ClearFrame(color byte) {
	SetMemoryArea(0, 0, res.WIDTH-1, res.HEIGHT-1)
	SetMemoryPointer(0, 0)
	SendCommand(cmd.WRITE_RAM)
	/* send the color data */
	var img = bytes.Repeat([]byte{color}, res.WIDTH/8*res.HEIGHT)
	SendData(img...)
}

func RandomizeFrame() {
	SetMemoryArea(0, 0, res.WIDTH-1, res.HEIGHT-1)
	SetMemoryPointer(0, 0)
	SendCommand(cmd.WRITE_RAM)
	/* send the color data */
	var img = make([]byte, res.WIDTH/8*res.HEIGHT)
	for i := range img {
		img[i] = byte(rand.Int())
	}
	SendData(img...)
}

func SetFrame(img []byte, x, y, imgWidth, imgHeight int) {
	var (
		xEnd int
		yEnd int
	)

	if len(img) == 0 || x < 0 || imgWidth < 0 || y < 0 || imgHeight < 0 {
		return
	}
	/* x point must be the multiple of 8 or the last 3 bits will be ignored */
	x &= 0xF8 // 11111000
	imgWidth &= 0xF8
	if x+imgWidth >= res.WIDTH {
		xEnd = res.WIDTH - 1
	} else {
		xEnd = x + imgWidth - 1
	}
	if y+imgHeight >= res.HEIGHT {
		yEnd = res.HEIGHT - 1
	} else {
		yEnd = y + imgHeight - 1
	}
	SetMemoryArea(x, y, xEnd, yEnd)
	SetMemoryPointer(x, y)
	SendCommand(cmd.WRITE_RAM)
	/* send the img data */
	for j := 0; y < yEnd; j++ {
		x = 0
		for i := 0; x < xEnd; i++ {
			SendData(img[i+j*(imgWidth/8)])
			x += 8
		}
		y++
	}
}

func DisplayFrame() {
	SendCommand(cmd.DISPLAY_UPDATE_CONTROL_2)
	SendData(0xC4)
	SendCommand(cmd.MASTER_ACTIVATION)
	SendCommand(cmd.TERMINATE_FRAME_READ_WRITE)
	WaitUntilIdle()
}

func SetMemoryArea(x_start, y_start, x_end, y_end int) {
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

func SetMemoryPointer(x, y int) {
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
	WaitUntilIdle()
}
