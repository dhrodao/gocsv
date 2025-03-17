package gocsv_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/dhrodao/gocsv"
	"github.com/stretchr/testify/assert"
)

func TestEncodeWithCarriageReturnAndQuotedString(t *testing.T) {
	expected := "Name,Age\r\n\"John, Francis\",25\r\nMichael,43\r\n"
	decoded := []A{
		{Name: "John, Francis", Age: 25},
		{Name: "Michael", Age: 43},
	}

	var buffer bytes.Buffer
	encoder := gocsv.NewEncoder(&buffer)
	encoder.CarriageReturn()
	assert.Nil(t, encoder.Encode(decoded))

	assert.Equal(t, expected, buffer.String())
}

func TestEncodeWithNewLineInQuotedString(t *testing.T) {
	expected := "Name,Age\r\n\"John, Francis\",25\r\n\"Michael\nnewline\",43\r\n"
	decoded := []A{
		{Name: "John, Francis", Age: 25},
		{Name: "Michael\nnewline", Age: 43},
	}

	var buffer bytes.Buffer
	encoder := gocsv.NewEncoder(&buffer)
	encoder.CarriageReturn()
	assert.Nil(t, encoder.Encode(decoded))

	assert.Equal(t, expected, buffer.String())
}

func TestEncodeWithCustomType(t *testing.T) {
	expected := "Name,Age,birthdate\nJohn,25,19990112\nMichael,50,19750101\n"
	decoded := []C{
		{A{"John", 25}, BirthDate{time.Date(1999, 12, 1, 0, 0, 0, 0, time.UTC)}, B{}},
		{A{"Michael", 50}, BirthDate{time.Date(1975, 1, 1, 0, 0, 0, 0, time.UTC)}, B{}},
	}

	var buffer bytes.Buffer
	assert.Nil(t, gocsv.NewEncoder(&buffer).Encode(decoded))

	assert.Equal(t, expected, buffer.String())
}
