# gocsv

This program provides go services to Decode/Encode files in CSV format.

Features:
* [RFC 4180](https://datatracker.ietf.org/doc/html/rfc4180) compliant Decoder/Encoder
* mapping to strings, integers, floats and boolean values
* buffered Decoder/Encoder
* support Marshal/Unmarshal custom structures

Decoding snippet:
```go:examples/decode/example_decode.go
package main

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/dhrodao/gocsv"
)

type Person struct {
	Name string `csv:"Name"`
	Age  int64  `csv:"Age"`
}

func main() {
	input := "John;25\nMichael;50"
	reader := strings.NewReader(input)

	decoder := gocsv.NewDecoder(reader)
	decoder.SetReader(func() gocsv.CSVReader {
		r := csv.NewReader(reader)
		r.Comma = ';'
		return r
	})

	var records []*Person
	if err := decoder.Decode(&records); err != nil {
		panic("Error decoding!")
	}

	for _, record := range records {
		fmt.Println(record)
	}
}
```

Encoding snippet:
```go:examples/encode/example_encode.go
package main

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/dhrodao/gocsv"
)

type Person struct {
	Name string `csv:"Name"`
	Age  int64  `csv:"Age"`
}

func main() {
	var buffer bytes.Buffer
	decoder := gocsv.NewEncoder(&buffer)
	decoder.SetWriter(func() gocsv.CSVWriter {
		w := csv.NewWriter(&buffer)
		w.UseCRLF = true
		return w
	})

	var records []*Person
	if err := decoder.Encode(&records); err != nil {
		panic("Error decoding!")
	}

	fmt.Println(buffer.String())
}
```

Custom Marshal/Unmarshal Structures:
```Go
// Custom Structure to Marshal/Unmarshal
type BirthDate struct {
	time.Time
}

// Implement Marshaler interface
func (d *BirthDate) MarshalCSV() (string, error) {
	return d.Time.Format("20060201"), nil
}

// Implement Unmarshaler interface
func (d *BirthDate) UnmarshalCSV(str string) (err error) {
	d.Time, err = time.Parse("20060201", str)
	return err
}

type C struct {
	Name      string    `csv:"Name"`
	Age       string    `csv:"Age`
	BirthDate BirthDate `csv:"birthdate"`
}
```
