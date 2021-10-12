package main

import (
	"soulworker-server/database"
	. "soulworker-server/global"
	. "soulworker-server/network"
	"time"
)

func main() {
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
	//defer profile.Start(profile.MemProfile).Stop()

	Log("Starting...")

	if len(KeyTable) < 1 {
		panic("Key table has not been initialized.")
	}

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
