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

func parseValue(value []byte, t rune, order Order) interface{} {
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
	case Float16: // 2-byte float
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
	return nil
}

func buildValue(value interface{}, t rune, order rune) ([]byte, error) {
	fmt.Println("Build", t, value)
	switch t {
	case Char, SChar, UChar:
		// Type assertion to ensure value is of type int8 or uint8
		switch v := value.(type) {
		case int8:
			return []byte{byte(v)}, nil
		case uint8:
			return []byte{v}, nil
		default:
			return nil, fmt.Errorf("unsupported type for character: %T", value)
		}
	case Bool:
		switch v := value.(type) {
		case bool:
			if v {
				// true is represented as 1
				return []byte{1}, nil
			} else {
				// false is represented as 0
				return []byte{0}, nil
			}
		default:
			return nil, fmt.Errorf("unsupported type for character: %T", value)
		}

	case Short, UShort:
		// Type assertion to ensure value is of type int16 or uint16
		if v, ok := value.(int16); ok {
			bytes := make([]byte, 2)
			if order == '<' {
				binary.LittleEndian.PutUint16(bytes, uint16(v))
			} else if order == '>' {
				binary.BigEndian.PutUint16(bytes, uint16(v))
			} else {
				return nil, fmt.Errorf("unknown byte order: %v", order)
			}
			return bytes, nil
		} else if v, ok := value.(uint16); ok {
			bytes := make([]byte, 2)
			if order == '<' {
				binary.LittleEndian.PutUint16(bytes, v)
			} else if order == '>' {
				binary.BigEndian.PutUint16(bytes, v)
			} else {
				return nil, fmt.Errorf("unknown byte order: %v", order)
			}
			return bytes, nil
		}
		return nil, fmt.Errorf("unsupported type for Short/UShort: %T", value)
	// case Short:
	// 	// Type assertion to ensure value is of type int16
	// 	if v, ok := value.(int16); ok {
	// 		// Convert int16 to byte slice (little-endian)
	// 		return []byte{byte(v), byte(v >> 8)}, nil
	// 	}
	// 	return nil, fmt.Errorf("unsupported type for Short: %T", value)
	// case UShort:
	// 	// Type assertion to ensure value is of type uint16
	// 	if v, ok := value.(uint16); ok {
	// 		// Convert uint16 to byte slice (little-endian)
	// 		return []byte{byte(v), byte(v >> 8)}, nil
	// 	}
	// 	return nil, fmt.Errorf("unsupported type for UShort: %T", value)
	case Long:
		// Type assertion to ensure value is of type int32
		if v, ok := value.(int32); ok {
			// Convert int32 to byte slice (little-endian)
			bytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(bytes, uint32(v))
			return bytes, nil
		}
		return nil, fmt.Errorf("unsupported type for Long: %T", value)
	case ULong:
		// Type assertion to ensure value is of type uint32
		if v, ok := value.(uint32); ok {
			// Convert uint32 to byte slice (little-endian)
			bytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(bytes, v)
			return bytes, nil
		}
		return nil, fmt.Errorf("unsupported type for ULong: %T", value)
	case LongLong:
		// Type assertion to ensure value is of type int64
		if v, ok := value.(int64); ok {
			// Convert int64 to byte slice (little-endian)
			bytes := make([]byte, 8)
			binary.LittleEndian.PutUint64(bytes, uint64(v))
			return bytes, nil
		}
		return nil, fmt.Errorf("unsupported type for LongLong: %T", value)
	case ULongLong:
		// Type assertion to ensure value is of type uint64
		if v, ok := value.(uint64); ok {
			// Convert uint64 to byte slice (little-endian)
			bytes := make([]byte, 8)
			binary.LittleEndian.PutUint64(bytes, v)
			return bytes, nil
		}
		return nil, fmt.Errorf("unsupported type for ULongLong: %T", value)
	case SSizeT:
		// Type assertion to ensure value is of type int
		if v, ok := value.(int); ok {
			// Convert int to byte slice (little-endian)
			bytes := make([]byte, 4) // Assuming 4 bytes for SSizeT
			binary.LittleEndian.PutUint32(bytes, uint32(v))
			return bytes, nil
		}
		return nil, fmt.Errorf("unsupported type for SSizeT: %T", value)
	case SizeT:
		// Type assertion to ensure value is of type uint
		if v, ok := value.(uint); ok {
			// Convert uint to byte slice (little-endian)
			bytes := make([]byte, 4) // Assuming 4 bytes for SizeT
			binary.LittleEndian.PutUint32(bytes, uint32(v))
			return bytes, nil
		}
		return nil, fmt.Errorf("unsupported type for SizeT: %T", value)
	case Float16:
		// Type assertion to ensure value is of type float32
		if v, ok := value.(float32); ok {
			// Convert float32 to 16-bit floating-point number (half-precision float)
			intValue := math.Float32bits(v)
			halfPrecisionValue := uint16(intValue >> 16) // Use upper 16 bits
			// Convert 16-bit floating-point number to byte slice (little-endian)
			bytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(bytes, halfPrecisionValue)
			return bytes, nil
		}
		return nil, fmt.Errorf("unsupported type for Float16: %T", value)
	case Float:
		// Type assertion to ensure value is of type float32
		if v, ok := value.(float32); ok {
			// Convert float32 to byte slice (little-endian)
			bytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(bytes, math.Float32bits(v))
			return bytes, nil
		}
		return nil, fmt.Errorf("unsupported type for Float: %T", value)
	case Double:
		// Type assertion to ensure value is of type float64
		if v, ok := value.(float64); ok {
			// Convert float64 to byte slice (little-endian)
			bytes := make([]byte, 8)
			binary.LittleEndian.PutUint64(bytes, math.Float64bits(v))
			return bytes, nil
		}
		return nil, fmt.Errorf("unsupported type for Double: %T", value)
	// case CharP:
	// 	// Type assertion to ensure value is a string
	// 	if v, ok := value.(string); ok {
	// 		// Convert string to byte slice
	// 		bytes := []byte(v)
	// 		// Return the byte slice representing the string
	// 		return bytes, nil
	// 	}
	// 	return nil, fmt.Errorf("unsupported type for CharP: %T", value)
	// case VoidP:
	// 	// Type assertion to ensure value is a pointer
	// 	if v, ok := value.(uintptr); ok {
	// 		// Convert pointer to byte slice
	// 		ptrBytes := (*[unsafe.Sizeof(v)]byte)(unsafe.Pointer(v))[:]
	// 		return ptrBytes, nil
	// 	}
	// 	return nil, fmt.Errorf("unsupported type for VoidP: %T", value)
	default:
		return nil, fmt.Errorf("unknown type: %v", t)
	}
}

func Unpack(format string, data []byte) ([]interface{}, error) {

	order := LittleEndian
	var num_str string
	var parsed []interface{}
	reader := bytes.NewReader(data)

	if _, ok := OrderMap[rune(format[0])]; ok {
		order = OrderMap[rune(format[0])]
	}

	for _, t := range format {

		if unicode.IsDigit(t) {
			num_str += string(t)
			continue
		}

		if unicode.IsLetter(t) {

			if _, ok := TypesNames[rune(t)]; !ok {
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
				parsed = append(parsed, value)
				continue
			}

			for i := 0; i < num; i++ {

				if rawValue, err := readValue(reader, t); err != nil {
					return nil, fmt.Errorf("EOF: data content size less than format requires")
				} else {
					if value := parseValue(rawValue, t, order); value != nil {
						parsed = append(parsed, value)
					}
				}
			}
		}
	}
	return parsed, nil

}

// func Pack(format string, intf []interface{}) ([]byte, error) {

// 	order := '<'
// 	var num_str string
// 	var builded []byte

// 	if _, ok := OrderMap[rune(format[0])]; ok {
// 		order = rune(format[0])
// 	}

// 	i := 0

// 	for _, t := range format {

// 		if unicode.IsDigit(t) {
// 			num_str += string(t)
// 			continue
// 		}

// 		if unicode.IsLetter(t) {

// 			if _, ok := TypesNames[rune(t)]; !ok {
// 				return nil, fmt.Errorf("error: bad char ('%c') in struct format", t)
// 			}

// 			num, err := strconv.Atoi(num_str)
// 			if err != nil {
// 				num = 1
// 			}
// 			num_str = ""

// 			if t == String {
// 				if str, ok := intf[i].(string); ok {
// 					data := []byte(str)[:num]
// 					builded = append(builded, data...)
// 					i += 1
// 					continue
// 				} else {
// 					return nil, fmt.Errorf("value %v on index %d have to be a string type", intf[i], i)
// 				}
// 			}

// 			if t == PadByte {
// 				for i := 0; i < num; i++ {
// 					builded = append(builded, 0x00)
// 				}
// 				continue
// 			}

// 			for i := 0; i < num; i++ {
// 				if data, err := buildValue(intf[i], t); err != nil {
// 					return nil, fmt.Errorf("value %v on index %d have to be a %c type", intf[i], i, t)
// 				} else {
// 					builded = append(builded, data...)
// 					i += 1
// 				}
// 			}

// 			fmt.Println("Step", i, builded, t)
// 		}

// 	}

// 	fmt.Println(order, num_str, builded)

// 	return []byte{}, nil
// }
