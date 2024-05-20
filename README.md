# PyStruct
## Go lang implementation of python's struct module


### Installation

``` bash
> go get github.com/o-murphy/pystruct-go
```
 > [!TIP]
 > Tested only with small values


### Usage
```go
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
	}

	// pack
	byteArray2, err := pystruct.Pack(`<3sf`, intf)
	if err == nil {
		fmt.Println(byteArray2)
	}

	if bytes.Equal(byteArray, byteArray2) {
		fmt.Println("Equal")
	}

}
```