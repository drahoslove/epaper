package model

type Model struct {
	Ink
	Res
	Lut
	Cmd
}

type Ink struct {
	COLORED   byte
	UNCOLORED byte
}

type Res struct {
	WIDTH  int
	HEIGHT int
}

type Lut struct {
	FULL    []byte
	PARTIAL []byte
}

type Cmd struct {
	DRIVER_OUTPUT_CONTROL                byte
	BOOSTER_SOFT_START_CONTROL           byte
	GATE_SCAN_START_POSITION             byte
	DEEP_SLEEP_MODE                      byte
	DATA_ENTRY_MODE_SETTING              byte
	SW_RESET                             byte
	TEMPERATURE_SENSOR_CONTROL           byte
	MASTER_ACTIVATION                    byte
	DISPLAY_UPDATE_CONTROL_1             byte
	DISPLAY_UPDATE_CONTROL_2             byte
	WRITE_RAM                            byte
	WRITE_VCOM_REGISTER                  byte
	WRITE_LUT_REGISTER                   byte
	SET_DUMMY_LINE_PERIOD                byte
	SET_GATE_TIME                        byte
	BORDER_WAVEFORM_CONTROL              byte
	SET_RAM_X_ADDRESS_START_END_POSITION byte
	SET_RAM_Y_ADDRESS_START_END_POSITION byte
	SET_RAM_X_ADDRESS_COUNTER            byte
	SET_RAM_Y_ADDRESS_COUNTER            byte
	TERMINATE_FRAME_READ_WRITE           byte
}
