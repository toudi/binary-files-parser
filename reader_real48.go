package binary_files_parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

type Real48 [6]byte

func (r Real48) ToFloat64() float64 {
	// https://stackoverflow.com/questions/2506942/convert-delphi-real48-to-c-sharp-double
	var result float64 = 0.0
	var exponent float64
	var mantissa float64

	if r[0] != 0 {
		exponent = float64(r[0]) - 129
		mantissa = 0.0

		for i := 1; i < 5; i += 1 {
			mantissa += float64(r[i])
			mantissa *= 0.00390625 // mantissa /= 256
		}

		mantissa += float64(r[5] & 0x7f)
		mantissa *= 0.0078125 // mantissa /= 128
		mantissa += 1.0

		if r[5]&0x80 == 0x80 {
			mantissa = -mantissa
		}

		result = mantissa * math.Pow(2, exponent)
	}

	return result
}

func decodeReal48(src io.Reader, dest reflect.Value, structField reflect.StructField, byteOrder binary.ByteOrder) (int, error) {
	var realValue Real48
	if _, err := src.Read(realValue[:]); err != nil {
		return -1, fmt.Errorf("unable to read into Real48 value: %v", err)
	}

	dest.SetFloat(realValue.ToFloat64())
	return 6, nil
}
