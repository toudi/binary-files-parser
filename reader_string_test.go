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

func TestPascalStringWithinRecordReader(t *testing.T) {
	var srcBytes = []byte{0x04, 0x41, 0x41, 0x41, 0x41, 0x01, 0x02, 0x03}
	//                    ^^^^                          ^^^^^^^^^^^^^^^^
	//                      |                                  +---> garbage that needs to be stripped
	//                      +---> actual length
	var reader = bytes.NewReader(srcBytes)
	var bytesRead int
	var err error

	var dstString string
	var dstValue = reflect.ValueOf(&dstString).Elem()

	bytesRead, err = decodeString(reader, dstValue, reflect.StructField{Tag: `len:"7"`}, binary.LittleEndian)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if bytesRead != 8 {
		t.Errorf("unexpected number of bytes read: %d; want 8", bytesRead)
	}
	if dstString != "AAAA" {
		t.Errorf("unexpected end string: %s; want AAAA", dstString)
	}
}
