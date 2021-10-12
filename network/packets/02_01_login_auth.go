package packets

import (
	"bytes"
	"encoding/binary"

	"soulworker-server/network/util"
)

type LoginAuthRequest struct {
	Username string
	Password string
	MAC      string
	Leftover []byte
}

func (p *LoginAuthRequest) unmarshal(buffer *bytes.Buffer) error {
	// Login auth request (KRSW)
	// Client -> Server
	// ID=0x0201, Size=78, Total=85
	// 00000000  0e 00 58 00 58 00 58 00  58 00 58 00 58 00 58 00  |..X.X.X.X.X.X.X.|
	// 00000010  14 00 58 00 58 00 58 00  58 00 58 00 58 00 58 00  |..X.X.X.X.X.X.X.|
	// 00000020  58 00 58 00 58 00 22 00  58 00 58 00 2d 00 58 00  |X.X.X.".X.X.-.X.|
	// 00000030  58 00 2d 00 58 00 58 00  2d 00 58 00 58 00 2d 00  |X.-.X.X.-.X.X.-.|
	// 00000040  58 00 58 00 2d 00 58 00  58 00 83 96 98 00        |X.X.-.X.X.....|

	p.Username = util.ReadStringUTF16(buffer)
	//if PacketType(p.id()) == Login_LoginAuthRequest {
	//	p.Password = util.ReadStringUTF16(buffer)
	//	p.MAC = util.ReadStringUTF16(buffer)
	//} else {
	//	p.MAC = util.ReadStringUTF8(buffer)
	//}
	p.Password = util.ReadStringUTF16(buffer)
	p.MAC = util.ReadStringUTF16(buffer)

	p.Leftover = buffer.Next(4)

	return nil
}

type LoginAuthResponse struct {
	AccountId    uint32
	Unknown0     byte
	MAC          string
	ErrorMessage string
	ErrorCode    uint32
	Unknown1     byte
	Unknown2     string
	SessionKey   uint64
	Unknown3     byte
	Unknown4     uint16
	Unknown5     byte
	Unknown6     byte
	Unknown7     byte
}

func (p *LoginAuthResponse) id() PacketType {
	return Login_LoginAuthResponse
}

func (p *LoginAuthResponse) marshal() ([]byte, error) {
	// Login auth response
	// 00000000  a2 c9 00 00 01 58 58 2d  58 58 2d 58 58 2d 58 58  |.....58-58-58-58|
	// 00000010  2d 58 58 2d 58 58 00 00  00 00 00 00 00 01 0e 00  |-58-58..........|
	// 00000020  58 00 58 00 58 00 58 00  58 00 58 00 58 00 53 9f  |X.X.X.X.X.X.X.S.|
	// 00000030  2a 00 00 00 00 00 00 00  00 00 00 00              |*...........|

	loginRes := new(bytes.Buffer)

	_ = binary.Write(loginRes, binary.LittleEndian, p.AccountId) // 4 bytes
	p.Unknown0 = 1
	_ = binary.Write(loginRes, binary.LittleEndian, p.Unknown0)   // 0x01
	util.WriteStringUTF8NoLength(loginRes, p.MAC)                 // MAC address
	util.WriteStringUTF16(loginRes, p.ErrorMessage)               // errorMessage
	_ = binary.Write(loginRes, binary.LittleEndian, p.ErrorCode)  // 1 if bad login, 0 otherwise
	_ = binary.Write(loginRes, binary.LittleEndian, p.Unknown1)   // Unknown
	util.WriteStringUTF16(loginRes, p.Unknown2)                   // unknown
	_ = binary.Write(loginRes, binary.LittleEndian, p.SessionKey) // SessionKey
	_ = binary.Write(loginRes, binary.LittleEndian, p.Unknown3)   // Unknown
	_ = binary.Write(loginRes, binary.LittleEndian, p.Unknown4)   // Unknown
	_ = binary.Write(loginRes, binary.LittleEndian, p.Unknown5)   // Unknown
	_ = binary.Write(loginRes, binary.LittleEndian, p.Unknown6)   // Unknown
	_ = binary.Write(loginRes, binary.LittleEndian, p.Unknown7)   // Unknown

	return loginRes.Bytes(), nil
}
