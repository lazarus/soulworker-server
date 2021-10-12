package packets

import (
	"bytes"
	"encoding/binary"

	. "soulworker-server/network/util"

	. "soulworker-server/global"
)

type ServerConnectRequest struct {
	SelectedServer   uint16
	SelectedServerIP string
}

func (p *ServerConnectRequest) unmarshal(buffer *bytes.Buffer) error {
	// Join server request
	// Client -> Server
	// ID=0x0205, Size=2, Total=9
	// 00000000  01 00                                             |..|

	if err := binary.Read(buffer, binary.LittleEndian, &p.SelectedServer); err != nil {
		return err
	} else {
		if int(p.SelectedServer) > len(ServerMap) {
			panic("Invalid server choice")
		}
		p.SelectedServerIP = ServerMap[p.SelectedServer-1].GetIP()
	}

	return nil
}

type ServerConnectResponse struct {
	SelectedServerIP string
}

func (p *ServerConnectResponse) marshal() ([]byte, error) {
	// Server -> Client
	// ID=0x0211, Size=18, Total=25
	// 00000000  0e 00 31 32 33 2e 34 35  36 2e 37 38 39 2e 30 32  |..123.456.789.02|
	// 00000010  76 27                                             |v'|

	serverSelect := new(bytes.Buffer)
	WriteStringUTF8NoTrailing(serverSelect, p.SelectedServerIP)
	//serverSelect.Write([]byte{byte(len(serverConnectRequest.SelectedServerIP)), 0x00})
	//serverSelect.Write([]byte(serverConnectRequest.SelectedServerIP))
	serverSelect.Write(GetPortBytes(GameAuthPort))

	return serverSelect.Bytes(), nil
}

func (p *ServerConnectResponse) id() PacketType {
	return Login_ServerConnectResponse
}
