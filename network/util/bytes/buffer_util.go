package bytes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"runtime"
	"strings"
	"unicode/utf16"
)

func trace() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	n := f.Name()
	return n[strings.LastIndex(n, ".")+1:]
}

type Buffer struct {
	*bytes.Buffer
	binary.ByteOrder
}

func NewBuffer(buf []byte) *Buffer {
	b := &Buffer{
		bytes.NewBuffer(buf),
		binary.LittleEndian,
	}
	return b
}

func (b *Buffer) ReadByte() byte {
	return b.ReadUint8()
}

func (b *Buffer) ReadUint8() uint8 {
	var res uint8
	if err := binary.Read(b, b.ByteOrder, &res); err != nil {
		fmt.Printf("[!] %s failed. %s\n", trace(), err.Error())
	}
	return res
}

func (b *Buffer) ReadInt16() int16 {
	var res int16
	if err := binary.Read(b, b.ByteOrder, &res); err != nil {
		fmt.Printf("[!] %s failed. %s\n", trace(), err.Error())
	}
	return res
}

func (b *Buffer) ReadUint16() uint16 {
	var res uint16
	if err := binary.Read(b, b.ByteOrder, &res); err != nil {
		fmt.Printf("[!] %s failed. %s\n", trace(), err.Error())
	}
	return res
}

func (b *Buffer) ReadInt() int {
	return int(b.ReadInt32())
}

func (b *Buffer) ReadInt32() int32 {
	var res int32
	if err := binary.Read(b, b.ByteOrder, &res); err != nil {
		fmt.Printf("[!] %s failed. %s\n", trace(), err.Error())
	}
	return res
}

func (b *Buffer) ReadUint32() uint32 {
	var res uint32
	if err := binary.Read(b, b.ByteOrder, &res); err != nil {
		fmt.Printf("[!] %s failed. %s\n", trace(), err.Error())
	}
	return res
}

func (b *Buffer) ReadFloat32() float32 {
	var res float32
	if err := binary.Read(b, b.ByteOrder, &res); err != nil {
		fmt.Printf("[!] %s failed. %s\n", trace(), err.Error())
	}
	return res
}

func (b *Buffer) ReadFloat64() float64 {
	var res float64
	if err := binary.Read(b, b.ByteOrder, &res); err != nil {
		fmt.Printf("[!] %s failed. %s\n", trace(), err.Error())
	}
	return res
}

func (b *Buffer) ReadUTF8() string {
	stringLength := int(b.ReadUint16())
	if stringLength <= 0 {
		return ""
	}
	str := b.Next(stringLength)
	if len(str) != stringLength {
		fmt.Printf("[!] ReadStringUTF8 failed. Incorrect size: got %d, expected %d.\n", len(str), stringLength)
		fmt.Printf("%#+v\n", str)
	}
	return string(str[:stringLength])
}

func (b *Buffer) ReadUTF16() string {
	length := b.ReadUint16()
	utf := make([]uint16, length)
	for i := 0; i < int(length); i++ {
		utf[i] = b.ByteOrder.Uint16(b.Next(2))
	}
	return string(utf16.Decode(utf))
}
