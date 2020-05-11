package packets

import "encoding/binary"

type PacketType uint16

const (
	Login_LoginAuthRequest      PacketType = 0x0201
	Login_LoginAuthRequestGF    PacketType = 0x0218
	Login_LoginAuthResponse     PacketType = 0x0202
	Login_ServerListRequest     PacketType = 0x0203
	Login_ServerListResponse    PacketType = 0x0204
	Login_ServerConnectRequest  PacketType = 0x0205
	Login_ServerConnectResponse PacketType = 0x0211
	Login_ServerOptionsResponse PacketType = 0x0231

	Character_EnterCharServerRequest  PacketType = 0x0213
	Character_EnterCharServerResponse PacketType = 0x0214
	Character_CharacterListRequest    PacketType = 0x0311
	Character_CharacterListResponse   PacketType = 0x0312
	Character_CreateCharacterRequest  PacketType = 0x0301
	Character_SelectCharacterRequest  PacketType = 0x0313
	Character_SelectCharacterResponse PacketType = 0x0314

	Game_EnterGameServerRequest  PacketType = 0x0321
	Game_EnterGameServerResponse PacketType = 0x0322
	Game_CharacterInfoRequest    PacketType = 0x0331
	Game_CharacterInfoResponse   PacketType = 0x0332

	Server_WorldCurrentDate PacketType = 0x0403
	Server_WorldVersion     PacketType = 0x0404
)

func PacketTypeFromBytes(bytes []byte) PacketType {
	return PacketType(binary.LittleEndian.Uint16(bytes))
}
