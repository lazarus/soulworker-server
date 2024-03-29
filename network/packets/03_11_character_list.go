package packets

import (
	"bytes"
	"encoding/binary"

	. "soulworker-server/network/structures"
)

type CharacterListRequest struct {
}

func (c *CharacterListRequest) unmarshal(buffer *bytes.Buffer) error {
	// Client -> Server
	// ID=0x0311, Size=8, Total=15
	// 00000000  28 ff bb 00 00 00 00 00                           |(.......|

	return nil
}

type CharacterListResponse struct {
	CharInfo CharacterInfo
}

func (c *CharacterListResponse) marshal() ([]byte, error) {
	// Commented data is what is sent when the user has 0 characters
	// Server -> Client
	// ID=0x0312, Size=15, Total=22
	// 00000000  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00     |...............|
	buf := new(bytes.Buffer)

	// charName := "Austin"
	// buf.Write([]byte{0x03, 0x37, 0x71, 0x01, 0x00})
	// util.WriteString2(buf, charName)
	// buf.Write([]byte{
	// 	0x03, 0x00, 0x17, 0x05,
	// 	0xFF, 0x08, 0xB5, 0x14, 0xE6, 0x0C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, /* Character Level */ 55, 0x00,
	// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x73, 0x2D, 0xC3, 0x06, 0x00, 0xFF,
	// 	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF5, 0x8C, 0xC6, 0x0D, 0x15,
	// 	0x65, 0xCE, 0x0D, 0xFD, 0xB3, 0xCE, 0x0D, 0xC1, 0x08, 0xCF, 0x0D, 0xA1, 0x32, 0xD4, 0x0D, 0xFF,
	// 	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	// 	0xFF, 0xFF, 0xFF, 0xAD, 0x3F, 0xC5, 0x0D, 0x7F, 0x87, 0xE4, 0x0D, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	// 	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x0F, 0x89, 0xE4, 0x0D, 0xA1, 0x41, 0xC5, 0x0D, 0xF3,
	// 	0x14, 0xED, 0x0D, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x11, 0x40, 0xC5, 0x0D, 0xE3,
	// 	0x87, 0xE4, 0x0D, 0x95, 0x12, 0xC9, 0x0D, 0x75, 0xC4, 0xC8, 0x0D, 0x00, 0x00, 0x00, 0x00, 0x00,
	// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x3F, 0x00, 0x00, 0x80,
	// 	0x3F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0xDA, 0x1B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	// 	0x00, 0x00, 0x00, 0x00,
	// })
	if c.CharInfo.AccountId == 0 {
		buf.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	} else {
		//buf.Write([]byte{charCount}) // Character Count
		buf.Write([]byte{1}) // Character Count

		//for _, char := range charList {
		//	char.Write(buf)
		//}

		c.CharInfo.Write(buf)

		_ = binary.Write(buf, binary.LittleEndian, c.CharInfo.Id) //char.Id)
		_ = binary.Write(buf, binary.LittleEndian, byte(0))
		_ = binary.Write(buf, binary.LittleEndian, byte(0))
		_ = binary.Write(buf, binary.LittleEndian, uint64(0))
		_ = binary.Write(buf, binary.LittleEndian, byte(0))
		_ = binary.Write(buf, binary.LittleEndian, uint64(0))
	}

	return buf.Bytes(), nil
}

func (c *CharacterListResponse) id() PacketType {
	return Character_CharacterListResponse
}
