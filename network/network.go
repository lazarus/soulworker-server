package network

import (
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"net"

	"github.com/davecgh/go-spew/spew"

	"fmt"
	"time"

	. "soulworker-server/network/packets"
	"soulworker-server/network/structures"

	"soulworker-server/global"
)

// This class provides the base interface for all networking servers

// Network - Network struct
type Network struct {
	Name        string
	Port        uint16
	dataHandler func(*Connection, PacketType, interface{}) int
	Clients     map[*Connection]int
}

// Connection - Connection info
type Connection struct {
	conn          net.Conn
	writeQueue    chan PacketResponse
	accountId     uint32
	characterInfo *structures.CharacterInfo
}

// Start - Starts the network
func (network Network) Start() {
	sock, err := net.Listen("tcp", fmt.Sprintf(":%d", network.Port))
	if err != nil {
		panic(err)
	}
	global.Log(network.Name, fmt.Sprintf("Listening on %s.", sock.Addr().String()))

	connections := make(chan net.Conn)
	network.Clients = make(map[*Connection]int)
	clientCount := 0

	// Networking loop for accepting connections, handled in a separate thread
	go func() {
		for {
			conn, err := sock.Accept()
			if err != nil {
				fmt.Println("Could not connect", err.Error())
				continue
			}
			connections <- conn
		}
	}()

	// Main networking loop for handling connections and setting up their listen loop
	for conn := range connections {
		global.Log(network.Name, fmt.Sprintf("Client connected from %s.", conn.RemoteAddr()))
		_ = conn.SetDeadline(time.Time{})

		client := &Connection{
			conn:       conn,
			writeQueue: make(chan PacketResponse),
		}
		network.Clients[client] = clientCount
		clientCount++
		client.listen(&network)
	}
}

// Helper function to initiate the read and write loops for a given network connection
func (connection *Connection) listen(network *Network) {
	go connection.writeCycle()
	go connection.readCycle(network)
}

// Write cycle for a given network connection
func (connection *Connection) writeCycle() {
	defer func() { // Gets called on panic(), a quick and easy catch all clause to close a connection if it errors
		if r := recover(); r != nil {
			fmt.Println("Client disconnected:", r)
			_ = connection.conn.Close()
		}
	}()
	for packet := range connection.writeQueue { // Equivalent to socket select
		if raw, err := MarshalPacket(packet); err != nil {
			fmt.Println("Unable to marshal")
			panic(err)
		} else {
			if raw.Len() == 0 {
				panic("no data to send")
			}

			fmt.Println(connection.conn.LocalAddr(), "->", connection.conn.RemoteAddr())
			fmt.Printf("ID=0x%04X, Size=%d, Total=%d\n", binary.BigEndian.Uint16(raw.Bytes()[0:2]), raw.Len(), raw.Len()+5)
			fmt.Println(hex.Dump(raw.Bytes()[2:]))

			_, err = connection.conn.Write(Encrypt(raw))
			if err != nil {
				panic(err)
			}
		}
	}
}

// Read cycle for a given network connection
func (connection *Connection) readCycle(network *Network) {
	defer func() { // Gets called on panic(), a quick and easy catch all clause to close a connection if it errors
		if r := recover(); r != nil {
			fmt.Println("Client disconnected:", r)
			_ = connection.conn.Close()
		}
	}()
	for {
		// Packet header structure
		var header struct {
			Magic  uint16 // 0x02 0x00
			Size   uint16 // 0xXX 0xXX
			Sender uint8  // 0x01
		}

		res, _ := ioutil.ReadAll(connection.conn)
		spew.Dump(res)

		// Read 5 byte header into struct
		_ = binary.Read(connection.conn, binary.LittleEndian, &header)

		if header.Magic != 2 {
			panic(fmt.Sprintf("Expected a magic of 2, got %d.", header.Magic))
			//continue // Not a packet we recognize or no data to read, skip
		}

		fmt.Printf("%+v\n", header)

		if header.Size > 4095 {
			panic(fmt.Sprintf("Expected a buffer size of less than 4095, got %d.", header.Size))
		}

		packetBytes := make([]byte, header.Size-5)
		bytesRead, err := connection.conn.Read(packetBytes)
		if err != nil {
			fmt.Println("Error while reading packet data, exiting.")
			panic(err)
		}

		if bytesRead != int(header.Size-5) {
			panic(fmt.Sprintf("Expected %d bytes, got %d.", header.Size-5, bytesRead))
		}

		//packetID, buffer := Decrypt(packetBytes)
		raw := Decrypt(packetBytes[:bytesRead])

		if raw.Len() == 0 {
			continue
		}

		fmt.Println(connection.conn.RemoteAddr(), "->", connection.conn.LocalAddr())
		if packetId, packet, err := UnmarshalPacket(raw); err != nil {
			fmt.Println(err)
			fmt.Println("Unable to unmarshall")
			//panic(err)
		} else {
			// Refer to the proper networking server to handle the data
			network.dataHandler(connection, PacketType(packetId), packet)
		}
	}
} // sub_3D17D0
