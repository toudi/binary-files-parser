package binary_files_parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

func (d *BinaryDecoder) unpackRecursive(src io.Reader, dst interface{}, structField reflect.StructField, endianess binary.ByteOrder) (int, error) {
	var bytesRead int = 0
	var err error

	dstValue := reflect.ValueOf(dst)
	dstType := dstValue.Type()

	if dstValue.Kind() == reflect.Ptr {
		dstValue = dstValue.Elem()
		dstType = dstType.Elem()
	}

	valueType := dstValue.Kind()

	switch valueType {
	case reflect.Struct:
		// let's unpack the struct field by field.
		for i := 0; i < dstValue.NumField(); i += 1 {
			structField := dstValue.Field(i)
			// fmt.Printf("%d'th field has a name of %s and is of type %v\n", i, structField.Type().Name(), structField.Kind())

			// let's check if this is a registered tag codec
			binTag := dstType.Field(i).Tag.Get("bin")

			if binTag != "" {
				decoder, exists := d.tagCodecs[binTag]
				if exists {
					bytesRead, err = decoder(src, structField, dstType.Field(i), endianess)
					if err != nil {
						return -1, fmt.Errorf("unable to unpack field: %v", err)
					}

					continue
				}
			}

			if bytesRead, err = d.unpackRecursive(src, structField.Addr().Interface(), dstType.Field(i), endianess); err != nil {
				return -1, fmt.Errorf("unable to unpack struct: %v", err)
			}

		}
	case reflect.Array:
		declaredSize := dstType.Len()
		// fmt.Printf("begin to unpack array; declaredSize=%d\n", declaredSize)
		for i := 0; i < int(declaredSize); i += 1 {
			if bytesRead, err = d.unpackRecursive(src, dstValue.Index(i).Addr().Interface(), structField, endianess); err != nil {
				return -1, fmt.Errorf("unable to unpack array at position %d: %v", i, err)
			}
		}
	default:
		decoder, exists := d.codecs[valueType]
		if !exists {
			return -1, fmt.Errorf("unsupported value type: %v", valueType)
		}
		bytesRead, err = decoder(src, dstValue, structField, endianess)
	}

	return bytesRead, err
}
