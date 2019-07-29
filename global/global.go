package global

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// KeyTable - Array of bytes for XORing data
// ServerMap - Array of Servers
var (
	KeyTable = []byte{
		0x57, 0x19, 0xC6, 0x2D, 0x56, 0x68, 0x3A, 0xCC,
		0x60, 0x3B, 0x0B, 0xB1, 0x90, 0x5C, 0x4A, 0xF8,

		0x80, 0x28, 0xB1, 0x45, 0xB6, 0x85, 0xE7, 0x4C,
		0x06, 0x2D, 0x55, 0x83, 0xAF, 0x44, 0x99, 0x95,

		0xD9, 0x98, 0xBF, 0xAE, 0x53, 0x43, 0x63, 0xC8,
		0x4A, 0x71, 0x80, 0x9D, 0x0B, 0xA1, 0x70, 0x8A,

		0x0F, 0x54, 0x9C, 0x1B, 0x06, 0xC0, 0xEA, 0x3C,
		0xC0, 0x88, 0x71, 0x48, 0xB3, 0xB9, 0x45, 0x78,
	}
	ServerMap     []Server
	LoginPort     uint16 = 10000
	GameAuthPort  uint16 = 10100
	GameWorldPort uint16 = 10200
)

// Server - Server with a name and IP
type Server struct {
	Name string
	IP   string
}

// Log - Logs a message with a title
func Log(params ...string) {
	if len(params) == 1 {
		fmt.Printf("[%s]\t%s\n", "SoulWorker", params[0])
	} else {
		fmt.Printf("[%s]\t%s\n", params[0], params[1])
	}
}

// GetName - Gets the name of a Game Server
func (server Server) GetName() string {
	return server.Name
}

// GetIP - Gets the IP of a Game Server
func (server Server) GetIP() string {
	return server.IP
}

// Packet - Game Packet with an ID and Data
type Packet struct {
	ID   uint16
	Data *bytes.Buffer
}

// Encrypt - Encodes a packet to bytes
func (packet Packet) Encrypt() []byte {

	buffer := new(bytes.Buffer)

	var header struct {
		Magic  uint16 // 0x02 0x00
		Size   uint16 // 0xXX 0xXX
		Sender uint8  // 0x01
	}
	header.Magic = 0x02
	header.Size = uint16(packet.Data.Len() + 7) // Data + 5 byte Header + 2 byte ID
	header.Sender = 0x01

	binary.Write(buffer, binary.LittleEndian, header.Magic)
	binary.Write(buffer, binary.LittleEndian, header.Size)
	binary.Write(buffer, binary.LittleEndian, header.Sender)

	var pac = make([]byte, packet.Data.Len()+2)
	copy(pac[0:2], []byte{uint8(packet.ID>>8) & 0xFF, uint8(packet.ID & 0xFF)})
	copy(pac[2:], packet.Data.Bytes())

	for i := 0; i < len(pac); i++ {
		var byte1 = pac[i]
		var byte2 = KeyTable[4*int(header.Magic)-3*(i/3)+i]
		// fmt.Printf("%02X ^ %02X = %02X\n", byte1, byte2, byte1 ^ byte2)
		buffer.WriteByte(byte1 ^ byte2)
	}

	return buffer.Bytes()
}
