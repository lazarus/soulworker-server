package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
)

type BroochSets struct {
	NumberOfSets uint32
	Sets []BroochSet

	HashLength uint16
	Hash string
}

type BroochSet struct {
	Uint16_1	uint16
	Len1	uint16
	String1	string
	Uint16_2	uint16
	Uint16_3	uint16
	Uint16_4	uint16
	Len2	uint16
	String2	string
}

func main() {
	path := "resources/Tb_Broach_Set.res"

	data, err := ioutil.ReadFile(path)
	buffer := bytes.NewBuffer(data)

	if err != nil {
		return
	}

	// fmt.Println(hex.Dump(data))

	broochSets := &BroochSets{}

	binary.Read(buffer, binary.LittleEndian, &broochSets.NumberOfSets)

	numSets := broochSets.NumberOfSets

	broochSets.Sets = make([]BroochSet, numSets)

	for i := uint32(0); i < numSets; i++ {
		broochSet := &BroochSet{}

		binary.Read(buffer, binary.LittleEndian, &broochSet.Uint16_1)

		binary.Read(buffer, binary.LittleEndian, &broochSet.Len1)
		String1 := make([]byte, broochSet.Len1 * 2)
		binary.Read(buffer, binary.LittleEndian, &String1)
		broochSet.String1 = string(String1)

		binary.Read(buffer, binary.LittleEndian, &broochSet.Uint16_2)
		binary.Read(buffer, binary.LittleEndian, &broochSet.Uint16_3)
		binary.Read(buffer, binary.LittleEndian, &broochSet.Uint16_4)

		binary.Read(buffer, binary.LittleEndian, &broochSet.Len2)
		String2 := make([]byte, broochSet.Len2 * 2)
		binary.Read(buffer, binary.LittleEndian, &String2)
		broochSet.String1 = string(String2)

		broochSets.Sets[i] = *broochSet
	}

	binary.Read(buffer, binary.LittleEndian, &broochSets.HashLength)

	hashBuffer := make([]byte, broochSets.HashLength)
	binary.Read(buffer, binary.LittleEndian, &hashBuffer)

	broochSets.Hash = string(hashBuffer)

	fmt.Printf("%#+v\n", broochSets)
}
