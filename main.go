package main

import (
	"bytes"
	"encoding/binary"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"hash/crc32"
	"log"
)

func main() {
	log.Print("Connecting to MQTT...")

	opts := MQTT.NewClientOptions().AddBroker("tcp://192.168.1.21:1883")
	opts.SetClientID("go-simple")

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	purple := Color{127, 0, 100}
	dark_purple := purple.Scale(0.1)

	c := new(bytes.Buffer)
	c.Write(FadeRGB(dark_purple, purple, 5*1000))
	c.Write(FillSolidRGB(purple, 10*1000))
	c.Write(FadeRGB(purple, dark_purple, 5*1000))
	c.Write(FillSolidRGB(dark_purple, 10*1000))
	c.Write(FadeRGB(dark_purple, purple, 5*1000))

	c.Write(Rainbow(20 * 1000))
	c.Write(SetFlags(FLAG_REVERSE_GHUE))
	c.Write(Rainbow(20 * 1000))

	c.Write(FadeToBlack(1, 2*1000))

	c.Write(ConfettiDefault(20 * 1000))

	packet := WrapHeaderCRC(PACKET_SET_CMDS, c.Bytes())

	log.Printf("Packet length: %v", len(packet))
	log.Print(packet)
	token := client.Publish("device/led0/cmd", 0, true, packet)
	token.Wait()
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

func (c Color) Scale(scalar float32) Color {
	c.R = uint8(float32(c.R) * scalar)
	c.G = uint8(float32(c.G) * scalar)
	c.B = uint8(float32(c.B) * scalar)
	return c
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

type FillSolidParams struct {
	color Color
}

func FillSolidRGB(color Color, duration uint16) []byte {
	return NewCommand(CMD_FILL_SOLID_RGB, duration).WithParams(
		FillSolidParams{color})
}

func Off(duration uint16) []byte {
	return NewCommand(CMD_OFF, duration).Bytes()
}

type FadeToBlackParams struct {
	fadeBy uint8
}

func FadeToBlack(fadeBy uint8, duration uint16) []byte {
	return NewCommand(CMD_FADE_TO_BLACK, duration).WithParams(
		FadeToBlackParams{fadeBy})
}

type ConfettiParams struct {
	count  uint8
	fadeBy uint8
}

func ConfettiDefault(duration uint16) []byte {
	return NewCommand(CMD_CONFETTI, duration).Bytes()
}

func Confetti(count uint8, fadeBy uint8, duration uint16) []byte {
	return NewCommand(CMD_CONFETTI, duration).WithParams(
		ConfettiParams{count, fadeBy})
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
