package pystruct

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

type cOrder rune
type cFormatRune rune

const ( // order and alignment
	tLittleEndian    cOrder = '<'
	tBigEndian       cOrder = '>'
	tNativeOrderSize cOrder = '@'
	tNativeOrder     cOrder = '='
	tNetworkOrder    cOrder = '!'
)

var cOrderMap = map[rune]cOrder{
	'<': tLittleEndian,
	'>': tBigEndian,
	'@': tNativeOrderSize,
	'=': tNativeOrder,
	'!': tNetworkOrder,
}

// var cOrderStringMap = map[cOrder]string{
// 	LittleEndian:    "LittleEndian",
// 	BigEndian:       "BigEndian",
// 	NetworkOrder:    "NetworkOrder",
// 	NativeOrder:     "NativeOrder",
// 	NativeOrderSize: "NativeOrderSize",
// 	UnknownOrder:    "UnknownOrder",
// }

const (
	// PadByte   CFormatRune = 'x' // no value
	tChar      cFormatRune = 'c' // bytes of length 1
	tSChar     cFormatRune = 'b' // signed char -> 1 byte integer
	tUChar     cFormatRune = 'B' // unsigned char -> 1 byte integer
	tBool      cFormatRune = '?' // 1 byte bool
	tShort     cFormatRune = 'h' // short -> 2 byte integer
	tUShort    cFormatRune = 'H' // unsigned short -> 2 byte integer
	tInt       cFormatRune = 'i' // int -> 4 byte integer
	tUInt      cFormatRune = 'I' // unsigned int -> 4 byte integer
	tLong      cFormatRune = 'l' // long -> 4 byte integer
	tULong     cFormatRune = 'L' // unsigned long -> 4 byte integer
	tLongLong  cFormatRune = 'q' // long long -> 8 byte integer
	tULongLong cFormatRune = 'Q' // unsigned long long -> 8 byte integer
	// SSizeT    CFormatRune = 'n' // ssize_t -> integer
	// SizeT     CFormatRune = 'N' // size_t -> integer
	// Float16   CFormatRune = 'e' // 2 byte float
	tFloat32 cFormatRune = 'f' // 4 byte float
	tDouble  cFormatRune = 'd' // 8 byte float
	tString  cFormatRune = 's' // -> byteArray
	// CharP     CFormatRune = 'p' // -> byteArray
	// VoidP     CFormatRune = 'P' // -> integer
)

var cFormatMap = map[cFormatRune]cFormatRune{
	// 'x': PadByte,
	'c': tChar,
	'b': tSChar,
	'B': tUChar,
	'?': tBool,
	'h': tShort,
	'H': tUShort,
	'i': tInt,
	'I': tUInt,
	'l': tLong,
	'L': tULong,
	'q': tLongLong,
	'Q': tULongLong,
	// 'n': SSizeT,
	// 'N': SizeT,
	// 'e': Float16,
	'f': tFloat32,
	'd': tDouble,
	's': tString,
	// 'p': CharP,
	// 'P': VoidP,
}

var cFormatStringMap = map[cFormatRune]string{
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

var formatAlignmentMap = map[cFormatRune]int{
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
	if ord, ok := cOrderMap[order]; ok {
		switch ord {
		case tLittleEndian:
			return binary.LittleEndian, nil
		case tBigEndian, tNetworkOrder:
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

func parseValue(buffer []byte, cFmtRune cFormatRune, endian binary.ByteOrder) interface{} {

	// var endian binary.ByteOrder = binary.BigEndian
	// if order == LittleEndian {
	// 	endian = binary.LittleEndian
	// }

	switch cFmtRune {
	case tChar:
		return rune(buffer[0])
	case tSChar:
		return int8(buffer[0])
	case tUChar:
		return uint8(buffer[0])
	case tBool:
		return buffer[0] != 0
	case tShort:
		return int16(endian.Uint16(buffer))
	case tUShort:
		return endian.Uint16(buffer)
	case tInt, tLong:
		return int32(endian.Uint32(buffer))
	case tUInt, tULong:
		return endian.Uint32(buffer)
	case tLongLong:
		return int64(endian.Uint64(buffer))
	case tULongLong:
		return endian.Uint64(buffer)
	case tFloat32: // 4-byte float
		return math.Float32frombits(endian.Uint32(buffer))
	case tDouble: // 8-byte float
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

func buildValue(value interface{}, cFmtRune cFormatRune, endian binary.ByteOrder) []byte {

	buffer := new(bytes.Buffer)
	ref_val := reflect.ValueOf(value)

	switch cFmtRune {
	case tChar, tSChar:
		switch value.(type) {
		case rune:
			binary.Write(buffer, endian, byte(ref_val.Int()))
		}
	case tUChar:
		switch value.(type) {
		case uint8:
			binary.Write(buffer, endian, byte(ref_val.Uint()))
		}
	case tBool:
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
	case tShort:
		switch ref_val.Type().Kind() {
		case reflect.Int8, reflect.Int16:
			binary.Write(buffer, endian, int16(ref_val.Int()))
		}
	case tUShort:
		switch ref_val.Type().Kind() {
		case reflect.Uint8, reflect.Uint16:
			binary.Write(buffer, endian, uint16(ref_val.Int()))
		}
	case tInt, tLong:
		switch ref_val.Type().Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32:
			binary.Write(buffer, endian, int32(ref_val.Int()))
		}
	case tUInt, tULong:
		switch ref_val.Type().Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32:
			binary.Write(buffer, endian, uint32(ref_val.Int()))
		}
	case tLongLong:
		switch ref_val.Type().Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			binary.Write(buffer, endian, int64(ref_val.Int()))
		}
	case tULongLong:
		switch ref_val.Type().Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			binary.Write(buffer, endian, uint64(ref_val.Int()))
		}
	case tFloat32:
		switch ref_val.Type().Kind() {
		case reflect.Float32, reflect.Float64:
			binary.Write(buffer, endian, float32(ref_val.Float()))
		}
	case tDouble:
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
