package main

import (
	"fmt"

	pystruct "github.com/o-murphy/pystruct-go"
)

func main() {

	// unpack
	byteArray := []byte{97, 98, 99, 100, 101, 102, 103}
	if intf, err := pystruct.Unpack(`<3sf`, byteArray); err == nil {
		fmt.Println(intf...)
	}

	// pack
	if byteArray2, err := pystruct.Pack(`<3sf`, intf); err == nil {
		fmt.Println(byteArray2)
	}

}
