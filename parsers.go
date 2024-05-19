package pystruct

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Order rune
type CFormatRune rune

const ( // order and alignment
	LittleEndian    Order = '<'
	BigEndian       Order = '>'
	NativeOrderSize Order = '@'
	NativeOrder     Order = '='
	NetworkOrder    Order = '!'
	UnknownOrder    Order = 'e'
)

var OrderMap = map[rune]Order{
	'<': LittleEndian,
	'>': BigEndian,
	'@': NativeOrderSize,
	'=': NativeOrder,
	'!': NetworkOrder,
}

var OrderStringMap = map[Order]string{
	LittleEndian:    "LittleEndian",
	BigEndian:       "BigEndian",
	NetworkOrder:    "NetworkOrder",
	NativeOrder:     "NativeOrder",
	NativeOrderSize: "NativeOrderSize",
	UnknownOrder:    "UnknownOrder",
}

const (
	PadByte   CFormatRune = 'x' // no value
	Char      CFormatRune = 'c' // bytes of length 1
	SChar     CFormatRune = 'b' // signed char -> 1 byte integer
	UChar     CFormatRune = 'B' // unsigned char -> 1 byte integer
	Bool      CFormatRune = '?' // 1 byte bool
	Short     CFormatRune = 'h' // short -> 2 byte integer
	UShort    CFormatRune = 'H' // unsigned short -> 2 byte integer
	Int       CFormatRune = 'i' // int -> 4 byte integer
	UInt      CFormatRune = 'I' // unsigned int -> 4 byte integer
	Long      CFormatRune = 'l' // long -> 4 byte integer
	ULong     CFormatRune = 'L' // unsigned long -> 4 byte integer
	LongLong  CFormatRune = 'q' // long long -> 8 byte integer
	ULongLong CFormatRune = 'Q' // unsigned long long -> 8 byte integer
	SSizeT    CFormatRune = 'n' // ssize_t -> integer
	SizeT     CFormatRune = 'N' // size_t -> integer
	Float16   CFormatRune = 'e' // 2 byte float
	Float32   CFormatRune = 'f' // 4 byte float
	Double    CFormatRune = 'd' // 8 byte float
	String    CFormatRune = 's' // -> byteArray
	CharP     CFormatRune = 'p' // -> byteArray
	VoidP     CFormatRune = 'P' // -> integer
)

var CFormatMap = map[CFormatRune]CFormatRune{
	'x': PadByte,
	'c': Char,
	'b': SChar,
	'B': UChar,
	'?': Bool,
	'h': Short,
	'H': UShort,
	'i': Int,
	'I': UInt,
	'l': Long,
	'L': ULong,
	'q': LongLong,
	'Q': ULongLong,
	'n': SSizeT,
	'N': SizeT,
	'e': Float16,
	'f': Float32,
	'd': Double,
	's': String,
	'p': CharP,
	'P': VoidP,
}

var CFormatStringMap = map[CFormatRune]string{
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
	'e': "Float16",
	'f': "Float32",
	'd': "Double",
	's': "String",
	'p': "CharP",
	'P': "VoidP",
}

var SizeMap = map[CFormatRune]int{
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

func getNativeOrder() Order {
	var nativeEndian binary.ByteOrder
	if nativeEndian = binary.LittleEndian; nativeEndian == binary.BigEndian {
		return BigEndian
	} else {
		return LittleEndian
	}
}

func getOrder(order rune) (Order, error) {
	if ord, ok := OrderMap[order]; ok {
		switch ord {
		case LittleEndian:
			return LittleEndian, nil
		case BigEndian:
			return BigEndian, nil
		case NetworkOrder:
			return BigEndian, nil
		default:
			return getNativeOrder(), nil
		}
	}
	return UnknownOrder, fmt.Errorf("error: bad char ('%c') in struct format", order)
}

func parseChar(buffer []byte) rune {
	return rune(buffer[0])
}

func buildChar(value rune) []byte {
	return []byte{byte(value)}
}

func parseSChar(buffer []byte) int8 {
	return int8(buffer[0])
}

func buildSChar(value int8) []byte {
	return []byte{byte(value)}
}

func parseUChar(buffer []byte) uint8 {
	return uint8(buffer[0])
}

func buildUChar(value uint8) []byte {
	return []byte{value}
}

func parseBool(buffer []byte) bool {
	return buffer[0] != 0
}

func buildBool(value bool) []byte {
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

func parseString(buffer []byte) string {
	return string(buffer)
}

func buildString(value string) []byte {
	return []byte(value)
}

func parseValue(buffer []byte, cFormatRune CFormatRune, order Order) interface{} {
	switch cFormatRune {
	case PadByte:
		return nil
	case Char:
		return parseChar(buffer)
	case SChar:
		return parseSChar(buffer)
	case UChar:
		return parseUChar(buffer)
	case Bool:
		return parseBool(buffer)
	case Short:
		return parseShort(buffer, order)
	case UShort:
		return parseUShort(buffer, order)
	case Int, Long:
		return parseIntLong(buffer, order)
	case UInt, ULong:
		return parseUIntULong(buffer, order)
	case LongLong, SSizeT:
		return parseLongLong(buffer, order)
	case ULongLong, SizeT:
		return parseULongLong(buffer, order)
	case Float16: // 2-byte float
		return parseFloat16(buffer, order)
	case Float32: // 4-byte float
		return parseFloat32(buffer, order)
	case Double: // 8-byte float
		return parseDouble(buffer, order)
	case String:
		return parseString(buffer)
	// case CharP:
	// 	n := bytes.IndexByte(value, 0)
	// 	if n == -1 {
	// 		n = len(value)
	// 	}
	// 	return string(value[:n])
	// case VoidP:
	// 	if order == BigEndian {
	// 		return binary.BigEndian.Uint64(value)
	// 	} else if order == LittleEndian {
	// 		return binary.LittleEndian.Uint64(value)
	// 	}
	default:
		return nil
	}
}

func buildValue(value interface{}, cFormatRune CFormatRune, order Order) []byte {
	switch cFormatRune {
	case Char:
		switch v := value.(type) {
		case rune:
			return buildChar(v)
		default:
			return nil
		}
	case SChar:
		switch v := value.(type) {
		case int8:
			return buildSChar(v)
		default:
			return nil
		}
	case UChar:
		switch v := value.(type) {
		case uint8:
			return buildUChar(v)
		default:
			return nil
		}
	case Bool:
		switch v := value.(type) {
		case bool:
			return buildBool(v)
		default:
			return nil
		}
	case Short:
		switch v := value.(type) {
		case int16:
			return buildShort(v, order)
		default:
			return nil
		}
	case UShort:
		switch v := value.(type) {
		case uint16:
			return buildUShort(v, order)
		default:
			return nil
		}
	case Int, Long:
		switch v := value.(type) {
		case int32:
			return buildIntLong(v, order)
		default:
			return nil
		}
	case UInt, ULong:
		switch v := value.(type) {
		case uint32:
			return buildUIntULong(v, order)
		default:
			return nil
		}
	case LongLong:
		switch v := value.(type) {
		case int64:
			return buildLongLong(v, order)
		default:
			return nil
		}
	case ULongLong:
		switch v := value.(type) {
		case uint64:
			return buildULongLong(v, order)
		default:
			return nil
		}
	case Float16:
		switch v := value.(type) {
		case float32:
			return buildFloat16(v, order)
		default:
			return nil
		}
	case Float32:
		switch v := value.(type) {
		case float32:
			return buildFloat32(v, order)
		default:
			return nil
		}
	case Double:
		switch v := value.(type) {
		case float64:
			return buildDouble(v, order)
		default:
			return nil
		}
	default:
		return nil
	}
}
