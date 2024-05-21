package pystruct

import (
	"encoding/binary"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var formatPattern string = `^([@<>=!])?((\d*[cbBhHiIqQlLfds])+)$`
var groupPattern string = `(\d*)([cbBhHiIqQlLfds])`
var formatRegexp *regexp.Regexp
var groupRegexp *regexp.Regexp

func init() {
	// Compile the regular expression
	fmtRe, err := regexp.Compile(formatPattern)
	if err != nil {
		panic(fmt.Sprint("Error compiling format regex:", err))
	}
	formatRegexp = fmtRe

	// Define the regex for individual groups with separate capture for number and char
	groupRe, err := regexp.Compile(groupPattern)
	if err != nil {
		panic(fmt.Sprint("Error compiling group regex:", err))
	}
	groupRegexp = groupRe
}

type formatGroup struct {
	number int
	format CFormatRune
}

func (f *formatGroup) alignment() int {
	return alignmentMap[f.format]
}

func strip(format string) string {
	return strings.ReplaceAll(format, " ", "")
}

func parseFormat(format string) (binary.ByteOrder, []formatGroup, error) {
	var order binary.ByteOrder = getNativeOrder()
	var formatGroups []formatGroup

	format = strip(format)

	// Find the entire match with submatches
	matches := formatRegexp.FindStringSubmatch(format)
	if len(matches) == 0 {
		return nil, nil, fmt.Errorf("struct.error: Unexpected struct format %s", format)
	}

	// Extract and print the prefix if present
	if prefix := matches[1]; prefix != "" {
		order, _ = getOrder(rune(prefix[0]))
	}

	// Extract the groups (matches[2])
	// Find all individual groups
	individualMatches := groupRegexp.FindAllStringSubmatch(matches[2], -1)

	// Print each group with parsed number and char
	for _, match := range individualMatches {
		var number int

		numberStr := match[1]
		formatRune := CFormatRune(rune(match[2][0]))

		if numberStr == "" {
			number = 1
		} else {
			number, _ = strconv.Atoi(numberStr) // check on err not needed cause of match `\d` regexp
		}
		formatGroups = append(formatGroups, formatGroup{number, formatRune})
	}
	return order, formatGroups, nil
}

func parseFormatAndCalcSize(format string) (binary.ByteOrder, []formatGroup, int, error) {
	order, groups, err := parseFormat(format)
	if err != nil {
		return nil, nil, -1, err
	}
	size := 0
	for _, group := range groups {
		size += group.number * group.alignment()
	}
	return order, groups, size, nil
}

// Return the size of the struct
// (and hence of the bytes object produced by pack(format, ...))
// corresponding to the format string format
func NewCalcSize(format string) (int, error) {
	_, _, size, err := parseFormatAndCalcSize(format)
	if err != nil {
		return -1, err
	}
	return size, nil
}

// Unpack from buffer starting at position offset, according to the format string format.
// The result is an []interface{} even if it contains exactly one item.
// The buffer’s size in bytes, starting at position offset,
// must be at least the size required by the format, as reflected by CalcSize().
func NewUnpackFrom(format string, buffer []byte, offset int) ([]interface{}, error) {
	var parsedValues []interface{}

	order, groups, size, err := parseFormatAndCalcSize(format)

	if err != nil {
		return nil, err
	}

	if len(buffer)-offset != size {
		return nil, fmt.Errorf("struct.error: unpack requires a buffer of %d bytes", size)
	}

	for _, group := range groups {
		if group.format == String {
			bytesShift := group.alignment() * group.number
			value := parseString(buffer[offset : offset+bytesShift])
			offset += bytesShift
			// fmt.Printf("Fmt: %c, %v, shift->%d\n", group.format, value, bytesShift)
			parsedValues = append(parsedValues, value)
		} else {
			bytesShift := group.alignment()
			for num := 0; num < group.number; num++ {
				value := parseValue(buffer[offset:offset+bytesShift], group.format, order)
				offset += bytesShift
				// fmt.Printf("Fmt: %c, %v, shift->%d\n", group.format, value, bytesShift)
				parsedValues = append(parsedValues, value)
			}
		}
	}
	return parsedValues, nil
}

// Unpack from the buffer buffer (presumably packed by Pack(format, ...))
// according to the format string format. The result is an []interface{} even if it contains exactly one item.
// The buffer’s size in bytes must match the size required by the format, as reflected by CalcSize().
func NewUnpack(format string, buffer []byte) ([]interface{}, error) {
	return NewUnpackFrom(format, buffer, 0)
}

// // Iteratively unpack from the buffer buffer according to the format string format.
// // This function returns an iterator which will read equally sized chunks from the buffer until all its contents have been consumed.
// // The buffer’s size in bytes must be a multiple of the size required by the format, as reflected by CalcSize()
func NewIterUnpack(format string, buffer []byte) (<-chan interface{}, <-chan error) {
	parsedValues := make(chan interface{})
	errors := make(chan error)

	go func() {
		defer close(parsedValues)
		defer close(errors)

		offset := 0

		order, groups, size, err := parseFormatAndCalcSize(format)

		if err != nil {
			errors <- err
		}

		if len(buffer)-offset != size {
			errors <- fmt.Errorf("struct.error: unpack requires a buffer of %d bytes", size)
		}

		for _, group := range groups {
			if group.format == String {
				bytesShift := group.alignment() * group.number
				value := parseString(buffer[offset : offset+bytesShift])
				offset += bytesShift
				// fmt.Printf("Fmt: %c, %v, shift->%d\n", group.format, value, bytesShift)
				parsedValues <- value
			} else {
				bytesShift := group.alignment()
				for num := 0; num < group.number; num++ {
					value := parseValue(buffer[offset:offset+bytesShift], group.format, order)
					offset += bytesShift
					// fmt.Printf("Fmt: %c, %v, shift->%d\n", group.format, value, bytesShift)
					parsedValues <- value
				}
			}
		}
	}()

	return parsedValues, errors
}
