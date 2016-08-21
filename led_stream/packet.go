package led_stream

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"log"
)

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
