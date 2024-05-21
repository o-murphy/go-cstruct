package main

import (
	"bytes"
	"fmt"

	pystruct "github.com/o-murphy/pystruct-go"
)

func main() {

	// unpack
	byteArray := []byte{97, 98, 99, 100, 101, 102, 103}
	intf, err := pystruct.Unpack(`<3sf`, byteArray)
	if err == nil {
		fmt.Println(intf...)
	} else {
		fmt.Println(err)
	}

	// pack
	byteArray2, err := pystruct.Pack(`<3sf`, intf...)
	if err == nil {
		fmt.Println(byteArray2)
	} else {
		fmt.Println(err)
	}

	if bytes.Equal(byteArray, byteArray2) {
		fmt.Println("Equal")
	} else {
		fmt.Println(err)
	}

	// or use struct
	s, err := pystruct.NewStruct(`<3sf`)
	if err == nil {
		byteArray, err := s.Pack(intf...)
		if err != nil {
			fmt.Println(intf...)
		} else {
			intf, err := s.Unpack(byteArray)
			if err != nil {
				fmt.Println(intf...)
			}
		}
	} else {
		fmt.Println(err)
	}

}
