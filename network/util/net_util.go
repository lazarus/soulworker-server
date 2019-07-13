package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/sys/windows"
	"strings"
)

func ReadString(buffer *bytes.Buffer) string {
	// 0e 00 4c 61 7a 61 72 75 73 00 |..Lazarus.|
	stringLength := int(binary.LittleEndian.Uint16(buffer.Next(2)))
	// 0e 00 (14)
	str, err := buffer.ReadString(0x00)

	if err != nil {
		fmt.Printf("[!] ReadString failed. %s\n", err.Error())
		return ""
	}

	if len(str) != stringLength + 1 {
		fmt.Printf("[!] ReadString failed. Incorrect size: got %d, expected %d.\n", len(str), stringLength)
		fmt.Printf("%#+v\n", str)
	}

	return str[:stringLength]
}

func ReadString2(buffer *bytes.Buffer) string {
	var str strings.Builder

	// 00000000  0e 00 4c 00 61 00 7a 00  61 00 72 00 75 00 73 00  |..L.a.z.a.r.u.s.|
	stringLength := int(binary.LittleEndian.Uint16(buffer.Next(2))) / 2
	// 0e 00 (14)

	for i := 0; i < stringLength; i++ {
		var chr uint16
		err := binary.Read(buffer, binary.LittleEndian, &chr)

		if err != nil {
			fmt.Printf("[!] ReadString2 failed. %s\n", err.Error())
			return ""
		}

		str.WriteString(windows.UTF16ToString([]uint16{chr}))
	}

	return str.String()
}

func WriteString(buffer *bytes.Buffer, string string) {
	stringLength := uint16(len(string))

	b := make([]byte, stringLength + 2)

	binary.LittleEndian.PutUint16(b, stringLength)

	for i, v := range string {
		b[i + 2] = byte(v)
	}

	buffer.Write(b)
	buffer.Write([]byte{0x00})
}

func WriteString2 (buffer *bytes.Buffer, string string) {
	stringLength := uint16(len(string) * 2)

	b := make([]byte, stringLength + 2)

	binary.LittleEndian.PutUint16(b, stringLength)

	for i, v := range string {
		b[i*2 + 2] = byte(v)
		b[i*2 + 2 + 1] = byte(0)
	}

	buffer.Write(b)
}