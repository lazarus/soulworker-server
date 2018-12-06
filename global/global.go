package global

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// KeyTable - Array of bytes for XORing data
// ServerMap - Array of Servers
var (
	KeyTable      []byte
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

	var bytes = make([]byte, packet.Data.Len()+2)
	copy(bytes[0:2], []byte{uint8(packet.ID>>8) & 0xFF, uint8(packet.ID & 0xFF)})
	copy(bytes[2:], packet.Data.Bytes())

	for i := 0; i < len(bytes); i++ {
		var byte1 = bytes[i]
		var byte2 = KeyTable[4*int(header.Magic)-3*(i/3)+i]
		// fmt.Printf("%02X ^ %02X = %02X\n", byte1, byte2, byte1 ^ byte2)
		buffer.WriteByte(byte1 ^ byte2)
	}

	return buffer.Bytes()
}
