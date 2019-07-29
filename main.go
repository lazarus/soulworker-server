package main

import (
	"time"

	"soulworker-server/global"
	"soulworker-server/network"

	"encoding/hex"
	"fmt"
)

func quickDecode(arr []byte) {
	packetID, buffer := network.Decrypt(arr)
	fmt.Printf("ID=0x%04X, Size=%d, Total=%d\n", packetID, buffer.Len(), buffer.Len()+7)
	fmt.Println(hex.Dump(buffer.Bytes()))
}

func main() {
	//buf := new(bytes.Buffer)

	//util.WriteString2(buf, "Name")
	//util.WriteString(buf, "Name")

	//fmt.Printf("%#+v\n", buf.Bytes())

	//fmt.Println(util.ReadString2(buf))
	//fmt.Println(util.ReadString(buf))

	//packets := [][]byte{
	//}
	//
	//for i := 0; i < len(packets); i++ {
	//	quickDecode(packets[i][5:])
	//}

	global.ServerMap = []global.Server{
		{Name: "Yggdrasil", IP: "127.0.0.1"},
	}

	start()
}

func start() {
	global.Log("Starting...")

	loginNetwork := network.NewLoginNetwork()
	gameNetwork := network.NewGameNetwork()
	gameWorld := network.NewGameWorld()

	go loginNetwork.Start()
	go gameNetwork.Start()
	go gameWorld.Start()

	for {
		time.Sleep(10 * time.Second)
	}
}
