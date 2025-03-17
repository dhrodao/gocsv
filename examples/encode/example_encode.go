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
