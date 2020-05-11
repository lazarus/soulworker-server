package network

import (
	"../database"
	"../global"
	"./structures"
	"./util"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
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

// process - Processes data from the network from the given connection, with the given packetId and packet buffer contents
// It returns an abstract integer value
func (gameNetwork *GameNetwork) process(channel *Connection, packetID uint16, buffer *bytes.Buffer) int {
	if buffer.Len() == 0 {
		return 0
	}

	if packetID == 0x0213 {
		// Client -> Server
		// ID=0x0213, Size=15, Total=22
		// 00000000  af b7 0f 00 02 00 28 ff  bb 00 00 00 00 00 00     |......(........|
		var ClientEnterServerRequest struct {
			accountId uint32
			unknown0 uint16
			sessionKey uint64
			unknown1 byte
		}

		_ = binary.Read(buffer, binary.LittleEndian, &channel.accountId)
		_ = binary.Read(buffer, binary.LittleEndian, &ClientEnterServerRequest.unknown0)
		_ = binary.Read(buffer, binary.LittleEndian, &ClientEnterServerRequest.sessionKey)
		_ = binary.Read(buffer, binary.LittleEndian, &ClientEnterServerRequest.unknown1)

		if id := database.VerifySessionKey(channel.accountId, ClientEnterServerRequest.sessionKey); id == 0 {
			log.Fatal("Invalid sessionKey", ClientEnterServerRequest.sessionKey, "for account", channel.accountId)
		} else {
			fmt.Println("Session key is okay!")
			log.Println("session key before reconnect", ClientEnterServerRequest.sessionKey)
		}

		//spew.Dump(ClientEnterServerRequest)

		//fmt.Println("HELLOacct id", channel.accountId)
		//fmt.Printf("%#+v\n", channel)
		//fmt.Println("sssssss")

		// Server -> Client
		// ID=0x0214, Size=5, Total=12
		// 00000000  00 af b7 0f 00                                    |.....|
		buf := new(bytes.Buffer)
		buf.Write([]byte{0x00})
		_ = binary.Write(buf, binary.LittleEndian, ClientEnterServerRequest.accountId)

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
		buf2.Write([]byte{0x01, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00})

		channel.writeQueue <- global.Packet{ID: 0x0107, Data: buf2}
	} else if packetID == 0x0347 {
		// Nothing atm
		// Client -> Server
		// ID=0x0347, Size=4, Total=11
		// 00000000  01 00 00 00                                       |....|
	} else if packetID == 0x0313 {
		// Select character
		// Client -> Server
		// ID=0x0313, Size=17, Total=24
		//00000000  02 8f 31 00 00 00 00 00  01 00 00 00 00 00 00 00  |..1.............|
		//00000010  00                                                |.|

		var ClientSelectCharacterReq struct {
			characterId uint32
			unknown0 uint32
			unknown1 byte
			unknown2 uint32
			unknown3 uint32
		}

		_ = binary.Read(buffer, binary.LittleEndian, &ClientSelectCharacterReq.characterId)
		_ = binary.Read(buffer, binary.LittleEndian, &ClientSelectCharacterReq.unknown0)
		_ = binary.Read(buffer, binary.LittleEndian, &ClientSelectCharacterReq.unknown1)
		_ = binary.Read(buffer, binary.LittleEndian, &ClientSelectCharacterReq.unknown2)
		_ = binary.Read(buffer, binary.LittleEndian, &ClientSelectCharacterReq.unknown3)

		// Server -> Client
		// ID=0x0315, Size=92, Total=99
		//00000000  02 8f 31 00 91 d2 0e 00  02 02 02 00 45 37 20 00  |..1.........E7 .|
		//00000010  45 37 20 00 fd 79 00 00  79 52 02 00 00 00 00 00  |E7 ..y..yR......|
		//00000020  00 00 00 00 0e 00 32 30  36 2e 32 35 33 2e 31 37  |......206.253.17|
		//00000030  35 2e 38 32 0e 2b ff ff  00 00 00 00 00 00 00 00  |5.82.+..........|
		//00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
		//00000050  00 00 00 00 00 00 00 00  00 00 00 00              |............|
		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.LittleEndian, &ClientSelectCharacterReq.characterId)
		_ = binary.Write(buf, binary.LittleEndian, &channel.accountId)
		_ = binary.Write(buf, binary.LittleEndian, uint32(0))
		_ = binary.Write(buf, binary.LittleEndian, uint32(0))
		_ = binary.Write(buf, binary.LittleEndian, uint32(0))
		_ = binary.Write(buf, binary.LittleEndian, uint32(0))
		_ = binary.Write(buf, binary.LittleEndian, global.GameWorldPort)
		_ = binary.Write(buf, binary.LittleEndian, make([]byte, 10))

		util.WriteStringUTF8NoTrailing(buf, global.ServerMap[0].GetIP()) // Errors with a trailing null byte
		_ = binary.Write(buf, binary.LittleEndian, global.GameWorldPort)

		position := structures.WorldPosition{
			MapId: 10003,
			X: 10444.9951,
			Y: 10179.7461,
			Z: 100.325394,
			Orientation: 100.0,
		}
		position.Write(buf)

		_ = binary.Write(buf, binary.LittleEndian, byte(0))

		_ = binary.Write(buf, binary.LittleEndian, byte(0))
		_ = binary.Write(buf, binary.LittleEndian, byte(0))
		_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // error code?
		_ = binary.Write(buf, binary.LittleEndian, byte(0))
		_ = binary.Write(buf, binary.LittleEndian, uint32(0))

		channel.writeQueue <- global.Packet{ID: 0x0314, Data: buf}
	} else if packetID == 0x0301 {

		charInfo := structures.CharacterInfo{}
		charInfo.Read(buffer)

		//fmt.Printf("%#+v\n", charInfo)

		//spew.Dump(charInfo)

		//

		//char := &structures.CharacterModel{
		//	AccountId:  1337,
		//	Index:      charInfo.Index,
		//	Name:       charInfo.Username,
		//	Class:      charInfo.CharSelection,
		//	Level:      charInfo.Level,
		//	Appearance: charInfo.Appearance,
		//
		//	MapId:      10003,
		//	X:          10444.9951,
		//	Y:          10179.7461,
		//	Z:          100.325394,
		//}
		//
		//
		//charList := []structures.CharacterInfo {
		//	{
		//		Id:            3001,//char.Id,
		//		Username:      char.Name,
		//		CharSelection: char.Class,
		//		Appearance:    char.Appearance,
		//		Level:         char.Level,
		//		Index:         char.Index,
		//	},
		//}

		//fmt.Println("acc id", channel.accountId)
		//fmt.Printf("%#+v\n", channel)
		//fmt.Println("acc id", &channel.accountId)


		charInfo.AccountId = channel.accountId
		charId := database.InsertCharacterToDb(&charInfo)
		_ /* charCount */ = byte(database.FetchUserCharacterCount(int(channel.accountId)))

		charInfo.Id = uint32(charId)
		//fmt.Println("Newly inserted character id is", charInfo.Id)

		buf := new(bytes.Buffer)
		//buf.Write([]byte{charCount}) // Character Count
		buf.Write([]byte{1}) // Character Count

		charInfo.Unknown6 = 110011301 + (1000000 * uint32(charInfo.CharSelection))

		//for _, char := range charList {
		//	char.Write(buf)
		//}


		spew.Dump(charInfo.AccountId)

		charInfo.Write(buf)

		_ = binary.Write(buf, binary.LittleEndian, charInfo.Id)//char.Id)
		_ = binary.Write(buf, binary.LittleEndian, byte(0))
		_ = binary.Write(buf, binary.LittleEndian, byte(0))
		_ = binary.Write(buf, binary.LittleEndian, uint64(0))
		_ = binary.Write(buf, binary.LittleEndian, byte(0))
		_ = binary.Write(buf, binary.LittleEndian, uint64(0))

		channel.writeQueue <- global.Packet{ID: 0x0312, Data: buf}

		// Server -> Client
		// ID=0x0107, Size=14, Total=21
		// 00000000  01 00 01 00 01 01 00 00  01 00 00 01 00 01        |..............|
		buf2 := new(bytes.Buffer)
		buf2.Write([]byte{0x01, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00})

		channel.writeQueue <- global.Packet{ID: 0x0107, Data: buf2}
	}

	return 0
}
