package pystruct

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func readValue_old(reader *bytes.Reader, t CFormatRune) ([]byte, error) {
	value := []byte{}

	for i := 0; i < FormatAlignmentMap[t]; i++ {
		b, err := reader.ReadByte()
		if err == io.EOF {
			return nil, fmt.Errorf("EOF: data content size less than format requires")
		}
		value = append(value, b)
	}
	return value, nil
}

func checkFormatAndBufSize_old(format string, expected_size int) error {
	fmt_size, err := CalcSize_old(format)
	switch {
	case err != nil:
		return err
	case fmt_size != expected_size:
		return fmt.Errorf("struct.error: unpack requires a buffer of %d bytes", fmt_size)
	default:
		return nil
	}
}

func addNum_old(num int, sRune rune) int {
	if add, err := strconv.Atoi(string(sRune)); err == nil {
		switch {
		case num == 0:
			return add
		default:
			return add * 10
		}
	}
	return 0
}

// Return the size of the struct
// (and hence of the bytes object produced by pack(format, ...))
// corresponding to the format string format
func CalcSize_old(format string) (int, error) {
	num := 0
	size := 0

	if _, ok := OrderMap[rune(format[0])]; ok {
		format = format[1:]
	}

	for _, sRune := range format {
		cFormatRune := CFormatRune(sRune)

		if add := addNum_old(num, sRune); add > 0 {
			num += add
			continue
		}

		if num == 0 {
			num = 1
		}

		if _, ok := CFormatMap[cFormatRune]; !ok {
			return -1, fmt.Errorf("struct.error: bad char ('%c') in struct format", cFormatRune)
		}

		size += num * FormatAlignmentMap[cFormatRune]
		num = 0
	}
	return size, nil
}

// Return a bytes object containing the values v1, v2, … packed according to the format string format.
// The arguments must match the values required by the format exactly.
func Pack_old(format string, intf ...interface{}) ([]byte, error) {

	if _, err := CalcSize_old(format); err != nil {
		return nil, err
	}

	num := 0
	index := 0
	var buffer []byte

	order, _ := getOrder(rune(format[0]))
	if order != nil {
		format = format[1:]
	} else {
		order = getNativeOrder()
	}

	for _, sRune := range format {
		if index+1 > len(intf) {
			return buffer, fmt.Errorf("struct.error: index error, number of interface items less than format requires")
		}

		cFormatRune := CFormatRune(sRune)

		if add := addNum_old(num, sRune); add > 0 {
			num += add
			continue
		}

		if num == 0 {
			num = 1
		}

		if _, ok := CFormatMap[cFormatRune]; !ok {
			return nil, fmt.Errorf("struct.error: bad char ('%c') in struct format", cFormatRune)
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
	if len(intf) > index {
		return buffer, fmt.Errorf("struct.error: found %d extra items that wouldn't be parsed", len(intf)-index)
	}
	return buffer, nil
}

// Pack the values v1, v2, … according to the format string format
// and write the packed bytes into the writable buffer
// starting at position offset. Note that offset is a required argument.
func PackInto_old(format string, buffer []byte, offset int, intf ...interface{}) ([]byte, error) {
	partBuf, err := Pack_old(format, intf...)
	if err != nil {
		return nil, err
	}

	if offset < 0 {
		return nil, fmt.Errorf("struct.error: offset have to be >= 0")
	}

	// Ensure buffer is large enough
	requiredLength := offset + len(partBuf)
	if requiredLength > len(buffer) {
		// Expand the buffer
		newBuffer := make([]byte, requiredLength)
		copy(newBuffer, buffer)
		buffer = newBuffer
	}

	copy(buffer[offset:], partBuf)
	return buffer, nil
}

// Unpack_old from the buffer buffer (presumably packed by Pack(format, ...))
// according to the format string format. The result is an []interface{} even if it contains exactly one item.
// The buffer’s size in bytes must match the size required by the format, as reflected by CalcSize().
func Unpack_old(format string, buffer []byte) ([]interface{}, error) {

	if err := checkFormatAndBufSize_old(format, len(buffer)); err != nil {
		return nil, err
	}

	num := 0
	var parsedValues []interface{}
	reader := bytes.NewReader(buffer)

	order, _ := getOrder(rune(format[0]))
	if order != nil {
		format = format[1:]
	} else {
		order = getNativeOrder()
	}

	for _, sRune := range format {
		cFormatRune := CFormatRune(sRune)

		if add := addNum_old(num, sRune); add > 0 {
			num += add
			continue
		}

		if num == 0 {
			num = 1
		}

		if _, ok := CFormatMap[cFormatRune]; !ok {
			return nil, fmt.Errorf("struct.error: bad char ('%c') in struct format", cFormatRune)
		}

		if cFormatRune == String {
			value := ""
			for i := 0; i < num; i++ {
				if rawValue, err := readValue_old(reader, cFormatRune); err != nil {
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

			if rawValue, err := readValue_old(reader, cFormatRune); err != nil {
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

// Iteratively unpack from the buffer buffer according to the format string format.
// This function returns an iterator which will read equally sized chunks from the buffer until all its contents have been consumed.
// The buffer’s size in bytes must be a multiple of the size required by the format, as reflected by CalcSize()
func IterUnpack_old(format string, buffer []byte) (<-chan interface{}, <-chan error) {

	parsedValues := make(chan interface{})
	errors := make(chan error)

	go func() {
		defer close(parsedValues)
		defer close(errors)

		if err := checkFormatAndBufSize_old(format, len(buffer)); err != nil {
			errors <- err
			return
		}

		num := 0
		reader := bytes.NewReader(buffer)

		order, _ := getOrder(rune(format[0]))
		if order != nil {
			format = format[1:]
		} else {
			order = getNativeOrder()
		}

		for _, sRune := range format {
			cFormatRune := CFormatRune(sRune)

			if add := addNum_old(num, sRune); add > 0 {
				num += add
				continue
			}

			if _, ok := CFormatMap[cFormatRune]; !ok {
				errors <- fmt.Errorf("struct.error: bad char ('%c') in struct format", cFormatRune)
				return
			}

			if num == 0 {
				num = 1
			}

			if cFormatRune == String {
				value := ""
				for i := 0; i < num; i++ {
					if rawValue, err := readValue_old(reader, cFormatRune); err != nil {
						errors <- err
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
				if rawValue, err := readValue_old(reader, cFormatRune); err != nil {
					errors <- err
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

// Unpack from buffer starting at position offset, according to the format string format.
// The result is an []interface{} even if it contains exactly one item.
// The buffer’s size in bytes, starting at position offset,
// must be at least the size required by the format, as reflected by CalcSize().
func UnpackFrom_old(format string, buffer []byte, offset int) ([]interface{}, error) {
	if offset >= len(buffer) {
		return nil, fmt.Errorf("offset is out of range")
	}
	return Unpack_old(format, buffer[offset:])
}

// Struct_old(fmt) --> compiled struct object
type Struct_old struct {
	format string
}

// bind method CalcSize for Struct instance
func (s *Struct_old) CalcSize() (int, error) {
	return CalcSize_old(s.format)
}

// bind method Pack for Struct instance
func (s *Struct_old) Pack(intf ...interface{}) ([]byte, error) {
	return Pack_old(s.format)
}

// bind method PackInto for Struct instance
func (s *Struct_old) PackInto(buffer []byte, offset int, intf ...interface{}) ([]byte, error) {
	return PackInto_old(s.format, buffer, offset)
}

// bind method Unpack for Struct instance
func (s *Struct_old) Unpack(buffer []byte) ([]interface{}, error) {
	return Unpack_old(s.format, buffer)
}

// bind method UnpackFrom for Struct instance
func (s *Struct_old) UnpackFrom(buffer []byte, offset int) ([]interface{}, error) {
	return UnpackFrom_old(s.format, buffer, offset)
}

// bind method IterUnpack for Struct instance
func (s *Struct_old) IterUnpack(format string, buffer []byte) (<-chan interface{}, <-chan error) {
	return IterUnpack_old(s.format, buffer)
}
