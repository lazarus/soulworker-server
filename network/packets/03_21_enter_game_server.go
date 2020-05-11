package packets

import (
	"bytes"
	"encoding/binary"

	. "../structures"
)

type EnterGameServerRequest struct {
	AccountId   uint32
	CharacterId uint32
	Unknown     uint64
	Unknown2    byte
	SessionKey  uint64
}

func (p *EnterGameServerRequest) unmarshal(buffer *bytes.Buffer) error {
	_ = binary.Read(buffer, binary.LittleEndian, &p.AccountId)
	_ = binary.Read(buffer, binary.LittleEndian, &p.CharacterId)
	_ = binary.Read(buffer, binary.LittleEndian, &p.Unknown)
	_ = binary.Read(buffer, binary.LittleEndian, &p.Unknown2)
	_ = binary.Read(buffer, binary.LittleEndian, &p.SessionKey)

	return nil
}

type EnterGameServerResponse struct {
	Unknown  uint32
	Result   byte // 1
	Position WorldPosition
	Unknown2 byte
	Unknown3 uint32
}

func (p *EnterGameServerResponse) marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, p.Unknown)
	_ = binary.Write(buf, binary.LittleEndian, p.Result)

	p.Position.Write(buf)

	_ = binary.Write(buf, binary.LittleEndian, p.Unknown2)
	_ = binary.Write(buf, binary.LittleEndian, p.Unknown3)

	return buf.Bytes(), nil
}

func (p *EnterGameServerResponse) id() PacketType {
	return Game_EnterGameServerResponse
}
