package structures

import (
	"../util"
	"bytes"
	"encoding/binary"
)

// send server world date
type ServerWorldCurDate struct {
	Timestamp uint64
	Year      uint16
	Month     uint16
	Day       uint16
	Hour      uint16
	Minute    uint16
	Second    uint16
}

// Builds the struct into a byte buffer
func (f *ServerWorldCurDate) Build() *bytes.Buffer {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, f.Timestamp)
	_ = binary.Write(buf, binary.LittleEndian, f.Year)
	_ = binary.Write(buf, binary.LittleEndian, f.Month)
	_ = binary.Write(buf, binary.LittleEndian, f.Day)
	_ = binary.Write(buf, binary.LittleEndian, f.Hour)
	_ = binary.Write(buf, binary.LittleEndian, f.Minute)
	_ = binary.Write(buf, binary.LittleEndian, f.Second)
	_ = binary.Write(buf, binary.LittleEndian, uint16(0))

	return buf
}

// send server world version
type ServerWorldVersion struct {
	Unknown  uint32
	Unknown2 uint32 // 1
	Unknown3 uint32 // 0x0322
	Unknown4 uint32 // 0x3BBB
}

// Builds the struct into a byte buffer
func (f *ServerWorldVersion) Build() *bytes.Buffer {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown2)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown3)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown4)

	return buf
}

// Object world position struct
type WorldPosition struct {
	MapId       uint16
	Unknown8    uint64
	X           float32
	Y           float32
	Z           float32
	Orientation float32
}
// Builds the struct into a byte buffer
func (f *WorldPosition) Write(buf *bytes.Buffer) {
	_ = binary.Write(buf, binary.LittleEndian, f.MapId)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown8)
	_ = binary.Write(buf, binary.LittleEndian, f.X)
	_ = binary.Write(buf, binary.LittleEndian, f.Y)
	_ = binary.Write(buf, binary.LittleEndian, f.Z)
	_ = binary.Write(buf, binary.LittleEndian, f.Orientation)
}

// send enter game server res
type ServerEnterGameServerRes struct {
	Unknown uint32
	Result 	byte // 1
	Position WorldPosition
	Unknown2 byte
	Unknown3 uint32
}
// Builds the struct into a byte buffer
func (f *ServerEnterGameServerRes) Build() *bytes.Buffer {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown)
	_ = binary.Write(buf, binary.LittleEndian, f.Result)

	f.Position.Write(buf)

	_ = binary.Write(buf, binary.LittleEndian, f.Unknown2)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown3)

	return buf
}
// Character information wrapper
type CharacterInfoEx struct {
	Character CharacterInfo
	Position WorldPosition
	Unknown0 float32
	Unknown1 float32
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
	Unknown0 uint32
	Unknown1 uint64
	Unknown2 byte
	Unknown3 byte
	Unknown4 byte
	Unknown5 byte
	Unknown6 uint32
	Unknown7 uint32
	Unknown8 string
	Unknown9 byte
	UnknownA uint64
	UnknownB uint32
	UnknownC uint32
	UnknownD uint32
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

	util.WriteStringUTF8(buf, f.Unknown8)

	_ = binary.Write(buf, binary.LittleEndian, f.Unknown9)
	_ = binary.Write(buf, binary.LittleEndian, f.UnknownA)
	_ = binary.Write(buf, binary.LittleEndian, f.UnknownB)
	_ = binary.Write(buf, binary.LittleEndian, f.UnknownC)
	_ = binary.Write(buf, binary.LittleEndian, f.UnknownD)

	return buf
}
// Server response for the character list
type ServerCharacterInfoRes struct {
	MyCharacterInfo MyCharacterInfoEx
	Unknown0 byte
	Unknown1 byte
	Unknown2 byte
}
// Builds the struct into a byte buffer
func (f *ServerCharacterInfoRes) Build() *bytes.Buffer {
	buf := f.MyCharacterInfo.Build()

	_ = binary.Write(buf, binary.LittleEndian, f.Unknown0)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown1)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown2)

	return buf
}
