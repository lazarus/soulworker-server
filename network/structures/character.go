package structures

import (
	"bytes"
	"encoding/binary"

	"soulworker-server/network/util"
)

// UnknownStructure A series of structs related to user information packet parsing
type UnknownStructure struct {
	Unknown1 int64
	Unknown2 int32
	Unknown3 uint32
	Unknown4 int64
	Unknown5 int32
	Unknown6 uint32
}

type UnknownStructure2 struct {
	Unknown1 uint32
	Unknown2 float32
	Unknown3 byte
	Unknown4 uint32
	Unknown5 byte
}

// Stats User stats
type Stats struct {
	CurrentHealth uint32
	MaximumHealth uint32
	CurrentSF     uint32
	MaximumSF     uint32
	Unknown1      uint32
	Unknown2      uint32
	Unknown3      uint32
	Stamina       uint32
	Unknown4      uint32
	Unknown5      uint32
	MoveSpeed     float32
	AttackSpeed   float32
}

// CharacterInfo Main character info struct
type CharacterInfo struct {
	Id uint32

	//UsernameLen  uint16
	//Username     []byte
	Username      string
	CharSelection byte // (1 - Haru, 2 - Erwin, 3 - Lily, 4 - Stella, 5 - Jin, 6 - Iris, 7 - Chii)
	Unknown0      byte
	UUID2         uint32
	// Maybe they actually should be uint16? :S
	/*HairStyle    uint16
	HairColor    uint16
	EyeColor     uint16
	SkinColor    uint16*/
	Appearance uint64
	Unknown1   uint64

	Level     byte
	Unknown2  byte
	AccountId uint32
	Unknown3  byte
	Unknown4  uint32
	Unknown5  byte
	Unknown6  uint32
	Unknown7  byte
	Unknown8  int32

	UnknownStructure []UnknownStructure

	Unknown9  uint32
	Unknown10 uint32
	Unknown11 uint32
	Unknown12 uint32
	//Unknown13Len uint16
	//Unknown13	 []byte
	Unknown13 string
	Unknown14 uint32

	Stats

	Unknown15 byte
	//Unknown16Len uint16
	//Unknown16	 []byte
	Unknown16   string
	Energy      uint16
	BonusEnergy uint16
	Unknown17   uint16
	Unknown18   byte
	Unknown19   uint32
	Unknown20   byte
	Unknown21   uint32

	UnknownCount      byte
	UnknownStructure2 []UnknownStructure2

	Index     byte
	Unknown22 uint32

	Outfit uint32
}

// Populates the fields of a CharacterInfo object from the contents of the byte buffer
func (charInfo *CharacterInfo) Read(buffer *bytes.Buffer) {

	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Id)

	// Base
	//_ = binary.Read(buffer, binary.LittleEndian, &charInfo.UsernameLen)
	//charInfo.Username = make([]byte, charInfo.UsernameLen)
	//_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Username)
	charInfo.Username = util.ReadStringUTF16(buffer)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.CharSelection)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown0)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.UUID2)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Appearance)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown1)

	// Other
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Level)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown2)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.AccountId)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown3)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown4)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown5)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown6)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown7)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown8)

	// Unknown Struct (Items?)
	charInfo.UnknownStructure = make([]UnknownStructure, 13)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.UnknownStructure)

	// Unknown
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown9)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown10)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown11)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown12)
	//_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown13Len)
	//charInfo.Unknown13 = make([]byte, charInfo.Unknown13Len)
	//_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown13)
	charInfo.Unknown13 = util.ReadStringUTF8(buffer)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown14)

	// Stats
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Stats)

	// Unknown
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown15)
	//_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown16Len)
	//charInfo.Unknown16 = make([]byte, charInfo.Unknown16Len)
	//_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown16)
	charInfo.Unknown16 = util.ReadStringUTF8(buffer)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Energy)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.BonusEnergy)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown17)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown18)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown19)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown20)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown21)

	// Unknown Chunk
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.UnknownCount)
	charInfo.UnknownStructure2 = make([]UnknownStructure2, charInfo.UnknownCount)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.UnknownStructure2)

	// Footer
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Index)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Unknown22)
	_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Outfit)
}

// Populates a byte buffer with the contents of a CharacterInfo object
func (charInfo *CharacterInfo) Write(buffer *bytes.Buffer) {

	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Id)

	// Base
	//_ = binary.Read(buffer, binary.LittleEndian, &charInfo.UsernameLen)
	//charInfo.Username = make([]byte, charInfo.UsernameLen)
	//_ = binary.Read(buffer, binary.LittleEndian, &charInfo.Username)
	util.WriteStringUTF16(buffer, charInfo.Username)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.CharSelection)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown0)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.UUID2)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Appearance)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown1)

	// Other
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Level)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown2)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.AccountId)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown3)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown4)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown5)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown6) // 110011301 + (1000000 * Class)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown7)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown8)

	// Unknown Struct (Items?)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.UnknownStructure)

	// Unknown
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown9)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown10)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown11)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown12)
	util.WriteStringUTF8(buffer, charInfo.Unknown13)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown14)

	// Stats
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Stats)

	// Unknown
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown15)
	util.WriteStringUTF8(buffer, charInfo.Unknown16)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Energy)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.BonusEnergy)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown17)

	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown18)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown19)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown20)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown21)

	// Unknown Chunk
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.UnknownCount)
	for _, s := range charInfo.UnknownStructure2 {
		_ = binary.Write(buffer, binary.LittleEndian, s)
	}

	// Footer
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Index)
	_ = binary.Write(buffer, binary.LittleEndian, charInfo.Unknown22)
}

// CharacterModel wrapper for easier read/write
type CharacterModel struct {
	Id         uint32
	AccountId  uint32
	Index      byte
	Name       string
	Class      byte
	Level      byte
	Appearance uint64

	MapId uint16
	X     float32
	Y     float32
	Z     float32
	O     float32
}
