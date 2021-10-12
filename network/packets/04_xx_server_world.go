package packets

import (
	"bytes"
	"encoding/binary"
)

// ServerWorldCurDate send server world date
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
func (f *ServerWorldCurDate) marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, f.Timestamp)
	_ = binary.Write(buf, binary.LittleEndian, f.Year)
	_ = binary.Write(buf, binary.LittleEndian, f.Month)
	_ = binary.Write(buf, binary.LittleEndian, f.Day)
	_ = binary.Write(buf, binary.LittleEndian, f.Hour)
	_ = binary.Write(buf, binary.LittleEndian, f.Minute)
	_ = binary.Write(buf, binary.LittleEndian, f.Second)
	_ = binary.Write(buf, binary.LittleEndian, uint16(0))

	return buf.Bytes(), nil
}

func (f *ServerWorldCurDate) id() PacketType {
	return Server_WorldCurrentDate
}

// ServerWorldVersion send server world version
type ServerWorldVersion struct {
	Unknown  uint32
	Unknown2 uint32 // 1
	Unknown3 uint32 // 0x0322
	Unknown4 uint32 // 0x3BBB
}

// Builds the struct into a byte buffer
func (f *ServerWorldVersion) marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown2)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown3)
	_ = binary.Write(buf, binary.LittleEndian, f.Unknown4)

	return buf.Bytes(), nil
}

func (f *ServerWorldVersion) id() PacketType {
	return Server_WorldVersion
}
