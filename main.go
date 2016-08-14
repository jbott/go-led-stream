package main

import (
	"bytes"
	"encoding/binary"
	"github.com/jacobsa/go-serial/serial"
	"hash/crc32"
	"log"
	"time"
)

func main() {
	log.Print("Opening serial port...")

	serial_options := serial.OpenOptions{
		PortName:              "/dev/cu.usbmodem31",
		BaudRate:              57600,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       2,
		InterCharacterTimeout: 0,
	}

	port, err := serial.Open(serial_options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	defer port.Close()

	log.Print("Waiting for arduino to reset...")
	time.Sleep(5 * time.Second)
	log.Print("DONE")

	cmd_buf := new(bytes.Buffer)
	cmd_buf.Write(SetFlags(FLAG_PAUSE_GHUE))
	cmd_buf.Write(Rainbow(5 * 1000))
	cmd_buf.Write(SetFlags(0))
	cmd_buf.Write(Rainbow(10 * 1000))
	cmd_buf.Write(SetFlags(FLAG_PAUSE_GHUE))
	cmd_buf.Write(Rainbow(5 * 1000))
	cmd_buf.Write(SetFlags(FLAG_REVERSE_GHUE))
	cmd_buf.Write(Rainbow(10 * 1000))
	packet := WrapHeaderCRC(PACKET_SET_CMDS, cmd_buf.Bytes())

	log.Print(packet)
	port.Write(packet)
	time.Sleep(10000 * time.Second)
}

type Color struct {
	R uint8
	G uint8
	B uint8
}

func CreateColor(r int, g int, b int) Color {
	return Color{R: uint8(r), G: uint8(g), B: uint8(b)}
}

func (c *Color) UpdateColor(r uint8, g uint8, b uint8) {
	c.R = r
	c.G = g
	c.B = b
}

const (
	CMD_SET_FLAGS = iota
	CMD_RAINBOW   = iota
	CMD_FADE_RGB  = iota
)

type Command struct {
	index       uint8
	duration    uint16
	params_size uint8
}

func NewCommand(index uint8, duration uint16) Command {
	return Command{index, duration, 0}
}

func (c Command) Bytes() []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, c)
	if err != nil {
		log.Fatalf("Error writing cmd struct: %v", err)
	}
	return buf.Bytes()
}

func (c Command) WithParams(params interface{}) []byte {
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

type FadeParams struct {
	start_color Color
	end_color   Color
}

func FadeRGB(start_color Color, end_color Color, duration uint16) []byte {
	return NewCommand(CMD_FADE_RGB, duration).WithParams(
		FadeParams{start_color, end_color})
}

func Rainbow(duration uint16) []byte {
	return NewCommand(CMD_RAINBOW, duration).Bytes()
}

const (
	FLAG_PAUSE_GHUE   = 1 << iota
	FLAG_REVERSE_GHUE = 1 << iota
)

type SetFlagsParams struct {
	flags uint8
}

func SetFlags(flags uint8) []byte {
	if flags == 0 {
		return NewCommand(CMD_SET_FLAGS, 0).Bytes()
	} else {
		return NewCommand(CMD_SET_FLAGS, 0).WithParams(
			SetFlagsParams{flags})
	}
}

const (
	PACKET_NONE     = iota
	PACKET_SET_CMDS = iota
)

type PacketHeader struct {
	magic         uint32
	packet_type   uint8
	packet_length uint8
}

func WrapHeaderCRC(packet_type uint8, data []byte) []byte {
	header := PacketHeader{0xDEADBEEF, packet_type, uint8(len(data))}
	checksum := crc32.ChecksumIEEE(data)

	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, header)
	if err != nil {
		log.Fatalf("Error writing header: %v", err)
	}

	n, err := buf.Write(data)
	if err != nil {
		log.Fatalf("Error writing data: %v", err)
	}
	if n != len(data) {
		log.Fatal("Did not write correct amount of data")
	}

	err = binary.Write(buf, binary.LittleEndian, checksum)
	if err != nil {
		log.Fatalf("Error writing checksum: %v", err)
	}

	return buf.Bytes()
}
