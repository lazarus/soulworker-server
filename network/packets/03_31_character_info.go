package packets

import (
	"bytes"
	"encoding/binary"

	. "../structures"
	. "../util"
)

type CharacterInfoRequest struct {
	CharacterId uint32
}

func (p *CharacterInfoRequest) unmarshal(buffer *bytes.Buffer) error {
	_ = binary.Read(buffer, binary.LittleEndian, &p.CharacterId)
	_ = buffer.Next(1)

	return nil
}

type CharacterInfoResponse struct {
	MyCharacterInfo MyCharacterInfoEx
	Unknown0        byte
	Unknown1        byte
	Unknown2        byte
}

func (f *CharacterInfoResponse) marshal() ([]byte, error) {
	buf := f.MyCharacterInfo.Build()

	_ = binary.Write(buf, binary.LittleEndian, f.Unknown0)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown1)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown2)

	return buf.Bytes(), nil
}

func (f *CharacterInfoResponse) id() PacketType {
	return Game_CharacterInfoResponse
}

// Character information wrapper
type CharacterInfoEx struct {
	Character *CharacterInfo
	Position  WorldPosition
	Unknown0  float32
	Unknown1  float32
}

// Builds the struct into a byte buffer
func (f *CharacterInfoEx) Write(buf *bytes.Buffer) {
	f.Character.Write(buf)
	f.Position.Write(buf)

	_ = binary.Write(buf, binary.LittleEndian, f.Unknown0)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown1)
}

// Character information wrapper for individual characters
type MyCharacterInfoEx struct {
	CharacterInfo CharacterInfoEx
	Unknown0      uint32
	Unknown1      uint64
	Unknown2      byte
	Unknown3      byte
	Unknown4      byte
	Unknown5      byte
	Unknown6      uint32
	Unknown7      uint32
	Unknown8      string
	Unknown9      byte
	UnknownA      uint64
	UnknownB      uint32
	UnknownC      uint32
	UnknownD      uint32
}

// Builds the struct into a byte buffer
func (f *MyCharacterInfoEx) Build() *bytes.Buffer {
	buf := new(bytes.Buffer)
	f.CharacterInfo.Write(buf)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown0)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown1)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown2)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown3)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown4)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown5)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown6)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown7)

	// either 3xlong or 24 bytes
	_ = binary.Write(buf, binary.LittleEndian, uint64(0))
	_ = binary.Write(buf, binary.LittleEndian, uint64(0))
	_ = binary.Write(buf, binary.LittleEndian, uint64(0))

	WriteStringUTF8(buf, f.Unknown8)

	_ = binary.Write(buf, binary.LittleEndian, f.Unknown9)
	_ = binary.Write(buf, binary.LittleEndian, f.UnknownA)
	_ = binary.Write(buf, binary.LittleEndian, f.UnknownB)
	_ = binary.Write(buf, binary.LittleEndian, f.UnknownC)
	_ = binary.Write(buf, binary.LittleEndian, f.UnknownD)

	return buf
}
