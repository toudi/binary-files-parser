package main

import (
	"bytes"
	"fmt"
	"time"

	bfp "github.com/toudi/binary-files-parser"
)

type NestedStruct struct {
	PasString string `bin:"p"`
}

type NestedStruct2 struct {
	SomeField  uint16
	SomeString string `bin:"p"`
}

type BinaryStruct struct {
	IntField  uint16
	Nested    NestedStruct
	Array     [10]byte
	Timestamp time.Time `bin:"ts"`
	Real48    float64   `bin:"r48"`
	Structs   [2]NestedStruct2
}

func main() {
	var srcBytes = []byte{
		0x00, 0x01, // uint16
		0x04, 0x41, 0x41, 0x41, 0x41, // string
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, // array of 10 bytes
		0x00, 0x00, 0x00, 0x00, // timestamp as uint32
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // real48
		// NestedStruct2, element 0
		0x00, 0x01, // uint16
		0x01, 0x41, // string
		// NestedStruct2, element 1
		0x00, 0x02, // uint16
		0x01, 0x42, // string
	}
	var reader = bytes.NewReader(srcBytes)

	var bin BinaryStruct

	var err error

	decoder := bfp.NewReader(reader)

	err = decoder.Unpack(&bin)

	if err != nil {
		fmt.Printf("Error unpacking structure: %v\n", err)
	}

	fmt.Printf("end structure: %+v\n", bin)
}
