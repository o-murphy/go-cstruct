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
	Number int
	Format CFormatRune
}

func (f *formatGroup) Alignment() int {
	return alignmentMap[f.Format]
}

func ParseFormat(format string) (binary.ByteOrder, []formatGroup, error) {
	var order binary.ByteOrder = getNativeOrder()
	var formatGroups []formatGroup

	format = strings.ReplaceAll(format, " ", "")

	// Find the entire match with submatches
	matches := formatRegexp.FindStringSubmatch(format)
	if len(matches) == 0 {
		return nil, nil, fmt.Errorf("struct.error: wrong struct format %s", format)
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
		// fmt.Printf("Group: %s, Number: %d, Char: %s\n", match[0], number, char)
		formatGroups = append(formatGroups, formatGroup{number, formatRune})
	}
	return order, formatGroups, nil
}

func CalcFormatSize(format string) (int, error) {
	_, groups, err := ParseFormat(format)
	if err != nil {
		return -1, err
	}
	count := 0
	for _, group := range groups {
		count += group.Number * group.Alignment()
	}
	return count, nil
}
