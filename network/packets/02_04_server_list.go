package packets

import (
	"bytes"
	"encoding/binary"
	"errors"

	"../../database"
	. "../../global"
	. "../util"
)

type ServerListRequest struct {
	AccountId uint32
}

func (p *ServerListRequest) unmarshal(buffer *bytes.Buffer) error {
	// Server list request, 0xaf 0xb7 0x0f 0x00 seems like an identifier, appears later
	// Client -> Server
	// ID=0x0203, Size=4, Total=11
	// 00000000  af b7 0f 00                                       |....|

	return binary.Read(buffer, binary.LittleEndian, &p.AccountId)
}

type ServerListResponse struct {
	AccountId uint32
}

func (p *ServerListResponse) id() PacketType {
	return Login_ServerListResponse
}

func (p *ServerListResponse) marshal() ([]byte, error) {
	// Server -> Client
	// ID=0x0204, Size=44, Total=51
	// 00000000  00 01 01 00 74 27 0d 00  4c 6f 73 74 20 4d 65 6d  |....t'..Lost Mem|
	// 00000010  6f 72 69 65 73 0c 00 31  32 2e 33 34 2e 35 36 2e  |ories..12.34.56.|
	// 00000020  37 38 39 01 00 00 00 45  01 00 00 03              |789....E....|

	if p.AccountId == 0 {
		return nil, errors.New("invalid account id")
	}
	serverList := new(bytes.Buffer)
	numServers := len(ServerMap)
	serverList.Write([]byte{0x00, byte(numServers)})
	for i := 0; i < numServers; i++ {
		server := ServerMap[i]
		serverList.Write([]byte{byte(i + 1), 0x00})
		serverList.Write(GetPortBytes(GameAuthPort))
		serverList.Write([]byte{byte(len(server.GetName())), 0x00})
		serverList.Write([]byte(server.GetName()))
		serverList.Write([]byte{byte(len(server.GetIP())), 0x00})
		serverList.Write([]byte(server.GetIP()))
		serverList.Write(
			[]byte{
				0x01, 0x00, 0x00, 0x00,
				0x00, 0x00, /* Number of people on the server, little endian, for population indicator */
				0x00, 0x00,
				byte(database.FetchUserCharacterCount(p.AccountId)), /* Number of characters the user has on the server */
			},
		)
	}

	return serverList.Bytes(), nil
}
