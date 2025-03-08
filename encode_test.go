package gocsv_test

import (
	"bytes"
	"testing"

	"github.com/dhrodao/gocsv"
	"github.com/stretchr/testify/assert"
)

func TestEncode2WithCarriageReturnAndQuotedString(t *testing.T) {
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

func TestEncode2WithNewLineInQuotedString(t *testing.T) {
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
