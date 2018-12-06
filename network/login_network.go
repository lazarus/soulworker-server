package network

import (
	"bytes"
	"../global"

	// "fmt"

	"encoding/binary"
)

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

// process - Processes data from the network
func (network *LoginNetwork) process(channel Connection, packetID uint16, buffer *bytes.Buffer) int {
	if buffer.Len() == 0 {
		return 0
	}

	if packetID == 0x2002 {
		// Basically just a hello message
		// Client -> Server
		// ID=0x2002, Size=4, Total=11
		// 00000000  00 00 00 00                                       |....|
	} else if packetID == 0x0218 {
		// Login auth request
		// Client -> Server
		// ID=0x0218, Size=110, Total=117
		authCodeLength := int(binary.LittleEndian.Uint16(buffer.Next(2)))
		/*authCode :=*/ buffer.Next(authCodeLength)

		macAddrLength := int(binary.LittleEndian.Uint16(buffer.Next(2)))
		macAddr := buffer.Next(macAddrLength)

		// Server -> Client
		// ID=0x0202, Size=75, Total=82
		// 00000000  af b7 0f 00 01 37 43 2d  36 37 2d 41 32 2d 39 34  |.....xx-xx-xx-xx|
		// 00000010  2d 38 45 2d 42 45 00 00  00 00 00 00 00 01 1e 00  |-xx-xx..........|
		// 00000020  6d 00 6f 00 72 00 6e 00  69 00 6e 00 67 00 66 00  |x.x.x.x.x.x.x.x.|
		// 00000030  69 00 72 00 65 00 39 00  32 00 39 00 33 00 28 ff  |x.x.x.x.x.x.x.(.|
		// 00000040  bb 00 00 00 00 00 00 00  00 00 00                 |...........|
		loginRes := new(bytes.Buffer)
		loginRes.Write([]byte{0xaf, 0xb7, 0x0f, 0x00, 0x01})
		loginRes.Write(macAddr)
		loginRes.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})

		username := []byte("Austin")
		usernameLength := byte(len(username) * 2)

		loginRes.Write([]byte{usernameLength, 0x00})

		for i := 0; i < len(username); i++ {
			loginRes.Write([]byte{username[i], 0x00})
		}

		loginRes.Write([]byte{0x28, 0xff, 0xbb, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

		channel.writeQueue <- global.Packet{ID: 0x0202, Data: loginRes}
	} else if packetID == 0x0203 {
		// Server list request, 0xaf 0xb7 0x0f 0x00 seems like an identifier, appears later
		// Client -> Server
		// ID=0x0203, Size=4, Total=11
		// 00000000  af b7 0f 00                                       |....|

		// Server -> Client
		// ID=0x0204, Size=300, Total=307
		// 00000000  00 07 01 00 75 27 0d 00  5b 44 45 5d 20 4c 61 6b  |....u'..[DE] Lak|
		// 00000010  61 73 68 69 6e 0e 00 31  39 34 2e 31 38 37 2e 31  |ashin..194.187.1|
		// 00000020  39 2e 31 33 30 01 00 00  00 8f 01 00 00 00 02 00  |9.130...........|
		// 00000030  76 27 0d 00 5b 4e 41 5d  20 54 65 6e 65 62 72 69  |v'..[NA] Tenebri|
		// 00000040  73 0e 00 32 30 36 2e 32  35 33 2e 31 37 33 2e 36  |s..206.253.173.6|
		// 00000050  32 01 00 00 00 89 01 00  00 00 03 00 77 27 0b 00  |2...........w'..|
		// 00000060  5b 45 4e 5d 20 43 61 6e  64 75 73 0e 00 31 39 34  |[EN] Candus..194|
		// 00000070  2e 31 38 37 2e 31 39 2e  31 33 30 01 00 00 00 d3  |.187.19.130.....|
		// 00000080  02 00 00 00 04 00 78 27  0a 00 5b 46 52 5d 20 52  |......x'..[FR] R|
		// 00000090  75 63 63 6f 0e 00 31 39  34 2e 31 38 37 2e 31 39  |ucco..194.187.19|
		// 000000a0  2e 31 33 30 01 00 00 00  51 01 00 00 00 05 00 79  |.130....Q......y|
		// 000000b0  27 0a 00 5b 50 4c 5d 20  47 72 61 63 65 0e 00 31  |'..[PL] Grace..1|
		// 000000c0  39 34 2e 31 38 37 2e 31  39 2e 31 33 30 01 00 00  |94.187.19.130...|
		// 000000d0  00 e3 00 00 00 00 06 00  7a 27 0e 00 5b 45 53 5d  |........z'..[ES]|
		// 000000e0  20 41 6d 61 72 79 6c 6c  69 73 0e 00 31 39 34 2e  | Amaryllis..194.|
		// 000000f0  31 38 37 2e 31 39 2e 31  33 30 01 00 00 00 bf 00  |187.19.130......|
		// 00000100  00 00 00 07 00 7b 27 0a  00 5b 49 54 5d 20 55 72  |.....{'..[IT] Ur|
		// 00000110  69 65 6c 0e 00 31 39 34  2e 31 38 37 2e 31 39 2e  |iel..194.187.19.|
		// 00000120  31 33 30 01 00 00 00 87  00 00 00 00              |130.........|
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
					0x00, 0x00, /* Number of people on the server, little endian */
					0x00, 0x00,
					0x00, /* Number of characers the user has on the server */
				},
			)
		}

		channel.writeQueue <- global.Packet{ID: 0x0204, Data: serverList}

		// Server -> Client
		// ID=0x0231, Size=78, Total=85
		// 00000000  35 31 30 31 30 30 30 31  30 31 30 31 31 31 31 31  |5101000101011111|
		// 00000010  31 35 31 30 31 30 30 30  30 31 30 31 20 20 20 20  |151010000101    |
		// 00000020  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 20  |                |
		// 00000030  20 20 20 20 20 20 20 20  20 20 20 20 20 20 20 00  |               .|
		// 00000040  00 00 00 00 00 01 00 01  01 01 01 01 01 01        |..............|
		outBuf := bytes.NewBuffer([]byte{
			0x35, 0x31, 0x30, 0x31, 0x30, 0x30, 0x30, 0x31, 0x30, 0x31, 0x30, 0x31, 0x31, 0x31, 0x31, 0x31,
			0x31, 0x35, 0x31, 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x31, 0x20, 0x20, 0x20, 0x20,
			0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
			0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		})
		// 0x0231 = receive_eSUB_CMD_OPTION_LOAD ??? Don't know what any of this is
		channel.writeQueue <- global.Packet{ID: 0x0231, Data: outBuf}
	} else if packetID == 0x0205 {
		// Join server request
		// Client -> Server
		// ID=0x0205, Size=2, Total=9
		// 00000000  02 00                                             |..|
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
		// 00000000  0e 00 32 30 36 2e 32 35  33 2e 31 37 33 2e 36 32  |..xxx.xxx.xxx.xx|
		// 00000010  76 27                                             |v'|
		serverSelect := new(bytes.Buffer)
		serverSelect.Write([]byte{byte(len(selectedServerIP)), 0x00})
		serverSelect.Write([]byte(selectedServerIP))
		serverSelect.Write(GetPortBytes(global.GameAuthPort))

		channel.writeQueue <- global.Packet{ID: 0x0211, Data: serverSelect}
	}

	return 0
}
