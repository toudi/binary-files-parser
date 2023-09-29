package binary_files_parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

func decodeString(src io.Reader, dest reflect.Value, structField reflect.StructField, byteOrder binary.ByteOrder) (int, error) {
	var buffer [256]byte
	// stringLength represents the actual string length, stored in the first byte
	var stringLength int
	// how many bytes are we supposed to read from src
	var bytesToRead int
	var err error
	// declaredLengthStruct comes from the struct tag and it represents the declared length,
	// as defined in pascal record type.
	// for instance:
	// var foo: string[200]
	// foo := "A"
	// when written to a file, would contain 2 bytes (0x01, 0x41) followed by 199 null (0x00) bytes.
	// it is important to read the null bytes as well because they belong to the record type. we can
	// safely ignore them in the later decoding part.
	var declaredLengthStruct string

	if _, err := src.Read(buffer[:1]); err != nil {
		return -1, fmt.Errorf("unable to read source bytes into buffer: %v", err)
	}
	stringLength = int(buffer[0])
	bytesToRead = stringLength
	declaredLengthStruct = structField.Tag.Get("len")
	if declaredLengthStruct != "" {
		// the structure has a fixed declared size
		bytesToRead, err = strconv.Atoi(declaredLengthStruct)
		if err != nil {
			return -1, fmt.Errorf("could not parse declared struct length: %v", err)
		}
	}
	if _, err = src.Read(buffer[1 : bytesToRead+1]); err != nil {
		return -1, fmt.Errorf("unable to read source bytes into buffer: %v", err)
	}

	dest.SetString(string(buffer[1 : stringLength+1]))

	return bytesToRead + 1, nil
}
