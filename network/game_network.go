package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"soulworker-server/global"
)

// GameNetwork - Container for the GameNetwork
type GameNetwork struct {
	Network
}

// NewGameNetwork - Creates a new GameNetwork instance
func NewGameNetwork() *GameNetwork {
	gameNetwork := &GameNetwork{
		Network{
			Name: "Game Network",
			Port: global.GameAuthPort,
		},
	}
	gameNetwork.dataHandler = gameNetwork.process
	return gameNetwork
}

// process - Processes data from the network
func (gameNetwork *GameNetwork) process(channel Connection, packetID uint16, buffer *bytes.Buffer) int {
	if buffer.Len() == 0 {
		return 0
	}

	if packetID == 0x0213 {
		// Client -> Server
		// ID=0x0213, Size=15, Total=22
		// 00000000  af b7 0f 00 02 00 28 ff  bb 00 00 00 00 00 00     |......(........|
		uuid := buffer.Next(4)

		// Server -> Client
		// ID=0x0214, Size=5, Total=12
		// 00000000  00 af b7 0f 00                                    |.....|
		buf := new(bytes.Buffer)
		buf.Write([]byte{0x00})
		buf.Write(uuid)

		channel.writeQueue <- global.Packet{ID: 0x0214, Data: buf}

		// Server -> Client
		// ID=0x0403, Size=22, Total=29
		// 00000000  2a 9e 1a 5d 00 00 00 00  e3 07 07 00 02 00 00 00  |*..]............|
		// 00000010  3a 00 22 00 00 00                                 |:."...|
		buf2 := new(bytes.Buffer)
		buf2.Write([]byte{0x2a, 0x9e, 0x1a, 0x5d, 0x00, 0x00, 0x00, 0x00,  0xe3, 0x07, 0x07, 0x00, 0x02, 0x00, 0x00, 0x00,
			0x3a, 0x00, 0x22, 0x00, 0x00, 0x00})

		channel.writeQueue <- global.Packet{ID: 0x0403, Data: buf2}
	} else if packetID == 0x0311 {
		// Client -> Server
		// ID=0x0311, Size=8, Total=15
		// 00000000  28 ff bb 00 00 00 00 00                           |(.......|

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
		buf.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

		channel.writeQueue <- global.Packet{ID: 0x0312, Data: buf}

		// Server -> Client
		// ID=0x0107, Size=14, Total=21
		// 00000000  01 00 01 00 01 01 00 00  01 00 00 01 00 01        |..............|
		buf2 := new(bytes.Buffer)
		buf2.Write([]byte{0x01, 0x00, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00,  0x00, 0x00, 0x00, 0x01, 0x00, 0x01})

		channel.writeQueue <- global.Packet{ID: 0x0107, Data: buf2}
	} else if packetID == 0x0347 {
		// Nothing atm
		// Client -> Server
		// ID=0x0347, Size=4, Total=11
		// 00000000  00 00 00 00                                       |....|
	} else if packetID == 0x0313 {
		// Select character
		// Client -> Server
		// ID=0x0313, Size=17, Total=24
		// 00000000  00 29 1e 00 00 00 00 00  01 00 00 00 00 00 00 00  |.)..............|
		// 00000010  00                                                |.|
		ucid := buffer.Next(4)
		uuid := []byte{0x00, 0x00, 0x00, 0x00}

		// Server -> Client
		// ID=0x0315, Size=92, Total=99
		// 00000000  00 29 1e 00 af b7 0f 00  02 02 02 00 7d 36 20 00  |.)..........}6 .|
		// 00000010  7d 36 20 00 cb 21 00 00 >77 52<02 00 00 00 00 00  |}6 ..!..wR......|
		// 00000020  00 00 00 00 0e 00 32 30  36 2e 32 35 33 2e 31 37  |......206.253.17|
		// 00000030  35 2e 38 32>0e 2b<ff ff  00 00 00 00 00 00 00 00  |5.82.+..........|
		// 00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
		// 00000050  00 00 00 00 00 00 00 00  00 00 00 00              |............|
		buf := new(bytes.Buffer)
		buf.Write(ucid)
		buf.Write(uuid)

		p := GetPortBytes(global.GameWorldPort) // 21111 in the packet above (0x5277), 10200 in the code

		buf.Write([]byte{
			0x01, 0x01, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x85, 0x00, 0x00, 0x03,  p[0], p[1], 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
		})

		gameWorldIP := global.ServerMap[0].GetIP()

		buf.Write([]byte{uint8(len(gameWorldIP)), 0x00})
		buf.Write([]byte(gameWorldIP))
		buf.Write(p) // 11022 in the packet above (0x2b0e), 10200 in the code
		buf.Write([]byte{
			                        0xff, 0xff, 0x00, 0x00,  0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,  0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,  0x00, 0x00,
		})

		channel.writeQueue <- global.Packet{ID: 0x0314, Data: buf}
	} else if packetID == 0x0301 {
		type CharacterInfo struct {
			UUID uint32

			UsernameLen uint16
			Username    []byte

			CharSelection uint16 // (1 - Haru, 2 - Erwin, 3 - Lily, 4 - Stella, 5 - Jin, 6 - Iris, 7 - Chii)

			UUID2 uint32

			// Maybe they actually should be uint16? :S
			HairStyle    uint16
			HairColor    uint16
			EyeColor     uint16
			SkinColor    uint16

			Unknownuint32Set1 [6]uint32
			UnknownByte1      byte
			Unknownuint32Set2 [127]uint32

			UnknownByte2		byte
			UnknownByte3		byte

			CharacterSlot		byte
			StandardOutfit		byte
			UnknownByte4		byte // level ?
			UnknownByte5		byte // xp ?
			UnknownByte6		byte // rank ?
		}

		ucid1 := make([]byte, 4)
		rand.Read(ucid1) // Generate a random ucid for now

		// somehow this gives erwin stuff? idk (SHOULD GIVE CHII NOW)
		ucid2 := []byte{0x25, 0x28, 0xad, 0x06}

		ucid3 := []byte{0x2d, 0xa9, 0x2c, 0x0d}

		ucid4 := []byte{0x21, 0xab, 0x2c, 0x0d} // ALWAYS THE SAME (PROBABLY COORDS FOR TUTORIAL)

		ucid5 := []byte{0x91, 0xa9, 0x2c, 0x0d} // ALWAYS THE SAME (PROBABLY COORDS FOR TUTORIAL)

		charInfo := &CharacterInfo{}
		// prepend 01
		binary.Read(buffer, binary.LittleEndian, &charInfo.UUID)

		binary.Read(buffer, binary.LittleEndian, &charInfo.UsernameLen)
		charInfo.Username = make([]byte, charInfo.UsernameLen)
		binary.Read(buffer, binary.LittleEndian, &charInfo.Username)

		binary.Read(buffer, binary.LittleEndian, &charInfo.CharSelection)

		binary.Read(buffer, binary.LittleEndian, &charInfo.HairStyle)
		binary.Read(buffer, binary.LittleEndian, &charInfo.HairColor)
		binary.Read(buffer, binary.LittleEndian, &charInfo.EyeColor)
		binary.Read(buffer, binary.LittleEndian, &charInfo.SkinColor)

		binary.Read(buffer, binary.LittleEndian, &charInfo.Unknownuint32Set1)
		binary.Read(buffer, binary.LittleEndian, &charInfo.UnknownByte1)
		binary.Read(buffer, binary.LittleEndian, &charInfo.Unknownuint32Set2)

		binary.Read(buffer, binary.LittleEndian, &charInfo.UnknownByte2)
		binary.Read(buffer, binary.LittleEndian, &charInfo.UnknownByte3)

		binary.Read(buffer, binary.LittleEndian, &charInfo.CharacterSlot)
		binary.Read(buffer, binary.LittleEndian, &charInfo.StandardOutfit)

		binary.Read(buffer, binary.LittleEndian, &charInfo.UnknownByte4)
		binary.Read(buffer, binary.LittleEndian, &charInfo.UnknownByte5)
		binary.Read(buffer, binary.LittleEndian, &charInfo.UnknownByte6)
		// replace last 4 bytes with uuid
		// append 00 00 00 00 00 00 00 00 00 00
		charInfo.UUID = binary.LittleEndian.Uint32(ucid1)
		charInfo.Unknownuint32Set1[4] = binary.LittleEndian.Uint32(ucid2)
		charInfo.Unknownuint32Set2[12] = binary.LittleEndian.Uint32(ucid3) // X
		charInfo.Unknownuint32Set2[18] = binary.LittleEndian.Uint32(ucid4) // Z
		charInfo.Unknownuint32Set2[22] = binary.LittleEndian.Uint32(ucid5) // Y ?

		//

		buf := new(bytes.Buffer)
		buf.Write([]byte{0x01})
		_ = binary.Write(buf, binary.LittleEndian, charInfo.UUID)

		_ = binary.Write(buf, binary.LittleEndian, charInfo.UsernameLen)
		_ = binary.Write(buf, binary.LittleEndian, charInfo.Username)

		_ = binary.Write(buf, binary.LittleEndian, charInfo.CharSelection)

		_ = binary.Write(buf, binary.LittleEndian, charInfo.HairStyle)
		_ = binary.Write(buf, binary.LittleEndian, charInfo.HairColor)
		_ = binary.Write(buf, binary.LittleEndian, charInfo.EyeColor)
		_ = binary.Write(buf, binary.LittleEndian, charInfo.SkinColor)

		_ = binary.Write(buf, binary.LittleEndian, charInfo.Unknownuint32Set1)
		_ = binary.Write(buf, binary.LittleEndian, charInfo.UnknownByte1)
		_ = binary.Write(buf, binary.LittleEndian, charInfo.Unknownuint32Set2)

		_ = binary.Write(buf, binary.LittleEndian, charInfo.UnknownByte2)
		_ = binary.Write(buf, binary.LittleEndian, charInfo.UnknownByte3)

		_ = binary.Write(buf, binary.LittleEndian, charInfo.CharacterSlot)

		_ = binary.Write(buf, binary.LittleEndian, charInfo.UUID)
		buf.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

		fmt.Printf("%#+v\n", charInfo)

		channel.writeQueue <- global.Packet{ID: 0x0312, Data: buf}

		// Server -> Client
		// ID=0x0107, Size=14, Total=21
		// 00000000  01 00 01 00 01 01 00 00  01 00 00 01 00 01        |..............|
		buf2 := new(bytes.Buffer)
		buf2.Write([]byte{0x01, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01})

		channel.writeQueue <- global.Packet{ID: 0x0107, Data: buf2}
	}

	return 0
}
