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
