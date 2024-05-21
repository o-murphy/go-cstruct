# PyStruct
## Go lang implementation of python's struct module

#### Interpret bytes as packed binary data


> [!TIP] 
> If you can't find something in `README.md` you can search references in 
> *[Python's struct documentation](https://docs.python.org/3/library/struct.html)*

> [!NOTE] 
> Some functionality not yet implemented, [details there](#not-yet-implemented)

> [!NOTE]
> Tested only with small values

### Contents
* [Installation](#installation)
* [Usage](#usage)
	* [Example](#example)
	* [Byte Order, Size, and Alignment](#byte-order-size-and-alignment)
	* [Format characters](#format-characters)
	* [Functions](#functions)
		* [func CalcSize](#func-calcsize)
		* [func Pack](#func-pack)
		* [func PackInto](#func-packinto)
		* [func Unpack](#func-unpack)
		* [func UnpackFrom](#func-unpackfrom)
		* [func IterUnpack](#func-iterunpack)
	* [Types](#types)
		* [Type PyStruct](#type-struct-1)
			* [func CalcSize](#func-calcsize-1)
			* [func Pack](#func-pack-1)
			* [func PackInto](#func-packinto-1)
			* [func Unpack](#func-unpack-1)
			* [func UnpackFrom](#func-unpackfrom-1)
			* [func IterUnpack](#func-iterunpack-1)


## Installation

``` bash
go get github.com/o-murphy/pystruct-go
```

## Usage
### Example
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
	byteArray2, err := pystruct.Pack(`<3sf`, intf...)
	if err == nil {
		fmt.Println(byteArray2)
	}

	if bytes.Equal(byteArray, byteArray2) {
		fmt.Println("Equal")
	}

	// or use struct
	s := pystruct.NewStruct(`<3sf`)
	byteArray, err := s.Pack(intf...)
	intf, err := s.Unpack(byteArray)

}
```

### Byte Order, Size, and Alignment
<details> 
  <summary>Details</summary>
   By default, C types are represented in the machine’s native format and byte order, and properly aligned by skipping pad bytes if necessary (according to the rules used by the C compiler). This behavior is chosen so that the bytes of a packed struct correspond exactly to the memory layout of the corresponding C struct. Whether to use native byte ordering and padding or standard formats depends on the application.
   
   Alternatively, the first character of the format string can be used to indicate the byte order, size and alignment of the packed data, according to the following table: 
</details>

|Character  |Byte order            |Size    |Alignment |
|:---------:|----------------------|--------|----------|
|@          |native                |native  |native    |
|=          |native                |standard|none      |
|<          |little-endian         |standard|none      |
|>          |big-endian            |standard|none      |
|!          |network (= big-endian)|standard|none      |

> [!TIP] 
> If the first character is not one of these, '@' is assumed.

> [!NOTE]
> Note The number 1023 (0x3ff in hexadecimal) has the following byte representations:
> ```
> 03 ff in big-endian (>)
> ff 03 in little-endian (<)
> ```
> [More info there...](https://docs.python.org/3/library/struct.html#struct-alignment)


### Format Characters
<details> 
  <summary>Details</summary>
	Format characters have the following meaning; the conversion between C and Python values should be obvious given their types. The ‘Standard size’ column refers to the size of the packed value in bytes when using standard size; that is, when the format string starts with one of '<', '>', '!' or '='. When using native size, the size of the packed value is platform-dependent. 
</details>

> [!NOTE]
> [More info there...](https://docs.python.org/3/library/struct.html#format-characters)

| Format | C Type              | Go Type           | Python Type       | Standard Size |
| :----: | ------------------- | ----------------- | ----------------- | ------------- |
|   c    | char                | byte              | bytes of length 1 | 1             |
|   b    | signed char         | int8              | integer           | 1             |
|   B    | unsigned char       | uint8             | integer           | 1             |
|   ?    | _Bool               | bool              | bool              | 1             |
|   h    | short               | int16             | integer           | 2             |
|   H    | unsigned short      | uint16            | integer           | 2             |
|   i    | int                 | int32             | integer           | 4             |
|   I    | unsigned int        | uint32            | integer           | 4             |
|   l    | long                | int32             | integer           | 4             |
|   L    | unsigned long       | uint32            | integer           | 4             |
|   q    | long long           | int64             | integer           | 8             |
|   Q    | unsigned long long  | uint64            | integer           | 8             |
|   f    | float               | float32 (float64) | float             | 4 (8)         |
|   d    | double              | float64 (float32) | float             | 8 (4)         |
|   s    | char[]              | string            | bytes             | Variable      |
|   x    | pad byte            | [N/A**](#na)      | no value          |               |
|   n    | ssize_t             | [N/A**](#na)      | integer           |               |
|   N    | size_t              | [N/A**](#na)      | integer           |               |
|   e    | [float16](#float16) | [N/A**](#na)      | float             |  2             |
|   p    | char[]              | [N/A**](#na)      | bytes             |               |
|   P    | void*               | [N/A**](#na)      | integer           |               |

### Functions
#### func CalcSize
```go
func CalcSize(format string) (int, error)
```
Return the size of the struct
(and hence of the bytes object produced by pack(format, ...))
corresponding to the format string format

> ```go
> size, err := pystruct.CalcSize(format)
> if err == nil {
>	fmt.Println(size)
> }
> ```

#### func Pack
```go
func Pack(format string, intf []interface{}) ([]byte, error)
```
Return a bytes object containing the values v1, v2, … packed according to the format string format.
The arguments must match the values required by the format exactly.

> ```go
> intf := []interface{}{"abc", 1.01}
> byteArray, err := pystruct.Pack(`<3sf`, intf)
> if err == nil {
>	fmt.Println(byteArray2)
> }
> ```


#### func PackInto
```go
func PackInto(format string, buffer []byte, offset int, intf ...interface{}) ([]byte, error)
```
Pack the values v1, v2, … according to the format string format
and write the packed bytes into the writable buffer
starting at position offset. Note that offset is a required argument.

> ```go
> intf := []interface{}{"abc", 1.01}
> buffer := []byte{0xff, 0xff, 0xff, 0xff}
> byteArray, err := pystruct.PackInto(`<3sf`, buffer, 3, intf)
> if err == nil {
>	fmt.Println(byteArray)
> }
> ```

#### func Unpack
```go
func Unpack(format string, buffer []byte) ([]interface{}, error)
```
Unpack from the buffer buffer (presumably packed by Pack(format, ...))
according to the format string format. The result is an []interface{} even if it contains exactly one item.
The buffer’s size in bytes must match the size required by the format, as reflected by CalcSize().

> ```go
> byteArray := []byte{97, 98, 99, 100, 101, 102, 103}
> intf, err := pystruct.Unpack(`<3sf`, byteArray)
> if err == nil {
>	fmt.Println(intf...)
> }
> ```

#### func UnpackFrom
```go
func UnpackFrom(format string, buffer []byte, offset int) ([]interface{}, error)
```
Unpack from buffer starting at position offset, according to the format string format.
The result is an []interface{} even if it contains exactly one item.
The buffer’s size in bytes, starting at position offset,
must be at least the size required by the format, as reflected by CalcSize().

> ```go
> byteArray := []byte{0, 0, 0, 97, 98, 99, 100, 101, 102, 103}
> intf, err := pystruct.UnpackFrom(`<3sf`, byteArray, 3)
> if err == nil {
>	fmt.Println(intf...)
> }
> ```

#### func IterUnpack
```go
func IterUnpack(format string, buffer []byte) (<-chan interface{}, <-chan error)
```
Iteratively unpack from the buffer buffer according to the format string format.
This function returns an iterator which will read equally sized chunks from the buffer until all its contents have been consumed.
The buffer’s size in bytes must be a multiple of the size required by the format, as reflected by CalcSize()

<!-- // iterator, errs := pystruct.IterUnpack(format, byteArray)
// for i, value := 0, <-iterator; errs == nil; i, value = i+1, <-iterator {
// 	fmt.Println(i, value)
// } -->

> ```go
> format := `<3si`
> byteArray := []byte{97, 98, 99, 100, 101, 102, 103}
> iterator, errs := pystruct.IterUnpack(format, byteArray)
>
> i := 0
> for value := range iterator {
> 	fmt.Println(i, value)
>   i++
> }
> ```

### Types
#### type PyStruct
```go
type PyStruct struct {
	format string
}
```
PyStruct(fmt) --> compiled PyStruct object

##### func NewStruct
```go
NewStruct(format string) (PyStruct, error)
```
NewStruct(fmt) --> compiled PyStruct object
Methods bellow this just binds for same named functions

##### func CalcSize
([⬆️CalcSize](#func-pack))
```go
func (s *PyStruct) CalcSize() (int, error)
```

##### func Pack
([⬆️Pack](#func-pack))
```go
func (s *PyStruct) Pack(intf ...interface{}) ([]byte, error)
```

##### func PackInto
([⬆️PackInto](#func-packinto))
```go
func (s *PyStruct) PackInto(buffer []byte, offset int, intf ...interface{}) ([]byte, error) 
```

##### func Unpack
([⬆️Unpack](#func-unpack))
```go
func (s *PyStruct) Unpack(buffer []byte) ([]interface{}, error)
```

##### func UnpackFrom
([⬆️UnpackFrom](#func-unpackfrom))
```go
func (s *PyStruct) UnpackFrom(buffer []byte, offset int) ([]interface{}, error) 
```

##### func IterUnpack
([⬆️IterUnpack](#func-iterunpack))
```go
func (s *PyStruct) IterUnpack(format string, buffer []byte) (<-chan interface{}, <-chan error)
```

### Not yet implemented
##### N/A
* `p` char[] - Pascal string
* `P` void*
* `x` PadByte
* `n` ssize_t
* `N` size_t
##### Float16
* `e` Float16 - *(IEEE 754 binary16 half precision float)*

### RISK NOTICE
> [!IMPORTANT]
> THE CODE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE MATERIALS OR THE USE OR OTHER DEALINGS IN THE MATERIALS.