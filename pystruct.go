package pystruct

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

func readValue(reader *bytes.Reader, t CFormatRune) ([]byte, error) {
	value := []byte{}

	for i := 0; i < SizeMap[t]; i++ {
		b, err := reader.ReadByte()
		if err == io.EOF {
			return nil, fmt.Errorf("EOF: data content size less than format requires")
		}
		value = append(value, b)
	}
	return value, nil
}

func CalcSize(format string) (int, error) {
	num := 0
	size := 0

	if _, ok := OrderMap[rune(format[0])]; ok {
		format = format[1:]
	}

	for _, sRune := range format {
		cFormatRune := CFormatRune(sRune)

		if add, err := strconv.Atoi(string(sRune)); err == nil {
			switch {
			case num == 0:
				num += add
			default:
				num += add * 10
			}
			continue
		}

		if num == 0 {
			num = 1
		}

		if _, ok := CFormatMap[cFormatRune]; !ok && !unicode.IsLetter(sRune) {
			return -1, fmt.Errorf("struct.error: bad char ('%c') in struct format", cFormatRune)
		}

		size += num * SizeMap[cFormatRune]
		num = 0
	}
	return size, nil
}

func checkFormatAndSize(format string, expected_size int) error {
	fmt_size, err := CalcSize(format)
	switch {
	case err != nil:
		return err
	case fmt_size < expected_size, expected_size > expected_size:
		return fmt.Errorf("struct.error: unpack requires a buffer of %d bytes", fmt_size)
	default:
		return nil
	}
}

func Pack(format string, intf []interface{}) ([]byte, error) {
	if err := checkFormatAndSize(format, len(intf)); err != nil {
		return nil, err
	}

	num := 0
	index := 0
	var buffer []byte

	order, err := getOrder(rune(format[0]))
	if err != nil {
		return nil, err
	}

	if _, ok := OrderMap[rune(format[0])]; ok {
		format = format[1:]
	}

	for _, sRune := range format {
		cFormatRune := CFormatRune(sRune)

		if add, err := strconv.Atoi(string(sRune)); err == nil {
			switch {
			case num == 0:
				num += add
			default:
				num += add * 10
			}
			continue
		}

		if _, ok := CFormatMap[cFormatRune]; !ok && !unicode.IsLetter(sRune) {
			return nil, fmt.Errorf("struct.error: bad char ('%c') in struct format", cFormatRune)
		}

		if num == 0 {
			num = 1
		}

		if cFormatRune == String {

			value := intf[index]

			switch v := value.(type) {
			case string:
				buffer = append(buffer, buildString(v)...)
			default:
				return nil, fmt.Errorf("struct.error: argument for 's' must be a bytes object")
			}
			num = 0
			index += 1
			continue
		}

		for i := 0; i < num; i++ {

			if data := buildValue(intf[index], cFormatRune, order); data != nil {
				buffer = append(buffer, data...)
				index += 1
			} else {
				return nil, fmt.Errorf("struct.error: required argument is not an %s", CFormatStringMap[cFormatRune])
			}

		}
		num = 0

	}

	return buffer, nil
}

func Unpack(format string, buffer []byte) ([]interface{}, error) {

	if err := checkFormatAndSize(format, len(buffer)); err != nil {
		return nil, err
	}

	num := 0
	var parsedValues []interface{}
	reader := bytes.NewReader(buffer)

	order, err := getOrder(rune(format[0]))
	if err != nil {
		return nil, err
	}

	if _, ok := OrderMap[rune(format[0])]; ok {
		format = format[1:]
	}

	for _, sRune := range format {
		cFormatRune := CFormatRune(sRune)

		if add, err := strconv.Atoi(string(sRune)); err == nil {
			switch {
			case num == 0:
				num += add
			default:
				num += add * 10
			}
			continue
		}

		if _, ok := CFormatMap[cFormatRune]; !ok && !unicode.IsLetter(sRune) {
			return nil, fmt.Errorf("struct.error: bad char ('%c') in struct format", cFormatRune)
		}

		if num == 0 {
			num = 1
		}

		if cFormatRune == String {
			value := ""
			for i := 0; i < num; i++ {
				if rawValue, err := readValue(reader, cFormatRune); err != nil {
					return nil, err
				} else {
					value += parseString(rawValue)
				}
			}
			parsedValues = append(parsedValues, value)
			num = 0
			continue
		}

		for i := 0; i < num; i++ {

			if rawValue, err := readValue(reader, cFormatRune); err != nil {
				return nil, err
			} else {
				if value := parseValue(rawValue, cFormatRune, order); value != nil {
					parsedValues = append(parsedValues, value)
				}
			}
		}
		num = 0

	}
	return parsedValues, nil

}

// Iterator function to unpack the data
func IterUnpack(format string, buffer []byte) (<-chan interface{}, <-chan error) {

	parsedValues := make(chan interface{})
	errors := make(chan error)

	go func() {
		defer close(parsedValues)
		defer close(errors)

		if err := checkFormatAndSize(format, len(buffer)); err != nil {
			errors <- err
			return
		}

		num := 0
		reader := bytes.NewReader(buffer)

		order, err := getOrder(rune(format[0]))
		if err != nil {
			errors <- err
			return
		}

		if _, ok := OrderMap[rune(format[0])]; ok {
			format = format[1:]
		}

		for _, sRune := range format {
			cFormatRune := CFormatRune(sRune)

			if add, err := strconv.Atoi(string(sRune)); err == nil {
				switch {
				case num == 0:
					num += add
				default:
					num += add * 10
				}
				continue
			}

			if _, ok := CFormatMap[cFormatRune]; !ok && !unicode.IsLetter(sRune) {
				errors <- fmt.Errorf("struct.error: bad char ('%c') in struct format", cFormatRune)
				return
			}

			if num == 0 {
				num = 1
			}

			if cFormatRune == String {
				value := ""
				for i := 0; i < num; i++ {
					if rawValue, err := readValue(reader, cFormatRune); err != nil {
						errors <- fmt.Errorf("EOF: data content size less than format requires")
						return
					} else {
						value += string(rawValue)
					}
				}
				parsedValues <- value
				num = 0
				continue
			}

			for i := 0; i < num; i++ {
				if rawValue, err := readValue(reader, cFormatRune); err != nil {
					errors <- fmt.Errorf("EOF: data content size less than format requires")
					return
				} else {
					if value := parseValue(rawValue, cFormatRune, order); value != nil {
						parsedValues <- value
					}
				}
			}
			num = 0
		}
	}()

	return parsedValues, errors
}

// UnpackFrom unpacks binary data from a buffer starting at a specified offset
func UnpackFrom(format string, buffer []byte, offset int) ([]interface{}, error) {
	if offset >= len(buffer) {
		return nil, fmt.Errorf("offset is out of range")
	}
	return Unpack(format, buffer[offset:])
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
