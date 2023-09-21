package binary_files_parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"time"
)

type PascalDate uint32

func (pd PascalDate) ToTimeStruct() time.Time {
	// https://forums.codeguru.com/showthread.php?300162-Unpacking-a-Pascal-PackTime-datetime
	var sec int = int((pd>>0)&0x1f) * 2      // 0..60, only even values
	var min int = int(pd>>5) & 0x3f          // 0..59
	var hour int = int(pd>>11) & 0x1f        // 0..23
	var day int = int(pd>>16) & 0x1f         // 1..31
	var month int = int(pd>>21) & 0xf        // 1..12
	var year int = int((pd>>25)&0x7f) + 1980 // 1980..2108
	return time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)
}

func decodeTimestamp(src io.Reader, dest reflect.Value, structField reflect.StructField, byteOrder binary.ByteOrder) (int, error) {
	var tmpVal PascalDate
	if err := binary.Read(src, byteOrder, &tmpVal); err != nil {
		return -1, fmt.Errorf("unable to decode uint32: %v", err)
	}
	dest.Set(reflect.ValueOf(tmpVal.ToTimeStruct()))
	return 4, nil
}
