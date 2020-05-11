package network

import (
	"../database"
	"../global"
	"./util"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
)

// This file contains the contents for the Login Network

// LoginNetwork - Container for the LoginNetwork
type LoginNetwork struct {
	Network
}

// NewLoginNetwork - Creates a new LoginNetwork instance
func NewLoginNetwork() *LoginNetwork {
	network := &LoginNetwork{
		Network{
			Name: "Login Network",
			Port: global.LoginPort,
		},
	}
	network.dataHandler = network.process
	return network
}

// process - Processes data from the network from the given connection, with the given packetId and packet buffer contents
// It returns an abstract integer value
func (network *LoginNetwork) process(channel *Connection, packetID uint16, buffer *bytes.Buffer) int {
	if buffer.Len() == 0 {
		return 0
	}

	if packetID == 0x2002 {
		// Basically just a hello message
		// Client -> Server
		// ID=0x2002, Size=4, Total=11
		// 00000000  00 00 00 00                                       |....|
	} else if packetID == 0x0201 /* KR */ || packetID == 0x0218 /* GF */ { // Login Auth Request
		// Login auth request (KRSW)
		// Client -> Server
		// ID=0x0201, Size=78, Total=85
		// 00000000  0e 00 58 00 58 00 58 00  58 00 58 00 58 00 58 00  |..X.X.X.X.X.X.X.|
		// 00000010  14 00 58 00 58 00 58 00  58 00 58 00 58 00 58 00  |..X.X.X.X.X.X.X.|
		// 00000020  58 00 58 00 58 00 22 00  58 00 58 00 2d 00 58 00  |X.X.X.".X.X.-.X.|
		// 00000030  58 00 2d 00 58 00 58 00  2d 00 58 00 58 00 2d 00  |X.-.X.X.-.X.X.-.|
		// 00000040  58 00 58 00 2d 00 58 00  58 00 83 96 98 00        |X.X.-.X.X.....|

		username := util.ReadStringUTF16(buffer)
		var password string
		var mac string
		if packetID == 0x0201 {
			password = util.ReadStringUTF16(buffer)
			mac = util.ReadStringUTF16(buffer)
		} else {
			mac = util.ReadStringUTF8(buffer)
		}

		leftover := buffer.Next(4)

		// Valid credentials => austin:coolman83
		fmt.Printf("[+] Received Login Request:\n\tUsername: %s\n\tPassword: %s\n\tMac Address: %s\n\tLeftovers: %+#v\n\n", username, password, mac, leftover)

		// Query db for username:password combo and if successful, continue
		if database.CanConnect() != nil {
			log.Println("[!] Could not connect to the database")
			return 0
		}

		accountId := database.VerifyLoginCredentials(username, password)

		errorCode := 0
		if accountId == 0 {
			errorCode = 1
		}

		channel.accountId = uint32(accountId)

		// Login auth response
		// 00000000  a2 c9 00 00 01 58 58 2d  58 58 2d 58 58 2d 58 58  |.....58-58-58-58|
		// 00000010  2d 58 58 2d 58 58 00 00  00 00 00 00 00 01 0e 00  |-58-58..........|
		// 00000020  58 00 58 00 58 00 58 00  58 00 58 00 58 00 53 9f  |X.X.X.X.X.X.X.S.|
		// 00000030  2a 00 00 00 00 00 00 00  00 00 00 00              |*...........|
		loginRes := new(bytes.Buffer)
		_ = binary.Write(loginRes, binary.LittleEndian, uint32(accountId)) // 4 bytes
		_ = binary.Write(loginRes, binary.LittleEndian, byte(1)) // 0x01
		util.WriteStringUTF8NoLength(loginRes, mac)
		util.WriteStringUTF16(loginRes, "") // errorMessage
		_ = binary.Write(loginRes, binary.LittleEndian, uint32(errorCode)) // 1 if bad login, 0 otherwise
		_ = binary.Write(loginRes, binary.LittleEndian, byte(0)) // Unknown
		util.WriteStringUTF16(loginRes, "") // unknown
		if accountId > 0 { // Session Key to keep track of a user between servers
			sessionKey := uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
			database.UpdateSessionKey(accountId, sessionKey)
			fmt.Println("Updating session key for id", accountId,":", sessionKey)
			_ = binary.Write(loginRes, binary.LittleEndian, sessionKey)
		} else {
			_ = binary.Write(loginRes, binary.LittleEndian, uint64(0))
		}
		_ = binary.Write(loginRes, binary.LittleEndian, byte(0)) // Unknown
		_ = binary.Write(loginRes, binary.LittleEndian, uint16(0)) // Unknown
		_ = binary.Write(loginRes, binary.LittleEndian, byte(0)) // Unknown
		_ = binary.Write(loginRes, binary.LittleEndian, byte(0)) // Unknown
		_ = binary.Write(loginRes, binary.LittleEndian, byte(0)) // Unknown

		channel.writeQueue <- global.Packet{ID: 0x0202, Data: loginRes}

		if accountId == 0 {
			return 0
		}

		// Server Options
		//serverOptions := new(bytes.Buffer)
		//serverOptions.Write(make([]byte, 64)) // 64 bytes of zero
		//for i := 0; i < 14; i++ {
		//	serverOptions.Write([]byte{0, 0})
		//}

		//channel.writeQueue <- global.Packet{ID: 0x0231, Data: serverOptions}


	} else if packetID == 0x0203 {
		// Server list request, 0xaf 0xb7 0x0f 0x00 seems like an identifier, appears later
		// Client -> Server
		// ID=0x0203, Size=4, Total=11
		// 00000000  af b7 0f 00                                       |....|

		// Server -> Client
		// ID=0x0204, Size=44, Total=51
		// 00000000  00 01 01 00 74 27 0d 00  4c 6f 73 74 20 4d 65 6d  |....t'..Lost Mem|
		// 00000010  6f 72 69 65 73 0c 00 31  32 2e 33 34 2e 35 36 2e  |ories..12.34.56.|
		// 00000020  37 38 39 01 00 00 00 45  01 00 00 03              |789....E....|
		serverList := new(bytes.Buffer)
		numServers := len(global.ServerMap)
		serverList.Write([]byte{0x00, byte(numServers)})
		for i := 0; i < numServers; i++ {
			server := global.ServerMap[i]
			serverList.Write([]byte{byte(i + 1), 0x00})
			serverList.Write(GetPortBytes(global.GameAuthPort))
			serverList.Write([]byte{byte(len(server.GetName())), 0x00})
			serverList.Write([]byte(server.GetName()))
			serverList.Write([]byte{byte(len(server.GetIP())), 0x00})
			serverList.Write([]byte(server.GetIP()))
			serverList.Write(
				[]byte{
					0x01, 0x00, 0x00, 0x00,
					0x00, 0x00, /* Number of people on the server, little endian, for population indicator */
					0x00, 0x00,
					byte(database.FetchUserCharacterCount(int(channel.accountId))), /* Number of characters the user has on the server */
				},
			)
		}

		channel.writeQueue <- global.Packet{ID: 0x0204, Data: serverList}

		// Server -> Client
		// ID=0x0231, Size=78, Total=85
		// 00000000  36 31 30 31 30 30 30 31  30 30 30 31 31 31 31 30  |6101000100011110|
		// 00000010  30 35 31 30 31 30 30 30  30 31 31 31 30 20 20 20  |0510100001110   |
		// 00000020  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
		// 00000030  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 00  |               .|
		// 00000040  01 01 01 01 01 01 01 01  01 01 01 01 01 01        |..............|
		outBuf := bytes.NewBuffer([]byte{
			0x36, 0x31, 0x30, 0x31, 0x30, 0x30, 0x30, 0x31,  0x30, 0x30, 0x30, 0x31, 0x31, 0x31, 0x31, 0x30,
			0x30, 0x35, 0x31, 0x30, 0x31, 0x30, 0x30, 0x30,  0x30, 0x31, 0x31, 0x31, 0x20, 0x20, 0x20, 0x20,
			0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,  0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
			0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,  0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x00,
			0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,  0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		})
		// 0x0231 = receive_eSUB_CMD_OPTION_LOAD ??? Don't know what any of this is
		channel.writeQueue <- global.Packet{ID: 0x0231, Data: outBuf}
	} else if packetID == 0x0205 {
		// Join server request
		// Client -> Server
		// ID=0x0205, Size=2, Total=9
		// 00000000  01 00                                             |..|
		selectedServer, err := buffer.ReadByte()
		if err != nil {
			panic(err)
		}
		if int(selectedServer) > len(global.ServerMap) {
			panic("Invalid server choice")
		}
		selectedServerIP := global.ServerMap[selectedServer-1].GetIP()

		// Server -> Client
		// ID=0x0211, Size=18, Total=25
		// 00000000  0e 00 31 32 33 2e 34 35  36 2e 37 38 39 2e 30 32  |..123.456.789.02|
		// 00000010  76 27                                             |v'|
		serverSelect := new(bytes.Buffer)
		serverSelect.Write([]byte{byte(len(selectedServerIP)), 0x00})
		serverSelect.Write([]byte(selectedServerIP))
		serverSelect.Write(GetPortBytes(global.GameAuthPort))

		channel.writeQueue <- global.Packet{ID: 0x0211, Data: serverSelect}
	}

	return 0
}
