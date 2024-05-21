package pystruct_test

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"unicode"

	pystruct "github.com/o-murphy/pystruct-go"
)

func TestCalcSize(t *testing.T) {

	v, ok := pystruct.CFormatMap['f']

	if !ok && !unicode.IsLetter('f') {
		t.Errorf("Is format, ok %c, %v", v, ok)
	}

	v, ok = pystruct.CFormatMap['-']

	if ok && unicode.IsLetter('-') {
		t.Errorf("Is not format, ok %c, %v", v, ok)
	}

	size, err := pystruct.CalcSize(`<3sf`)
	if err != nil {
		t.Errorf("Error occured: %s\n", err)
	} else {
		if size != 7 {
			t.Errorf("Expected: 7, got %d\n", size)
		}
	}

	fmt.Println("PASS: TestCalcSize")
}

func TestUnpack(t *testing.T) {
	byteArray := []byte{97, 98, 99, 100, 101, 102, 103}

	intf, err := pystruct.Unpack(`<3sf`, byteArray)
	if err != nil {
		t.Error("Unbound error:", err)
	}

	if len(intf) != 2 {
		t.Errorf("Expected 2 elements in the interface slice, got %d", len(intf))
	}

	if intf[0] != "abc" {
		t.Errorf("Wrong parsed value 0")
	}

	// Check the type of the second element in the interface slice
	if f, ok := intf[1].(float32); !ok {
		t.Errorf("Second element is not a float32: %f\n", f)
	}

	fmt.Println("PASS: TestUnpack")
}

func TestIterUnpack(t *testing.T) {
	format := `<3si`
	byteArray := []byte{97, 98, 99, 100, 101, 102, 103}

	// iterator, errs := pystruct.IterUnpack(format, byteArray)
	// for i, value := 0, <-iterator; errs == nil; i, value = i+1, <-iterator {
	// 	fmt.Println(i, value)
	// }

	iterator, errs := pystruct.IterUnpack(format, byteArray)

	i := 0
	for value := range iterator {
		switch i {
		case 0:
			switch v := value.(type) {
			case string:
				if v != "abc" {
					t.Errorf("Wring string")
				}
			default:
				t.Errorf("Unexpected type")
			}
		case 1:
			switch v := value.(type) {
			case int8:
				if v != 101 {
					t.Errorf("Wring string")
				}
			}
		}

		i += 1
	}

	for err := range errs {
		t.Error("Unbound error:", err)
	}

	fmt.Println("PASS: TestIterUnpack")
}

func TestUnpackFrom(t *testing.T) {

	byteArray := []byte{0, 0, 0, 97, 98, 99, 100, 101, 102, 103}

	intf, err := pystruct.UnpackFrom(`<3sf`, byteArray, 3)

	if err != nil {
		t.Error("Unbound error:", err)
	}

	if len(intf) != 2 {
		t.Errorf("Expected 2 elements in the interface slice, got %d", len(intf))
	}

	if intf[0] != "abc" {
		t.Errorf("Wrong parsed value %v", intf[0])
	}

	if f, ok := intf[1].(float32); !ok {
		t.Errorf("Second element is not a float32: %f\n", f)
	}

	fmt.Println("PASS: TestUnpackFrom")
}

func TestPack(t *testing.T) {
	// intf := []interface{}{"abc", 1.01, 3}
	intf := []interface{}{"abc", 1.01}
	// intf := []interface{}{"abc"}
	byteArray, err := pystruct.Pack("<3sf", intf...)
	expected := []byte{97, 98, 99, 174, 71, 129, 63}

	if err != nil {
		t.Error("Unbound error:", err)
	}

	if !bytes.Equal(byteArray, expected) {
		t.Errorf("Expected: %v\nActual: %v\n", expected, byteArray)
	}

	fmt.Println("PASS: TestPack")
}

func TestPackInto(t *testing.T) {
	intf := []interface{}{"abc", 1.01}
	buffer := []byte{0xff, 0xff, 0xff, 0xff}
	byteArray, err := pystruct.PackInto("<3sf", buffer, 2, intf...)

	expected := []byte{0xff, 0xff, 97, 98, 99, 174, 71, 129, 63}

	if err != nil {
		t.Error("Unbound error:", err)
	}

	if !bytes.Equal(byteArray, expected) {
		t.Errorf("Expected: %v\nActual: %v\n", expected, byteArray)
	}

	fmt.Println("PASS: TestPackInto")
}

func TestWrongOrder(t *testing.T) {
	byteArray := []byte{97, 98, 99, 100, 101, 102, 103}
	_, err := pystruct.Unpack(`3sf`, byteArray)

	if err != nil {
		t.Error("Get order error:", err)
	}

	_, err = pystruct.Unpack(`3<sf`, byteArray)
	if err == nil {
		t.Error("Order pos error:", err)
	}
	fmt.Println("PASS: TestWrongOrder")
}

func TestRegexp(t *testing.T) {
	// Define the regex pattern with capture groups for prefix and groups
	pattern := `^([@<>=!])?((\d*[cbBhHiIqQlLfds])+)$`

	// Compile the regular expression
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}

	// Test string
	testString := "@3s2bd"

	// Find the entire match with submatches
	matches := re.FindStringSubmatch(testString)
	if len(matches) == 0 {
		fmt.Println("No matches found")
		return
	}

	// Print the entire match
	fmt.Println("Entire match:", matches[0])

	// Extract and print the prefix if present
	prefix := matches[1]
	if prefix == "" {
		fmt.Println("Prefix: (none)")
	} else {
		fmt.Println("Prefix:", prefix)
	}

	// Extract the groups
	groups := matches[2]

	// Define the regex for individual groups with separate capture for number and char
	groupPattern := `(\d*)([cbBhHiIqQlLfds])`
	groupRe, err := regexp.Compile(groupPattern)
	if err != nil {
		fmt.Println("Error compiling group regex:", err)
		return
	}

	// Find all individual groups
	individualMatches := groupRe.FindAllStringSubmatch(groups, -1)

	// Print each group with parsed number and char
	for _, match := range individualMatches {
		numberStr := match[1]
		char := match[2]
		var number int
		if numberStr == "" {
			number = 1
		} else {
			number, err = strconv.Atoi(numberStr)
			if err != nil {
				fmt.Println("Error converting number:", err)
				continue
			}
		}
		fmt.Printf("Group: %s, Number: %d, Char: %s\n", match[0], number, char)
	}
}
