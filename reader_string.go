package binary_files_parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

func decodeString(src io.Reader, dest reflect.Value, structField reflect.StructField, byteOrder binary.ByteOrder) (int, error) {
	var buffer [256]byte
	if _, err := src.Read(buffer[:1]); err != nil {
		return -1, fmt.Errorf("unable to read source bytes into buffer: %v", err)
	}
	declaredLength := int(buffer[0])
	if _, err := src.Read(buffer[1 : declaredLength+1]); err != nil {
		return -1, fmt.Errorf("unable to read source bytes into buffer: %v", err)
	}
	dest.SetString(string(buffer[1 : declaredLength+1]))

	return declaredLength + 1, nil
}
