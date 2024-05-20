package pystruct_test

import (
	"fmt"
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
	format := `<3sf`
	byteArray := []byte{97, 98, 99, 100, 101, 102, 103}

	iterator, errs := pystruct.IterUnpack(format, byteArray)

	i := 0
	for value := range iterator {
		switch v := value.(type) {
		case string:
			if v != "abc" {
				t.Errorf("Wring string")
			}
		case float32:

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
		t.Errorf("Wrong parsed value 0")
	}

	if f, ok := intf[1].(float32); !ok {
		t.Errorf("Second element is not a float32: %f\n", f)
	}

	fmt.Println("PASS: TestUnpackFrom")
}

func TestPack(t *testing.T) {
	intf := []interface{}{"abc", 1.01}
	byteArray, err := pystruct.Pack("<3sf", intf)

	if err != nil {
		t.Error("Unbound error:", err)
	} else {
		fmt.Println("byteArray:", byteArray)
	}
}
