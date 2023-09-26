package binary_files_parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type Decoder func(src io.Reader, dest reflect.Value, structField reflect.StructField, byteOrder binary.ByteOrder) (int, error)

type BinaryDecoder struct {
	codecs    map[reflect.Kind]Decoder
	tagCodecs map[string]Decoder
	src       io.Reader
	ByteOrder binary.ByteOrder
}

// Reads as many bytes as it requires from the source reader and unpacks it to
// dst. Returns an error if it can't decode source data
func (d *BinaryDecoder) Unpack(dst interface{}) error {
	var err error

	_, err = d.unpackRecursive(d.src, dst, reflect.StructField{}, d.ByteOrder)

	if err != nil {
		return fmt.Errorf("could not read object: %v", err)
	}
	return nil
}

func (d *BinaryDecoder) RegisterCodec(dataType reflect.Kind, decoder Decoder) {
	d.codecs[dataType] = decoder
}

func (d *BinaryDecoder) RegisterTagCodec(identifier string, decoder Decoder) {
	d.tagCodecs[identifier] = decoder
}

func NewReader(src io.Reader) *BinaryDecoder {
	return &BinaryDecoder{
		src: src,
		codecs: map[reflect.Kind]Decoder{
			reflect.Uint8:  decodeUint8,
			reflect.Uint16: decodeUint16,
			reflect.Bool:   decodeBool,
		},
		tagCodecs: map[string]Decoder{
			"ts":  decodeTimestamp,
			"uts": decodeUnpackedTimestamp,
			"p":   decodeString,
			"r48": decodeReal48,
		},
		ByteOrder: binary.LittleEndian,
	}
}
