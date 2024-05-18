package cstruct

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"
	"unicode"
)

const ( // order and alignment
	LittleEndian = '<'
	BigEndian    = '>'
	// NativeOrderSize = '@'
	// NativeOrder     = '='
	// Network         = '!'
)

var OrderMap = map[rune]rune{
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
	E         = 'e' // 2 byte float
	Float     = 'f' // 4 byte float
	Double    = 'd' // 8 byte float
	String    = 's' // -> byteArray
	CharP     = 'p' // -> byteArray
	VoidP     = 'P' // -> integer
)

var TypesMap = map[rune]string{
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

func readValue(reader *bytes.Reader, t rune) ([]byte, error) {
	// fmt.Printf("Parsing: %c", t)
	value := []byte{}
	for i := 0; i < SizeMap[t]; i++ {
		b, err := reader.ReadByte()
		if err == io.EOF {
			return nil, err
		}
		value = append(value, b)
	}
	return value, nil
}

func parseValue(value []byte, t rune, order rune) interface{} {
	switch t {
	case PadByte:
		return nil
	case Char:
		return value[0]
	case SChar:
		return int8(value[0])
	case UChar:
		return uint8(value[0])
	case Bool:
		return value[0] != 0
	case Short:
		if order == BigEndian {
			return int16(value[0])<<8 | int16(value[1])
		} else if order == LittleEndian {
			return int16(value[1])<<8 | int16(value[0])
		}
	case UShort:
		if order == BigEndian {
			return binary.BigEndian.Uint16(value)
		} else if order == LittleEndian {
			return binary.LittleEndian.Uint16(value)
		}
	case Int, Long:
		if order == BigEndian {
			return int32(binary.BigEndian.Uint32(value))
		} else if order == LittleEndian {
			return int32(binary.LittleEndian.Uint32(value))
		}
	case UInt, ULong:
		if order == BigEndian {
			return binary.BigEndian.Uint32(value)
		} else if order == LittleEndian {
			return binary.LittleEndian.Uint32(value)
		}
	case LongLong, SSizeT:
		if order == BigEndian {
			return int64(binary.BigEndian.Uint64(value))
		} else if order == LittleEndian {
			return int64(binary.LittleEndian.Uint64(value))
		}
	case ULongLong, SizeT:
		if order == BigEndian {
			return binary.BigEndian.Uint64(value)
		} else if order == LittleEndian {
			return binary.LittleEndian.Uint64(value)
		}
	case E: // 2-byte float
		if order == BigEndian {
			return math.Float32frombits(uint32(binary.BigEndian.Uint16(value)))
		} else if order == LittleEndian {
			return math.Float32frombits(uint32(binary.LittleEndian.Uint16(value)))
		}
	case Float: // 4-byte float
		if order == BigEndian {
			return math.Float32frombits(binary.BigEndian.Uint32(value))
		} else if order == LittleEndian {
			return math.Float32frombits(binary.LittleEndian.Uint32(value))
		}
	case Double: // 8-byte float
		if order == BigEndian {
			return math.Float64frombits(binary.BigEndian.Uint64(value))
		} else if order == LittleEndian {
			return math.Float64frombits(binary.LittleEndian.Uint64(value))
		}
	case String:
		return string(value)
	case CharP:
		n := bytes.IndexByte(value, 0)
		if n == -1 {
			n = len(value)
		}
		return string(value[:n])
	case VoidP:
		if order == BigEndian {
			return binary.BigEndian.Uint64(value)
		} else if order == LittleEndian {
			return binary.LittleEndian.Uint64(value)
		}
	default:
		return nil
	}
	return nil
}

func Unpack(format string, data []byte) ([]interface{}, error) {

	var parsedStruct []interface{}

	reader := bytes.NewReader(data)
	order := '<'

	var num_str string

	if _, ok := OrderMap[rune(format[0])]; ok {
		order = rune(format[0])
	}

	for _, t := range format {

		if unicode.IsDigit(t) {
			num_str += string(t)
			continue
		}

		if unicode.IsLetter(t) {

			if _, ok := TypesMap[rune(t)]; !ok {
				return nil, fmt.Errorf("error: bad char ('%c') in struct format", t)
			}

			num, err := strconv.Atoi(num_str)
			if err != nil {
				num = 1
			}
			num_str = ""

			if t == String {
				value := ""
				for i := 0; i < num; i++ {
					if rawValue, err := readValue(reader, t); err != nil {
						return nil, fmt.Errorf("EOF: data content size less than format requires")
					} else {
						value += string(rawValue)
					}
				}
				parsedStruct = append(parsedStruct, value)
				continue
			}

			for i := 0; i < num; i++ {

				if rawValue, err := readValue(reader, t); err != nil {
					return nil, fmt.Errorf("EOF: data content size less than format requires")
				} else {
					if value := parseValue(rawValue, t, order); value != nil {
						parsedStruct = append(parsedStruct, value)
					}
				}
			}
		}
	}
	return parsedStruct, nil

}
