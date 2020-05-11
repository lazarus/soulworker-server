package main

import (
	"bytes"
	"time"

	"encoding/hex"
	"fmt"

	"./database"
	. "./global"
	. "./network"
	. "./network/packets"
)

// quickDecode quickly decodes a series of encoded bytes.
// This simulates a receive call from the server.
// It returns a byte Buffer containing the decoded byte contents.
func quickDecode(arr []byte) *bytes.Buffer {
	raw := Decrypt(arr)
	if packetID, _, err := UnmarshalPacket(raw); err != nil {
		fmt.Printf("err: %s\n", err)
	} else {
		fmt.Printf("ID=0x%04X, Size=%d, Total=%d\n", packetID, raw.Len(), raw.Len()+5)
		fmt.Println(hex.Dump(raw.Bytes()))
		return raw
	}
	return nil
}

func main() {
	//packets := [][]byte{
	//	{
	//	},
	//}
	//
	//for i := 0; i < len(packets); i++ {
	//	quickDecode(packets[i][5:])
	//
	//	//charInfo := &structures.CharacterInfo{}
	//	//charInfo.Read(buffer)
	//
	//	//fmt.Printf("%#+v\n", charInfo)
	//
	//	//spew.Dump(charInfo)
	//}

	// List of game servers
	// Includes the server name and ip address
	ServerMap = []Server{
		{Name: "Yggdrasil", IP: "127.0.0.1"},
	}

	// Start the servers
	start()
}

// start is the main method for starting the server
// It starts the three component servers--login network, game network, and game world--each in a new thread
// It then waits indefinitely until it is cancelled by the console
func start() {
	Log("Starting...")

	Log("Initializing database.")
	database.Open()
	if err := database.CanConnect(); err != nil {
		panic(err)
	}

	loginNetwork := NewLoginNetwork()
	gameNetwork := NewGameNetwork()
	gameWorld := NewGameWorld()

	go loginNetwork.Start()
	go gameNetwork.Start()
	go gameWorld.Start()

	for {
		time.Sleep(10 * time.Second)
	}
}
