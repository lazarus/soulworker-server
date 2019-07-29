package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
)

type BroochContainer struct {
	Count uint32
	Brooches []Brooch

	HashLength uint16
	Hash string
}

type Brooch struct {
	Uint32_1 uint32
	Uint32_2 uint32
}

func main() {
	path := "resources/TB_BroachData.res"

	data, err := ioutil.ReadFile(path)
	buffer := bytes.NewBuffer(data)

	if err != nil {
		return
	}

	broochContainer := &BroochContainer{}

	binary.Read(buffer, binary.LittleEndian, &broochContainer.Count)

	numBrooches := broochContainer.Count

	broochContainer.Brooches = make([]Brooch, numBrooches)

	for i := uint32(0); i < numBrooches; i++ {
		brooch := &Brooch{}

		binary.Read(buffer, binary.LittleEndian, &brooch.Uint32_1)
		binary.Read(buffer, binary.LittleEndian, &brooch.Uint32_2)

		broochContainer.Brooches[i] = *brooch
	}

	binary.Read(buffer, binary.LittleEndian, &broochContainer.HashLength)

	hashBuffer := make([]byte, broochContainer.HashLength)
	binary.Read(buffer, binary.LittleEndian, &hashBuffer)

	broochContainer.Hash = string(hashBuffer)

	fmt.Printf("%#+v\n", broochContainer)
}
