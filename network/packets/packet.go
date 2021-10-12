package packets

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	. "soulworker-server/global"
)

type PacketRequest interface {
	unmarshal(buffer *bytes.Buffer) error
}
type PacketResponse interface {
	marshal() ([]byte, error)
	id() PacketType
}

// Encrypt - Encodes a packet to bytes
func Encrypt(packet *bytes.Buffer) []byte {
	buffer := new(bytes.Buffer)

	var header struct {
		Magic  uint16 // 0x02 0x00
		Size   uint16 // 0xXX 0xXX
		Sender uint8  // 0x01
	}
	header.Magic = 0x02
	header.Size = uint16(packet.Len() + 5) // Data + 5 byte Header + 2 byte ID
	header.Sender = 0x01

	_ = binary.Write(buffer, binary.LittleEndian, header.Magic)
	_ = binary.Write(buffer, binary.LittleEndian, header.Size)
	_ = binary.Write(buffer, binary.LittleEndian, header.Sender)

	var pac = packet.Bytes()

	for i := 0; i < len(pac); i++ {
		var byte1 = pac[i]
		var byte2 = KeyTable[4*0x02-3*(i/3)+i]
		// fmt.Printf("%02X ^ %02X = %02X\n", byte1, byte2, byte1 ^ byte2)
		buffer.WriteByte(byte1 ^ byte2)
	}

	return buffer.Bytes()
}

// Decrypt - Decrypts packet data
// It returns the packetId and decrypted packet data
func Decrypt(data []byte) *bytes.Buffer {
	buffer := new(bytes.Buffer)

	for i := 0; i < len(data); i++ {
		byte1 := data[i]
		index := 4*0x02 - 3*(i/3) + i
		var byte2 = KeyTable[index]
		buffer.WriteByte(byte1 ^ byte2)
	}

	return buffer
}

func MarshalPacket(p PacketResponse) (*bytes.Buffer, error) {
	pb, err := p.marshal()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, p.id())
	buf.Write(pb)

	return buf, nil
}

func UnmarshalPacket(buffer *bytes.Buffer) (uint16, interface{}, error) {
	var packetID uint16
	err := binary.Read(buffer, binary.BigEndian, &packetID)
	if err != nil {
		fmt.Println("Failed to read packet id from buffer.")
		panic(err)
	}

	fmt.Printf("ID=0x%04X, Size=%d, Total=%d\n", packetID, buffer.Len(), buffer.Len()+5)
	fmt.Println(hex.Dump(buffer.Bytes()))

	var p PacketRequest
	switch t := PacketType(packetID); t {
	case Login_LoginAuthRequest, Login_LoginAuthRequestGF:
		p = new(LoginAuthRequest)
	case Login_ServerListRequest:
		p = new(ServerListRequest)
	case Login_ServerConnectRequest:
		p = new(ServerConnectRequest)
	case Character_EnterCharServerRequest:
		p = new(EnterCharServerRequest)
	case Character_CharacterListRequest:
		p = new(CharacterListRequest)
	case Character_CreateCharacterRequest:
		p = new(CreateCharacterRequest)
	case Character_SelectCharacterRequest:
		p = new(SelectCharacterRequest)
	case Game_EnterGameServerRequest:
		p = new(EnterGameServerRequest)
	case Game_CharacterInfoRequest:
		p = new(CharacterInfoRequest)
	default:
		//return 0, p, fmt.Errorf("packet: unrecognized Packet type: 0x%04X", t)
		break
	}

	if p != nil {
		if err := p.unmarshal(buffer); err != nil {
			return 0, p, err
		}
	}

	return packetID, p, nil
}
