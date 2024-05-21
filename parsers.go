package pystruct

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
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
	// PadByte   CFormatRune = 'x' // no value
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
	// SSizeT    CFormatRune = 'n' // ssize_t -> integer
	// SizeT     CFormatRune = 'N' // size_t -> integer
	// Float16   CFormatRune = 'e' // 2 byte float
	Float32 CFormatRune = 'f' // 4 byte float
	Double  CFormatRune = 'd' // 8 byte float
	String  CFormatRune = 's' // -> byteArray
	// CharP     CFormatRune = 'p' // -> byteArray
	// VoidP     CFormatRune = 'P' // -> integer
)

var CFormatMap = map[CFormatRune]CFormatRune{
	// 'x': PadByte,
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
	// 'n': SSizeT,
	// 'N': SizeT,
	// 'e': Float16,
	'f': Float32,
	'd': Double,
	's': String,
	// 'p': CharP,
	// 'P': VoidP,
}

var CFormatStringMap = map[CFormatRune]string{
	// 'x': "PadByte", // PadByte
	'c': "Char",
	'b': "SChar",
	'B': "UChar",
	'?': "Bool",
	'h': "Short",
	'H': "UShort",
	'i': "Int",
	'I': "UInt",
	'l': "Long",
	'L': "ULong",
	'q': "LongLong",
	'Q': "ULongLong",
	// 'n': "SSizeT",
	// 'N': "SizeT",
	// 'e': "Float16",
	'f': "Float32",
	'd': "Double",
	's': "String",
	// 'p': "CharP",
	// 'P': "VoidP",
}

var FormatAlignmentMap = map[CFormatRune]int{
	// 'x': 1,
	'c': 1, 'b': 1, 'B': 1, '?': 1,
	'h': 2, 'H': 2,
	'i': 4, 'I': 4, 'l': 4, 'L': 4,
	'q': 8, 'Q': 8,
	// 'n': 0,	// 'N': 0,  // 'e': 2,
	'f': 4, 'd': 8,
	's': 1,
	// 'p': 1,	// 'P': 0,
}

func getNativeOrder() binary.ByteOrder {
	var nativeEndian binary.ByteOrder
	if nativeEndian = binary.LittleEndian; nativeEndian == binary.BigEndian {
		return binary.LittleEndian
	} else {
		return binary.BigEndian
	}
}

func getOrder(order rune) (binary.ByteOrder, error) {
	if ord, ok := OrderMap[order]; ok {
		switch ord {
		case LittleEndian:
			return binary.LittleEndian, nil
		case BigEndian, NetworkOrder:
			return binary.BigEndian, nil
		default:
			return getNativeOrder(), nil
		}
	}
	return nil, fmt.Errorf("error: bad char ('%c') in struct format", order)
}

func parseString(buffer []byte) string {
	return string(buffer)
}

func buildString(value string) []byte {
	return []byte(value)
}

func parseValue(buffer []byte, cFormatRune CFormatRune, endian binary.ByteOrder) interface{} {

	// var endian binary.ByteOrder = binary.BigEndian
	// if order == LittleEndian {
	// 	endian = binary.LittleEndian
	// }

	switch cFormatRune {
	case Char:
		return rune(buffer[0])
	case SChar:
		return int8(buffer[0])
	case UChar:
		return uint8(buffer[0])
	case Bool:
		return buffer[0] != 0
	case Short:
		return int16(endian.Uint16(buffer))
	case UShort:
		return endian.Uint16(buffer)
	case Int, Long:
		return int32(endian.Uint32(buffer))
	case UInt, ULong:
		return endian.Uint32(buffer)
	case LongLong:
		return int64(endian.Uint64(buffer))
	case ULongLong:
		return endian.Uint64(buffer)
	case Float32: // 4-byte float
		return math.Float32frombits(endian.Uint32(buffer))
	case Double: // 8-byte float
		return math.Float64frombits(endian.Uint64(buffer))
	// TODO:
	// case PadByte
	// case Float16
	// case CharP:
	// case VoidP:
	default:
		return nil
	}
}

func buildValue(value interface{}, cFormatRune CFormatRune, endian binary.ByteOrder) []byte {

	buffer := new(bytes.Buffer)
	ref_val := reflect.ValueOf(value)

	switch cFormatRune {
	case Char, SChar:
		switch value.(type) {
		case rune:
			binary.Write(buffer, endian, byte(ref_val.Int()))
		}
	case UChar:
		switch value.(type) {
		case uint8:
			binary.Write(buffer, endian, byte(ref_val.Uint()))
		}
	case Bool:
		n := 0
		switch ref_val.Type().Kind() {
		case reflect.Bool:
			if ref_val.Bool() {
				n = 1
			}
			binary.Write(buffer, endian, byte(n))
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			if ref_val.Int() > 0 {
				n = 1
			}
			binary.Write(buffer, endian, byte(n))
		}
	case Short:
		switch ref_val.Type().Kind() {
		case reflect.Int8, reflect.Int16:
			binary.Write(buffer, endian, int16(ref_val.Int()))
		}
	case UShort:
		switch ref_val.Type().Kind() {
		case reflect.Uint8, reflect.Uint16:
			binary.Write(buffer, endian, uint16(ref_val.Int()))
		}
	case Int, Long:
		switch ref_val.Type().Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32:
			binary.Write(buffer, endian, int32(ref_val.Int()))
		}
	case UInt, ULong:
		switch ref_val.Type().Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32:
			binary.Write(buffer, endian, uint32(ref_val.Int()))
		}
	case LongLong:
		switch ref_val.Type().Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			binary.Write(buffer, endian, int64(ref_val.Int()))
		}
	case ULongLong:
		switch ref_val.Type().Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			binary.Write(buffer, endian, uint64(ref_val.Int()))
		}
	case Float32:
		switch ref_val.Type().Kind() {
		case reflect.Float32, reflect.Float64:
			binary.Write(buffer, endian, float32(ref_val.Float()))
		}
	case Double:
		switch ref_val.Type().Kind() {
		case reflect.Float32, reflect.Float64:
			binary.Write(buffer, endian, float64(ref_val.Float()))
		}
	// TODO:
	// case PadByte
	// case Float16:
	// case CharP:
	// case VoidP:
	default:
		return nil
	}
	return buffer.Bytes()
}
