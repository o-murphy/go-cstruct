package cstruct

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

func readValue(reader *bytes.Reader, t CType) ([]byte, error) {
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

func parseValue(buffer []byte, t CType, order Order) interface{} {
	switch t {
	case PadByte:
		return nil
	case Char:
		return parseChar(buffer, order)
	case SChar:
		return parseSChar(buffer, order)
	case UChar:
		return parseUChar(buffer, order)
	case Bool:
		return parseBool(buffer, order)
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
		return parseString(buffer, order)
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

func buildValue(value interface{}, t CType, order Order) []byte {
	switch t {
	case Char:
		switch v := value.(type) {
		case rune:
			return buildChar(v, order)
		default:
			return nil
		}
	case SChar:
		switch v := value.(type) {
		case int8:
			return buildSChar(v, order)
		default:
			return nil
		}
	case UChar:
		switch v := value.(type) {
		case uint8:
			return buildUChar(v, order)
		default:
			return nil
		}
	case Bool:
		switch v := value.(type) {
		case bool:
			return buildBool(v, order)
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

func Unpack(format string, data []byte) ([]interface{}, error) {

	order := LittleEndian
	var num_str string
	var parsed []interface{}
	reader := bytes.NewReader(data)

	if ord, ok := OrderMap[rune(format[0])]; ok {
		order = ord
	}

	for _, cTypeRune := range format {
		cType := CType(cTypeRune)

		if unicode.IsDigit(cTypeRune) {
			num_str += string(cType)
			continue
		}

		if unicode.IsLetter(cTypeRune) {

			if _, ok := TypesNames[cType]; !ok {
				return nil, fmt.Errorf("error: bad char ('%c') in struct format", cType)
			}

			num, err := strconv.Atoi(num_str)
			if err != nil {
				num = 1
			}
			num_str = ""

			if cType == String {
				value := ""
				for i := 0; i < num; i++ {
					if rawValue, err := readValue(reader, cType); err != nil {
						return nil, fmt.Errorf("EOF: data content size less than format requires")
					} else {
						value += string(rawValue)
					}
				}
				parsed = append(parsed, value)
				continue
			}

			for i := 0; i < num; i++ {

				if rawValue, err := readValue(reader, cType); err != nil {
					return nil, fmt.Errorf("EOF: data content size less than format requires")
				} else {
					if value := parseValue(rawValue, cType, order); value != nil {
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
