/*

	Driver for waveshare 2.9" e-paper display
	https://www.waveshare.com/w/upload/e/e6/2.9inch_e-Paper_Datasheet.pdf

*/
package model2in9

import (
	"github.com/drahoslove/epaper/spec"
)

var Module = spec.Module{
	Ink,
	Dimension,
	lut,
	command,
}

// Colors
var Ink = spec.Ink{
	byte(0),
	^byte(0),
}

// Display dimension
var Dimension = spec.Dim{
	128,
	296,
}

// commands
var command = spec.Cmd{
	DRIVER_OUTPUT_CONTROL:                0x01,
	BOOSTER_SOFT_START_CONTROL:           0x0C,
	GATE_SCAN_START_POSITION:             0x0F,
	DEEP_SLEEP_MODE:                      0x10,
	DATA_ENTRY_MODE_SETTING:              0x11,
	SW_RESET:                             0x12,
	TEMPERATURE_SENSOR_CONTROL:           0x1A,
	MASTER_ACTIVATION:                    0x20,
	DISPLAY_UPDATE_CONTROL_1:             0x21,
	DISPLAY_UPDATE_CONTROL_2:             0x22,
	WRITE_RAM:                            0x24,
	WRITE_VCOM_REGISTER:                  0x2C,
	WRITE_LUT_REGISTER:                   0x32,
	SET_DUMMY_LINE_PERIOD:                0x3A,
	SET_GATE_TIME:                        0x3B,
	BORDER_WAVEFORM_CONTROL:              0x3C,
	SET_RAM_X_ADDRESS_START_END_POSITION: 0x44,
	SET_RAM_Y_ADDRESS_START_END_POSITION: 0x45,
	SET_RAM_X_ADDRESS_COUNTER:            0x4E,
	SET_RAM_Y_ADDRESS_COUNTER:            0x4F,
	TERMINATE_FRAME_READ_WRITE:           0xFF,
}

var lut = spec.Lut{
	FULL: []byte{
		0x02, 0x02, 0x01, 0x11, 0x12, 0x12, 0x22, 0x22,
		0x66, 0x69, 0x69, 0x59, 0x58, 0x99, 0x99, 0x88,
		0x00, 0x00, 0x00, 0x00, 0xF8, 0xB4, 0x13, 0x51,
		0x35, 0x51, 0x51, 0x19, 0x01, 0x00,
	},
	PARTIAL: []byte{
		0x10, 0x18, 0x18, 0x08, 0x18, 0x18, 0x08, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x13, 0x14, 0x44, 0x12,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	},
}
