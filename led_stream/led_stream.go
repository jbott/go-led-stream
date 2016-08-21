package led_stream

import (
	"bytes"
	"encoding/binary"
	"github.com/lucasb-eyer/go-colorful"
	"log"
	"strconv"
)

type Entry struct {
	Cmd      string `json:"cmd"`
	Duration string `json:"duration"`

	Flag_pause_ghue   string `json:"flag_pause_ghue"`
	Flag_reverse_ghue string `json:"flag_reverse_ghue"`

	Color string `json:"color"`

	StartColor string `json:"start_color"`
	EndColor   string `json:"end_color"`

	FadeBy string `json:"fade_by"`

	Count string `json:"count"`
}

type Color struct {
	R uint8
	G uint8
	B uint8
}

func NewColor(color string) Color {
	from_color, err := colorful.Hex(color)
	if err != nil {
		log.Fatal(err)
	}
	return Color{
		uint8(255 * from_color.R),
		uint8(255 * from_color.G),
		uint8(255 * from_color.B),
	}
}

const (
	CMD_SET_FLAGS      = iota
	CMD_OFF            = iota
	CMD_FILL_SOLID_RGB = iota
	CMD_RAINBOW        = iota
	CMD_FADE_RGB       = iota
	CMD_FADE_TO_BLACK  = iota
	CMD_CONFETTI       = iota
)

const (
	STR_CMD_SET_FLAGS      = "SET_FLAGS"
	STR_CMD_OFF            = "OFF"
	STR_CMD_FILL_SOLID_RGB = "FILL_SOLID_RGB"
	STR_CMD_RAINBOW        = "RAINBOW"
	STR_CMD_FADE_RGB       = "FADE_RGB"
	STR_CMD_FADE_TO_BLACK  = "FADE_TO_BLACK"
	STR_CMD_CONFETTI       = "CONFETTI"
)

type Command struct {
	index       uint8
	duration    uint16
	params_size uint8
}

type FadeParams struct {
	start_color Color
	end_color   Color
}

const (
	FLAG_PAUSE_GHUE   = 1 << iota
	FLAG_REVERSE_GHUE = 1 << iota
)

type SetFlagsParams struct {
	flags uint8
}

type FillSolidParams struct {
	color Color
}

type FadeToBlackParams struct {
	fadeBy uint8
}

type ConfettiParams struct {
	count  uint8
	fadeBy uint8
}

func EntryToBytes(e Entry) []byte {
	c := Command{}

	if duration, err := strconv.ParseUint(e.Duration, 10, 16); err == nil {
		c.duration = uint16(duration)
	}

	switch {
	case e.Cmd == STR_CMD_SET_FLAGS:
		c.index = CMD_SET_FLAGS
		c.duration = 0

		params := SetFlagsParams{}

		if e.Flag_pause_ghue == "on" {
			params.flags += FLAG_PAUSE_GHUE
		}
		if e.Flag_reverse_ghue == "on" {
			params.flags += FLAG_REVERSE_GHUE
		}

		if params.flags != 0 {
			return c.WithParamsBytes(params)
		}

	case e.Cmd == STR_CMD_OFF:
		c.index = CMD_OFF

	case e.Cmd == STR_CMD_FILL_SOLID_RGB:
		c.index = CMD_FILL_SOLID_RGB
		params := FillSolidParams{NewColor(e.Color)}
		return c.WithParamsBytes(params)

	case e.Cmd == STR_CMD_RAINBOW:
		c.index = CMD_RAINBOW

	case e.Cmd == STR_CMD_FADE_RGB:
		c.index = CMD_FADE_RGB
		params := FadeParams{NewColor(e.StartColor), NewColor(e.EndColor)}
		return c.WithParamsBytes(params)

	case e.Cmd == STR_CMD_FADE_TO_BLACK:
		c.index = CMD_FADE_TO_BLACK
		if fadeBy, err := strconv.ParseUint(e.FadeBy, 10, 8); err == nil {
			return c.WithParamsBytes(FadeToBlackParams{uint8(fadeBy)})
		}

	case e.Cmd == STR_CMD_CONFETTI:
		c.index = CMD_CONFETTI
		params := ConfettiParams{}
		if fadeBy, err := strconv.ParseUint(e.FadeBy, 10, 8); err == nil {
			params.fadeBy = uint8(fadeBy)
		}
		if count, err := strconv.ParseUint(e.Count, 10, 8); err == nil {
			params.count = uint8(count)
		}
		return c.WithParamsBytes(params)

	}

	return c.Bytes()
}

func (c Command) WithParamsBytes(params interface{}) []byte {
	c.params_size = uint8(binary.Size(params))
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, c)
	if err != nil {
		log.Fatalf("Error writing cmd struct: %v", err)
	}

	if params != nil {
		err = binary.Write(buf, binary.LittleEndian, params)
		if err != nil {
			log.Fatalf("Error writing params struct: %v", err)
		}
	}

	return buf.Bytes()
}

func (c Command) Bytes() []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, c)
	if err != nil {
		log.Fatalf("Error writing cmd struct: %v", err)
	}
	return buf.Bytes()
}
