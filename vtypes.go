package cstruct

import (
	"encoding/binary"
	"math"
)

type Order rune

const ( // order and alignment
	LittleEndian Order = '<'
	BigEndian    Order = '>'
	// NativeOrderSize = '@'
	// NativeOrder     = '='
	// Network         = '!'
)

var OrderMap = map[rune]Order{
	'<': LittleEndian,
	'>': BigEndian,
	// '@': "NativeOrderSize",
	// '=': "NativeOrder",
	// '!': "Network",
}

const (
	PadByte   = 'x' // no value
	Char      = 'c' // bytes of length 1
	SChar     = 'b' // signed char -> 1 byte integer
	UChar     = 'B' // unsigned char -> 1 byte integer
	Bool      = '?' // 1 byte bool
	Short     = 'h' // short -> 2 byte integer
	UShort    = 'H' // unsigned short -> 2 byte integer
	Int       = 'i' // int -> 4 byte integer
	UInt      = 'I' // unsigned int -> 4 byte integer
	Long      = 'l' // long -> 4 byte integer
	ULong     = 'L' // unsigned long -> 4 byte integer
	LongLong  = 'q' // long long -> 8 byte integer
	ULongLong = 'Q' // unsigned long long -> 8 byte integer
	SSizeT    = 'n' // ssize_t -> integer
	SizeT     = 'N' // size_t -> integer
	Float16   = 'e' // 2 byte float
	Float     = 'f' // 4 byte float
	Double    = 'd' // 8 byte float
	String    = 's' // -> byteArray
	CharP     = 'p' // -> byteArray
	VoidP     = 'P' // -> integer
)

var TypesNames = map[rune]string{
	'x': "PadByte", // PadByte
	'c': "Char",
	'b': "SChar",
	'B': "UChar",
	'?': "_Bool",
	'h': "Short",
	'H': "UShort",
	'i': "Int",
	'I': "UInt",
	'l': "Long",
	'L': "ULong",
	'q': "LongLong",
	'Q': "ULongLong",
	'n': "SSizeT",
	'N': "SizeT",
	'e': "E",
	'f': "F",
	'd': "D",
	's': "String",
	'p': "CharP",
	'P': "VoidP",
}

var SizeMap = map[rune]int{
	'x': 1,
	'c': 1,
	'b': 1,
	'B': 1,
	'?': 1,
	'h': 2,
	'H': 2,
	'i': 4,
	'I': 4,
	'l': 4,
	'L': 4,
	'q': 8,
	'Q': 8,
	'n': 0,
	'N': 0,
	'e': 2,
	'f': 4,
	'd': 8,
	's': 1,
	'p': 1,
	'P': 0,
}

func parseChar(buffer []byte, order Order) rune {
	return rune(buffer[0])
}

func buildChar(value rune, order Order) []byte {
	return []byte{byte(value)}
}

func parseSChar(buffer []byte, order Order) int8 {
	return int8(buffer[0])
}

func buildSChar(value int8, order Order) []byte {
	return []byte{byte(value)}
}

func parseUChar(buffer []byte, order Order) uint8 {
	return uint8(buffer[0])
}

func buildUChar(value uint8, order Order) []byte {
	return []byte{value}
}

func parseBool(buffer []byte, order Order) bool {
	return buffer[0] != 0
}

func buildBool(value bool, order Order) []byte {
	if value {
		// true is represented as 1
		return []byte{1}
	} else {
		// false is represented as 0
		return []byte{0}
	}
}

func parseShort(buffer []byte, order Order) int16 {
	switch order {
	case BigEndian:
		return int16(buffer[0])<<8 | int16(buffer[1])
	case LittleEndian:
		return int16(buffer[1])<<8 | int16(buffer[0])
	default:
		return int16(buffer[1])<<8 | int16(buffer[0])
	}
}

func buildShort(value int16, order Order) []byte {
	byteValue := make([]byte, 2)
	switch order {
	case BigEndian:
		binary.BigEndian.PutUint16(byteValue, uint16(value))
	case LittleEndian:
		binary.LittleEndian.PutUint16(byteValue, uint16(value))
	default:
		binary.LittleEndian.PutUint16(byteValue, uint16(value))
	}
	return byteValue
}

func parseUShort(buffer []byte, order Order) uint16 {
	switch order {
	case BigEndian:
		return binary.BigEndian.Uint16(buffer)
	case LittleEndian:
		return binary.LittleEndian.Uint16(buffer)
	default:
		return binary.LittleEndian.Uint16(buffer)
	}
}

func buildUShort(value uint16, order Order) []byte {
	byteValue := make([]byte, 2)
	switch order {
	case BigEndian:
		binary.BigEndian.PutUint16(byteValue, value)
	case LittleEndian:
		binary.LittleEndian.PutUint16(byteValue, value)
	default:
		binary.LittleEndian.PutUint16(byteValue, value)
	}
	return byteValue
}

func parseIntLong(buffer []byte, order Order) int32 {
	switch order {
	case BigEndian:
		return int32(binary.BigEndian.Uint32(buffer))
	case LittleEndian:
		return int32(binary.LittleEndian.Uint32(buffer))
	default:
		return int32(binary.LittleEndian.Uint32(buffer))
	}
}

func buildIntLong(value int32, order Order) []byte {
	byteValue := make([]byte, 4)
	switch order {
	case BigEndian:
		binary.BigEndian.PutUint32(byteValue, uint32(value))
	case LittleEndian:
		binary.LittleEndian.PutUint32(byteValue, uint32(value))
	default:
		binary.LittleEndian.PutUint32(byteValue, uint32(value))
	}
	return byteValue
}

func parseUIntULong(buffer []byte, order Order) uint32 {
	switch order {
	case BigEndian:
		return binary.BigEndian.Uint32(buffer)
	case LittleEndian:
		return binary.LittleEndian.Uint32(buffer)
	default:
		return binary.LittleEndian.Uint32(buffer)
	}
}

func buildUIntULong(value uint32, order Order) []byte {
	byteValue := make([]byte, 4)
	switch order {
	case BigEndian:
		binary.BigEndian.PutUint32(byteValue, value)
	case LittleEndian:
		binary.LittleEndian.PutUint32(byteValue, value)
	default:
		binary.LittleEndian.PutUint32(byteValue, value)
	}
	return byteValue
}

func parseLongLong(buffer []byte, order Order) int64 {
	switch order {
	case BigEndian:
		return int64(binary.BigEndian.Uint64(buffer))
	case LittleEndian:
		return int64(binary.LittleEndian.Uint64(buffer))
	default:
		return int64(binary.LittleEndian.Uint64(buffer))
	}
}

func buildLongLong(value int64, order Order) []byte {
	byteValue := make([]byte, 8)
	switch order {
	case BigEndian:
		binary.BigEndian.PutUint64(byteValue, uint64(value))
	case LittleEndian:
		binary.LittleEndian.PutUint64(byteValue, uint64(value))
	default:
		binary.LittleEndian.PutUint64(byteValue, uint64(value))
	}
	return byteValue
}

func parseULongLong(buffer []byte, order Order) uint64 {
	switch order {
	case BigEndian:
		return binary.BigEndian.Uint64(buffer)
	case LittleEndian:
		return binary.LittleEndian.Uint64(buffer)
	default:
		return binary.LittleEndian.Uint64(buffer)
	}
}

func buildULongLong(value uint64, order Order) []byte {
	byteValue := make([]byte, 8)
	switch order {
	case BigEndian:
		binary.BigEndian.PutUint64(byteValue, value)
	case LittleEndian:
		binary.LittleEndian.PutUint64(byteValue, value)
	default:
		binary.LittleEndian.PutUint64(byteValue, value)
	}
	return byteValue
}

func parseFloat16(buffer []byte, order Order) float32 {
	switch order {
	case BigEndian:
		return math.Float32frombits(uint32(binary.BigEndian.Uint16(buffer)))
	case LittleEndian:
		return math.Float32frombits(uint32(binary.LittleEndian.Uint16(buffer)))
	default:
		return math.Float32frombits(uint32(binary.LittleEndian.Uint16(buffer)))
	}
}

func buildFloat16(value float32, order Order) []byte {
	bits := math.Float32bits(value)
	bytes := make([]byte, 4)
	switch order {
	case BigEndian:
		binary.BigEndian.PutUint32(bytes, bits)
	case LittleEndian:
		binary.LittleEndian.PutUint32(bytes, bits)
	default:
		binary.LittleEndian.PutUint32(bytes, bits)
	}
	return bytes[:2] // returning the first 2 bytes only
}

func parseFloat32(buffer []byte, order Order) float32 {
	switch order {
	case BigEndian:
		return math.Float32frombits(binary.BigEndian.Uint32(buffer))
	case LittleEndian:
		return math.Float32frombits(binary.LittleEndian.Uint32(buffer))
	default:
		return math.Float32frombits(binary.LittleEndian.Uint32(buffer))
	}
}

func buildFloat32(value float32, order Order) []byte {
	bits := math.Float32bits(value)
	bytes := make([]byte, 4)
	switch order {
	case BigEndian:
		binary.BigEndian.PutUint32(bytes, bits)
	case LittleEndian:
		binary.LittleEndian.PutUint32(bytes, bits)
	}
	return bytes
}

func parseDouble(buffer []byte, order Order) float64 {
	switch order {
	case BigEndian:
		return math.Float64frombits(binary.BigEndian.Uint64(buffer))
	case LittleEndian:
		return math.Float64frombits(binary.LittleEndian.Uint64(buffer))
	default:
		return math.Float64frombits(binary.LittleEndian.Uint64(buffer))
	}
}

func buildDouble(value float64, order Order) []byte {
	bits := math.Float64bits(value)
	bytes := make([]byte, 8)
	switch order {
	case BigEndian:
		binary.BigEndian.PutUint64(bytes, bits)
	case LittleEndian:
		binary.LittleEndian.PutUint64(bytes, bits)
	default:
		binary.LittleEndian.PutUint64(bytes, bits)
	}
	return bytes
}

func parseString(buffer []byte, order Order) string {
	return string(buffer)
}

func buildString(value string, order Order) []byte {
	return []byte(value)
}
