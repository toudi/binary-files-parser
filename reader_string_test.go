package binary_files_parser

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

func TestPascalStringReader(t *testing.T) {
	var srcBytes = []byte{0x04, 0x41, 0x41, 0x41, 0x41}
	var reader = bytes.NewReader(srcBytes)
	var bytesRead int
	var err error

	var dstString string
	var dstValue = reflect.ValueOf(&dstString).Elem()

	bytesRead, err = decodeString(reader, dstValue, reflect.StructField{}, binary.LittleEndian)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if bytesRead != 5 {
		t.Errorf("unexpected number of bytes read: %d; want 5", bytesRead)
	}
	if dstString != "AAAA" {
		t.Errorf("unexpected end string: %s; want AAAA", dstString)
	}
}
