package network

import (
	"fmt"
	"log"
	"time"

	"../database"
	"../global"
	. "../network/packets"
	"github.com/davecgh/go-spew/spew"
)

// GameNetwork - Container for the GameNetwork
type GameNetwork struct {
	Network
}

// NewGameNetwork - Creates a new GameNetwork instance
func NewGameNetwork() *GameNetwork {
	gameNetwork := &GameNetwork{
		Network{
			Name: "Char Network",
			Port: global.GameAuthPort,
		},
	}
	gameNetwork.dataHandler = gameNetwork.process
	return gameNetwork
}

// process - Processes data from the network from the given connection, with the given packetId and packet buffer contents
// It returns an abstract integer value
func (gameNetwork *GameNetwork) process(channel *Connection, packetID PacketType, packet interface{}) int {

	switch packetID {
	case Character_EnterCharServerRequest:
		enterGameServerRequest := packet.(*EnterCharServerRequest)
		channel.accountId = enterGameServerRequest.AccountId

		if id := database.VerifySessionKey(channel.accountId, enterGameServerRequest.SessionKey); id == 0 {
			log.Fatal("Invalid sessionKey", enterGameServerRequest.SessionKey, "for account", channel.accountId)
		} else {
			fmt.Println("Session key is okay!")
			log.Println("session key before reconnect", enterGameServerRequest.SessionKey)
		}

		channel.writeQueue <- &EnterCharServerResponse{
			AccountId: enterGameServerRequest.AccountId,
		}

		now := time.Now()
		channel.writeQueue <- &ServerWorldCurDate{
			Timestamp: uint64(now.Unix()),
			Year:      uint16(now.Year()),
			Month:     uint16(now.Month()),
			Day:       uint16(now.Day()),
			Hour:      uint16(now.Hour()),
			Minute:    uint16(now.Minute()),
			Second:    uint16(now.Second()),
		}
		break
	case Character_CharacterListRequest:
		// Fetch characters from db and pass to CharacterListResponse
		channel.writeQueue <- &CharacterListResponse{}

		// Server -> Client
		// ID=0x0107, Size=14, Total=21
		// 00000000  01 00 01 00 01 01 00 00  01 00 00 01 00 01        |..............|
		//buf2 := new(bytes.Buffer)
		//buf2.Write([]byte{0x01, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00})
		//
		//channel.writeQueue <- global.Packet{ID: 0x0107, Data: buf2}

		break
	case 0x0347:
		// Nothing atm
		// Client -> Server
		// ID=0x0347, Size=4, Total=11
		// 00000000  01 00 00 00                                       |....|

		break
	case Character_SelectCharacterRequest:
		selectCharacterRequest := packet.(*SelectCharacterRequest)

		channel.writeQueue <- &SelectCharacterResponse{
			CharacterId: selectCharacterRequest.CharacterId,
			AccountId:   channel.accountId,
		}

		break
	case Character_CreateCharacterRequest:
		createCharacterRequest := packet.(*CreateCharacterRequest)

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

		createCharacterRequest.CharInfo.AccountId = channel.accountId
		spew.Dump(createCharacterRequest.CharInfo.AccountId)

		charId := database.InsertCharacterToDb(&createCharacterRequest.CharInfo)
		_ /* charCount */ = byte(database.FetchUserCharacterCount(channel.accountId))

		createCharacterRequest.CharInfo.Id = uint32(charId)
		//fmt.Println("Newly inserted character id is", charInfo.Id)

		createCharacterRequest.CharInfo.Unknown6 = 110011301 + (1000000 * uint32(createCharacterRequest.CharInfo.CharSelection)) // Weapon

		channel.writeQueue <- &CharacterListResponse{CharInfo: createCharacterRequest.CharInfo}

		// Server -> Client
		// ID=0x0107, Size=14, Total=21
		// 00000000  01 00 01 00 01 01 00 00  01 00 00 01 00 01        |..............|
		//buf2 := new(bytes.Buffer)
		//buf2.Write([]byte{0x01, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00})
		//
		//channel.writeQueue <- global.Packet{ID: 0x0107, Data: buf2}
		break
	}

	return 0
}
