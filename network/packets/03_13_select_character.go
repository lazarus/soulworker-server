package packets

import (
	"bytes"
	"encoding/binary"

	"soulworker-server/network/structures"
	"soulworker-server/network/util"

	. "soulworker-server/global"
)

type SelectCharacterRequest struct {
	CharacterId uint32
	Unknown0    uint32
	Unknown1    byte
	Unknown2    uint32
	Unknown3    uint32
}

func (p *SelectCharacterRequest) unmarshal(buffer *bytes.Buffer) error {
	// Select character
	// Client -> Server
	// ID=0x0313, Size=17, Total=24
	//00000000  02 8f 31 00 00 00 00 00  01 00 00 00 00 00 00 00  |..1.............|
	//00000010  00                                                |.|

	_ = binary.Read(buffer, binary.LittleEndian, &p.CharacterId)
	_ = binary.Read(buffer, binary.LittleEndian, &p.Unknown0)
	_ = binary.Read(buffer, binary.LittleEndian, &p.Unknown1)
	_ = binary.Read(buffer, binary.LittleEndian, &p.Unknown2)
	_ = binary.Read(buffer, binary.LittleEndian, &p.Unknown3)

	return nil
}

type SelectCharacterResponse struct {
	CharacterId uint32
	AccountId   uint32
}

func (p *SelectCharacterResponse) marshal() ([]byte, error) {
	// Server -> Client
	// ID=0x0315, Size=92, Total=99
	//00000000  02 8f 31 00 91 d2 0e 00  02 02 02 00 45 37 20 00  |..1.........E7 .|
	//00000010  45 37 20 00 fd 79 00 00  79 52 02 00 00 00 00 00  |E7 ..y..yR......|
	//00000020  00 00 00 00 0e 00 32 30  36 2e 32 35 33 2e 31 37  |......206.253.17|
	//00000030  35 2e 38 32 0e 2b ff ff  00 00 00 00 00 00 00 00  |5.82.+..........|
	//00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
	//00000050  00 00 00 00 00 00 00 00  00 00 00 00              |............|

	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, &p.CharacterId)
	_ = binary.Write(buf, binary.LittleEndian, &p.AccountId)
	_ = binary.Write(buf, binary.LittleEndian, uint32(0))
	_ = binary.Write(buf, binary.LittleEndian, uint32(0))
	_ = binary.Write(buf, binary.LittleEndian, uint32(0))
	_ = binary.Write(buf, binary.LittleEndian, uint32(0))
	_ = binary.Write(buf, binary.LittleEndian, GameWorldPort)
	_ = binary.Write(buf, binary.LittleEndian, make([]byte, 10))

	util.WriteStringUTF8NoTrailing(buf, ServerMap[0].GetIP()) // Errors with a trailing null byte
	_ = binary.Write(buf, binary.LittleEndian, GameWorldPort)

	position := structures.WorldPosition{
		MapId:       10003,
		X:           10444.9951,
		Y:           10179.7461,
		Z:           100.325394,
		Orientation: 100.0,
	}
	position.Write(buf)

	_ = binary.Write(buf, binary.LittleEndian, byte(0))

	_ = binary.Write(buf, binary.LittleEndian, byte(0))
	_ = binary.Write(buf, binary.LittleEndian, byte(0))
	_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // error code?
	_ = binary.Write(buf, binary.LittleEndian, byte(0))
	_ = binary.Write(buf, binary.LittleEndian, uint32(0))

	return buf.Bytes(), nil
}

func (p *SelectCharacterResponse) id() PacketType {
	return Character_SelectCharacterResponse
}
