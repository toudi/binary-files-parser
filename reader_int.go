package binary_files_parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

var buffer [8]byte

func decodeBool(src io.Reader, dest reflect.Value, structField reflect.StructField, byteOrder binary.ByteOrder) (int, error) {
	bytesRead, err := src.Read(buffer[0:1])
	if err != nil {
		return -1, err
	}
	if buffer[0] > 0 {
		dest.SetBool(true)
	}
	return bytesRead, nil
}

func decodeUint8(src io.Reader, dest reflect.Value, structField reflect.StructField, byteOrder binary.ByteOrder) (int, error) {
	bytesRead, err := src.Read(buffer[0:1])
	if err != nil {
		return -1, err
	}
	dest.SetUint(uint64(buffer[0]))
	return bytesRead, nil
}

func decodeUint16(src io.Reader, dest reflect.Value, structField reflect.StructField, byteOrder binary.ByteOrder) (int, error) {
	var tmpVal uint16
	if err := binary.Read(src, byteOrder, &tmpVal); err != nil {
		return -1, fmt.Errorf("unable to decode uint16: %v", err)
	}
	dest.SetUint(uint64(tmpVal))
	return 2, nil
}
