package packets

import (
	"bytes"
	"encoding/binary"
)

type EnterCharServerRequest struct {
	AccountId  uint32
	Unknown0   uint16
	SessionKey uint64
	Unknown1   byte
}

func (p *EnterCharServerRequest) unmarshal(buffer *bytes.Buffer) error {
	// Client -> Server
	// ID=0x0213, Size=15, Total=22
	// 00000000  af b7 0f 00 02 00 28 ff  bb 00 00 00 00 00 00     |......(........|

	_ = binary.Read(buffer, binary.LittleEndian, &p.AccountId)
	_ = binary.Read(buffer, binary.LittleEndian, &p.Unknown0)
	_ = binary.Read(buffer, binary.LittleEndian, &p.SessionKey)
	_ = binary.Read(buffer, binary.LittleEndian, &p.Unknown1)

	return nil
}

type EnterCharServerResponse struct {
	AccountId uint32
}

func (p *EnterCharServerResponse) marshal() ([]byte, error) {
	// Server -> Client
	// ID=0x0214, Size=5, Total=12
	// 00000000  00 af b7 0f 00                                    |.....|

	buf := new(bytes.Buffer)
	buf.Write([]byte{0x00})
	_ = binary.Write(buf, binary.LittleEndian, p.AccountId)

	return buf.Bytes(), nil
}

func (p *EnterCharServerResponse) id() PacketType {
	return Character_EnterCharServerResponse
}
