package network

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"net"

	// "io"
	// "log"
	"fmt"
	"soulworker-server/global"
	"time"
)

func check(err error, message string) {
	if err != nil {
		panic(err)
	}
	fmt.Println(message)
}

// Network - Network struct
type Network struct {
	Name        string
	Port        uint16
	dataHandler func(Connection, uint16, *bytes.Buffer) int
	Clients     map[Connection]int
}

// GetPortBytes - Gets the short bytes of a port
func GetPortBytes(port uint16) []byte {
	return []byte{byte(port), byte(port >> 8)}
}

// Connection - Connection info
type Connection struct {
	conn       net.Conn
	writeQueue chan global.Packet
}

// Start - Starts the network
func (network Network) Start() {
	sock, err := net.Listen("tcp", fmt.Sprintf(":%d", network.Port))
	if err != nil {
		panic(err)
	}
	global.Log(network.Name, fmt.Sprintf("Listening on %s.", sock.Addr().String()))

	connections := make(chan net.Conn)
	network.Clients = make(map[Connection]int)
	clientCount := 0

	go func() {
		for {
			conn, err := sock.Accept()
			if err != nil {
				panic(err)
			}
			connections <- conn
		}
	}()

	for {
		select {
		case conn := <-connections:
			global.Log(network.Name, fmt.Sprintf("Client connected from %s.", conn.RemoteAddr()))
			conn.SetDeadline(time.Time{})

			client := Connection{
				conn:       conn,
				writeQueue: make(chan global.Packet),
			}
			network.Clients[client] = clientCount
			clientCount++
			client.listen(network)
		}
	}
}

func (connection Connection) listen(network Network) {
	go connection.writeCycle(network)
	go connection.readCycle(network)
}

func (connection Connection) writeCycle(network Network) {
	for {
		select {
		case packet := <-connection.writeQueue:
			if packet.Data.Len() == 0 {
				return
			}

			fmt.Println(connection.conn.LocalAddr(), "->", connection.conn.RemoteAddr())
			fmt.Printf("ID=0x%04X, Size=%d, Total=%d\n", packet.ID, packet.Data.Len(), packet.Data.Len()+7)
			fmt.Println(hex.Dump(packet.Data.Bytes()))

			connection.conn.Write(packet.Encrypt())
		}
	}
}

func (connection Connection) readCycle(network Network) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Client disconnected:", r)
			connection.conn.Close()
		}
	}()
	if len(global.KeyTable) < 1 {
		panic("Key table has not been initialized.")
	}
	for {
		var header struct {
			Magic  uint16 // 0x02 0x00
			Size   uint16 // 0xXX 0xXX
			Sender uint8  // 0x01
		}
		// Read 5 byte header into struct
		binary.Read(connection.conn, binary.LittleEndian, &header)

		if header.Magic != 2 {
			// panic(fmt.Sprintf("Expected a magic of 2, got %d.", header.Magic))
			continue // Not a packet we recognize or no data to read, skip
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

		packetID, buffer := Decrypt(packetBytes)

		fmt.Println(connection.conn.RemoteAddr(), "->", connection.conn.LocalAddr())
		fmt.Printf("ID=0x%04X, Size=%d, Total=%d\n", packetID, buffer.Len(), buffer.Len()+7)
		fmt.Println(hex.Dump(buffer.Bytes()))

		network.dataHandler(connection, packetID, buffer)
	}
} // sub_3D17D0

// Decrypt - Decrypts packet data
func Decrypt(data []byte) (uint16, *bytes.Buffer) {
	var magic uint8 = 0x02

	buffer := new(bytes.Buffer)

	for i := 0; i < len(data); i++ {
		byte1 := data[i]
		index := 4*int(magic)-3*(i/3)+i
		var byte2 = global.KeyTable[index]
		buffer.WriteByte(byte1 ^ byte2)
	}

	var packetID uint16
	err := binary.Read(buffer, binary.BigEndian, &packetID)
	if err != nil {
		fmt.Println("Failed to read packet id from buffer.")
		panic(err)
	}

	return packetID, buffer
}
