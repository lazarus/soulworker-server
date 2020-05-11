package network

import (
	"fmt"
	"log"
	"math/rand"

	"../database"
	"../global"
	. "./packets"
)

// This file contains the contents for the Login Network

// LoginNetwork - Container for the LoginNetwork
type LoginNetwork struct {
	Network
}

// NewLoginNetwork - Creates a new LoginNetwork instance
func NewLoginNetwork() *LoginNetwork {
	network := &LoginNetwork{
		Network{
			Name: "Login Network",
			Port: global.LoginPort,
		},
	}
	network.dataHandler = network.process
	return network
}

// process - Processes data from the network from the given connection, with the given packetId and packet buffer contents
// It returns an abstract integer value
func (network *LoginNetwork) process(channel *Connection, packetID PacketType, packet interface{}) int {
	//if packetID == 0x2002 {
	//	// Basically just a hello message
	//	// Client -> Server
	//	// ID=0x2002, Size=4, Total=11
	//	// 00000000  00 00 00 00                                       |....|
	//}
	switch packetID {
	case Login_LoginAuthRequest /* KR */, Login_LoginAuthRequestGF /* GF */ :
		loginAuthRequest := packet.(*LoginAuthRequest)

		// Valid credentials => austin:coolman83
		fmt.Printf("[+] Received Login Request:\n\tUsername: %s\n\tPassword: %s\n\tMac Address: %s\n\tLeftovers: %+#v\n\n", loginAuthRequest.Username, loginAuthRequest.Password, loginAuthRequest.MAC, loginAuthRequest.Leftover)

		// Query db for username:password combo and if successful, continue
		if database.CanConnect() != nil {
			log.Println("[!] Could not connect to the database")
			break
		}

		accountId := database.VerifyLoginCredentials(loginAuthRequest.Username, loginAuthRequest.Password)

		errorCode := 0
		if accountId == 0 {
			errorCode = 1
		}

		channel.accountId = accountId

		var sessionKey uint64
		if accountId > 0 { // Session Key to keep track of a user between servers
			sessionKey = uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
			database.UpdateSessionKey(accountId, sessionKey)
			fmt.Println("Updating session key for id", accountId, ":", sessionKey)
		}

		channel.writeQueue <- &LoginAuthResponse{
			AccountId:  accountId,
			MAC:        loginAuthRequest.MAC,
			ErrorCode:  uint32(errorCode),
			SessionKey: sessionKey,
		}
		break
	case Login_ServerListRequest:
		serverListRequest := packet.(*ServerListRequest)

		channel.writeQueue <- &ServerListResponse{
			AccountId: serverListRequest.AccountId,
		}

		channel.writeQueue <- &ServerOptionsResponse{}
		break
	case Login_ServerConnectRequest:
		serverConnectRequest := packet.(*ServerConnectRequest)

		channel.writeQueue <- &ServerConnectResponse{
			SelectedServerIP: serverConnectRequest.SelectedServerIP,
		}
		break
	}

	return 0
}
