package main

import (
	"bytes"
	"time"

	"./database"
	. "./global"
	. "./network"
	"encoding/hex"
	"fmt"
)

// quickDecode quickly decodes a series of encoded bytes.
// This simulates a receive call from the server.
// It returns a byte Buffer containing the decoded byte contents.
func quickDecode(arr []byte) *bytes.Buffer {
	packetID, buffer := Decrypt(arr)
	fmt.Printf("ID=0x%04X, Size=%d, Total=%d\n", packetID, buffer.Len(), buffer.Len()+7)
	fmt.Println(hex.Dump(buffer.Bytes()))
	return buffer
}

func main() {
	// Testing information
	//buf := new(bytes.Buffer)

	//util.WriteString2(buf, "Name")
	//util.WriteString(buf, "Name")

	//fmt.Printf("%#+v\n", buf.Bytes())

	//fmt.Println(util.ReadString2(buf))
	//fmt.Println(util.ReadString(buf))

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
	fmt.Print("[Database]\tChecking connection... ")
	if err := database.CanConnect(); err != nil {
		fmt.Println("Could not connect: ", err)
		return
	} else {
		fmt.Println("Connected!")
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
