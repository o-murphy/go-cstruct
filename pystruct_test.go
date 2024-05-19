package pystruct_test

import (
	"fmt"
	"testing"

	"pystruct"
)

func TestCalcSize(t *testing.T) {
	size, err := pystruct.CalcSize(`<3sf`)
	if err != nil {
		t.Errorf("Error occured: %s\n", err)
	} else {
		fmt.Printf("Calculated size: %d\n", size)

		if size != 7 {
			t.Errorf("Expected: 7, got %d\n", size)
		}

	}
}

func TestUnpack(t *testing.T) {
	byteArray := []byte{97, 98, 99, 100, 101, 102, 103}

	intf, err := pystruct.Unpack(`<3sf`, byteArray)
	if err == nil {
		fmt.Println(intf...)
	} else {
		fmt.Println(err)
		t.Errorf("Unbound error")
	}

	if len(intf) != 2 {
		t.Errorf("Expected 2 elements in the interface slice, got %d", len(intf))
	}

	if intf[0] != "abc" {
		t.Errorf("Wrong parsed value 0")
	}

	// Check the type of the second element in the interface slice
	if f, ok := intf[1].(float32); ok {
		fmt.Printf("Second element is a float32: %f\n", f)
	} else {
		t.Errorf("Second element is not a float32")
	}
}

func TestIterUnpack(t *testing.T) {
	format := `<3sf`
	byteArray := []byte{97, 98, 99, 100, 101, 102, 103}

	iterator, errs := pystruct.IterUnpack(format, &byteArray)

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
		fmt.Println(value)
	}

	for err := range errs {
		fmt.Println("Error:", err)
	}
}

func TestUnpackFrom(t *testing.T) {

	byteArray := []byte{0, 0, 0, 97, 98, 99, 100, 101, 102, 103}

	intf, err := pystruct.UnpackFrom(`<3sf`, byteArray, 3)
	if err == nil {
		fmt.Println(intf...)
	} else {
		fmt.Println(err)
		t.Errorf("Unbound error")
	}

	if len(intf) != 2 {
		t.Errorf("Expected 2 elements in the interface slice, got %d", len(intf))
	}

	if intf[0] != "abc" {
		t.Errorf("Wrong parsed value 0")
	}

	// Check the type of the second element in the interface slice
	if f, ok := intf[1].(float32); ok {
		fmt.Printf("Second element is a float32: %f\n", f)
	} else {
		t.Errorf("Second element is not a float32")
	}
}
