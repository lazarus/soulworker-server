package structures

import (
	"bytes"
	"encoding/binary"
)

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
