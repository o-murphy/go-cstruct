package pystruct

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

var (
	// DefaultTagName is the default tag name for struct fields which provides
	// a more granular to tweak certain structs. Lookup the necessary functions
	// for more info.
	DefaultTagName = "pystruct" // struct's field default tag name
)

func readValue(reader *bytes.Reader, t CFormatRune) ([]byte, error) {
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

func CalcSize(format string) (int, error) {
	numStr := ""
	size := 0

	order, err := getOrder(rune(format[0]))
	if err != nil {
		return -1, err
	}
	fmt.Printf("Order: %s\n", OrderStringMap[order])

	for _, sRune := range format {
		cFormatRune := CFormatRune(sRune)

		if unicode.IsDigit(sRune) {
			numStr += string(cFormatRune)
			continue
		}

		num, err := strconv.Atoi(numStr)
		if err != nil {
			num = 1
		}
		numStr = ""

		if unicode.IsLetter(sRune) {

			if _, ok := CFormatMap[cFormatRune]; !ok {
				return -1, fmt.Errorf("error: bad char ('%c') in struct format", cFormatRune)
			}

		}

		size += num * SizeMap[cFormatRune]
	}
	return size, nil
}

func Unpack(format string, buffer []byte) ([]interface{}, error) {

	var numStr string
	var parsedValues []interface{}
	reader := bytes.NewReader(buffer)

	order, err := getOrder(rune(format[0]))
	if err != nil {
		return nil, err
	}

	for _, sRune := range format {
		cFormatRune := CFormatRune(sRune)

		if unicode.IsDigit(sRune) {
			numStr += string(cFormatRune)
			continue
		}

		if unicode.IsLetter(sRune) {

			if _, ok := CFormatMap[cFormatRune]; !ok {
				return nil, fmt.Errorf("error: bad char ('%c') in struct format", cFormatRune)
			}

			num, err := strconv.Atoi(numStr)
			if err != nil {
				num = 1
			}
			numStr = ""

			if cFormatRune == String {
				value := ""
				for i := 0; i < num; i++ {
					if rawValue, err := readValue(reader, cFormatRune); err != nil {
						return nil, fmt.Errorf("EOF: data content size less than format requires")
					} else {
						value += string(rawValue)
					}
				}
				parsedValues = append(parsedValues, value)
				continue
			}

			for i := 0; i < num; i++ {

				if rawValue, err := readValue(reader, cFormatRune); err != nil {
					return nil, fmt.Errorf("EOF: data content size less than format requires")
				} else {
					if value := parseValue(rawValue, cFormatRune, order); value != nil {
						parsedValues = append(parsedValues, value)
					}
				}
			}
		}
	}
	return parsedValues, nil

}

// Iterator function to unpack the data
func IterUnpack(format string, buffer *[]byte) (<-chan interface{}, <-chan error) {
	parsedValues := make(chan interface{})
	errors := make(chan error)

	go func() {
		defer close(parsedValues)
		defer close(errors)

		var numStr string
		reader := bytes.NewReader(*buffer)

		order, err := getOrder(rune(format[0]))
		if err != nil {
			errors <- err
			return
		}

		for _, sRune := range format {
			cFormatRune := CFormatRune(sRune)

			if unicode.IsDigit(sRune) {
				numStr += string(cFormatRune)
				continue
			}

			if unicode.IsLetter(sRune) {

				if _, ok := CFormatMap[cFormatRune]; !ok {
					errors <- fmt.Errorf("error: bad char ('%c') in struct format", cFormatRune)
					return
				}

				num, err := strconv.Atoi(numStr)
				if err != nil {
					num = 1
				}
				numStr = ""

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
			}
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
