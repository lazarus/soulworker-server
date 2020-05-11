package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/sys/windows"
	"strings"
)

// ReadStringUTF8 reads a string from the buffer in UTF8 format.
// It returns the formatted string
func ReadStringUTF8(buffer *bytes.Buffer) string {
	// 0e 00 4c 61 7a 61 72 75 73 00 |..Lazarus.|
	stringLength := int(binary.LittleEndian.Uint16(buffer.Next(2)))
	// 0e 00 (14)

	if stringLength <= 0 {
		return ""
	}

	str := string(buffer.Next(stringLength))

	//if err != nil {
	//	fmt.Printf("[!] ReadStringUTF8 failed. %s\n", err.Error())
	//	return ""
	//}

	if len(str) != stringLength {
		fmt.Printf("[!] ReadStringUTF8 failed. Incorrect size: got %d, expected %d.\n", len(str), stringLength)
		fmt.Printf("%#+v\n", str)
	}

	return str[:stringLength]
}

// ReadStringUTF16 reads a string from the buffer in UTF16 format.
// It returns the formatted string
func ReadStringUTF16(buffer *bytes.Buffer) string {
	var str strings.Builder

	// 00000000  0e 00 4c 00 61 00 7a 00  61 00 72 00 75 00 73 00  |..L.a.z.a.r.u.s.|
	stringLength := int(binary.LittleEndian.Uint16(buffer.Next(2))) / 2
	// 0e 00 (14)

	for i := 0; i < stringLength; i++ {
		var chr uint16
		err := binary.Read(buffer, binary.LittleEndian, &chr)

		if err != nil {
			fmt.Printf("[!] ReadStringUTF16 failed. %s\n", err.Error())
			return ""
		}

		str.WriteString(windows.UTF16ToString([]uint16{chr}))
	}

	return str.String()
}

// WriteStringUTF8 writes a string to the buffer in UTF8 format
func WriteStringUTF8(buffer *bytes.Buffer, string string) {
	stringLength := uint16(len(string))

	b := make([]byte, stringLength + 2)

	binary.LittleEndian.PutUint16(b, stringLength)

	for i, v := range string {
		b[i + 2] = byte(v)
	}

	buffer.Write(b)

	if string != "" {
		buffer.Write([]byte{0x00})
	}
}

// WriteStringUTF8 writes a string to the buffer in UTF8 format with no length bytes at the beginning
// It's only used to write the mac address as of now
func WriteStringUTF8NoLength(buffer *bytes.Buffer, string string) {
	buffer.Write([]byte(string))
	buffer.Write([]byte{0x00})
}

// WriteStringUTF8 writes a string to the buffer in UTF8 format with no training null byte
func WriteStringUTF8NoTrailing(buffer *bytes.Buffer, string string) {
	_ = binary.Write(buffer, binary.LittleEndian, uint16(len(string)))
	buffer.Write([]byte(string))
}

// WriteStringUTF16 writes a string to the buffer in UTF16 format
func WriteStringUTF16(buffer *bytes.Buffer, string string) {
	stringLength := uint16(len(string) * 2)

	b := make([]byte, stringLength + 2)

	binary.LittleEndian.PutUint16(b, stringLength)

	for i, v := range string {
		b[i*2 + 2] = byte(v)
		b[i*2 + 2 + 1] = byte(0)
	}

	buffer.Write(b)
}