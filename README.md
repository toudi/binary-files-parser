# Short description

This project came to life when I had to parse some old binary files which were produced by a program written in pascal language. Pascal uses a slightly different system for encoding the following data types:

- string (they are not null-byte terminated, but instead begin with a single byte that describes the length and then it's followed by the contents)
- Real48 (which is a floating point format that's unique for Pascal and uses 6 bytes (as opposed to float32/float64 formats))
- DateTime (which is stored as uint32, however it is not equal to Unix timestamp)

When I tried to read the structures with `binary.Read` I quickly encountered a problem: the package does not work with variadic structs. instead, it expects a fixed size. I've posted a stackoverflow question [ https://stackoverflow.com/questions/77127659/hoping-to-use-the-reader-interface-in-order-to-make-my-binary-files-import-code ] and since then I've generalized it even more in hope that the code is cleaner

# Registering codecs

There are some built-in codecs for int/uint types which use `binary.Read` under the hood. It is also possible to register your own codecs by using the following sytax:

```golang
reader := NewReader(io.Reader)
reader.RegisterCodec(reflect.Type, decoderFunction)
reader.RegisterTagCodec("identifier", decoderFunction)
```

I've added the `RegisterTagCodec` function especially because of the timestamp structure. Essentially, there is no `reflect.Type` for the `Time` structure for this exact reason - it is a struct and not a type. However, because the decoder iterates over each of the fields within a struct we can also look at the `StructField` tags - so I used it as a feature of the decoder. This means that one can very easily register the tag codec like so:

```golang
type BinaryStruct struct {
	IntField  uint16
	Timestamp time.Time `bin:"ts"`
}

reader := NewReader ( ... )
reader.RegisterTagCodec("ts", decodePascalTimestamp)
```

if one does not add the `StructTag` then the decoder will simply recurse into the struct itself and try to unpack it's fields, which is sometimes OK, but sometimes can lead
to problems.
