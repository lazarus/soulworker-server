package main

// Server=localhost;Database=master;Trusted_Connection=True;
import (
	"errors"
	"time"

	"io/ioutil"

	"./global"
	"./network"

	"encoding/hex"
	"fmt"
)

func loadKeyTable(path string) error {
	if len(path) > 1 {
		buffer, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		global.KeyTable = buffer
		return nil
	}

	return errors.New("EmptyPath")
}

func quickDecode(arr []byte) {
	packetID, buffer := network.Decrypt(arr)
	fmt.Printf("ID=0x%04X, Size=%d, Total=%d\n", packetID, buffer.Len(), buffer.Len()+7)
	fmt.Println(hex.Dump(buffer.Bytes()))
}

func main() {
	loadKeyTable("resources/keyTable")

	global.ServerMap = []global.Server{
		global.Server{Name: "[NA] Austin Rocks", IP: "127.0.0.1"},
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
