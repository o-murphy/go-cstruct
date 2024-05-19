package cstruct_test

import (
	"cstruct"
	"fmt"
	"testing"
)

func TestUnpack(t *testing.T) {
	byteArray := []byte{97, 98, 99, 100, 101, 102, 103}

	intf, err := cstruct.Unpack(`<3sf`, byteArray)
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

// func TestPack(t *testing.T) {
// 	dataToPack := []interface{}{"abc", 1.01}

// 	byteArray, err := cstruct.Pack("<3sf", dataToPack)

// 	if err == nil {
// 		fmt.Println(byteArray)
// 	} else {
// 		fmt.Println(err)
// 	}

// }
