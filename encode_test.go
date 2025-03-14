package gocsv_test

import (
	"bytes"
	"testing"

	"github.com/dhrodao/gocsv"
	"github.com/stretchr/testify/assert"
)

// This test validates a basic CSV decoding
func TestEncodeWithHeaderAndCarriageReturn(t *testing.T) {
	input := "name,age\r\n\"John, Michael\",25\r\nJane,23\r\n"

	decoded := gocsv.Document{}
	*decoded.Header() = []string{"name", "age"}
	*decoded.Data() = []gocsv.Record{
		{"John, Michael", 25},
		{"Jane", 23},
	}

	var buffer bytes.Buffer
	encoder := gocsv.NewEncoder(&buffer)
	encoder.CarriageReturn()
	assert.Nil(t, encoder.Encode(&decoded))
	assert.Equal(t, input, buffer.String())

	t.Logf("Encoded document: %v", buffer.String())
}

func TestEncodeWithQuotedString(t *testing.T) {
	input := "a,b,c,d\n1,2,3,4\n1.2,2.3,3.4,5.6\n\"\"\"this\"\" is a \"\"test\"\"\",\"\"\"test\"\"\",test,test\n"

	decoded := gocsv.Document{}
	*decoded.Data() = []gocsv.Record{
		{"a", "b", "c", "d"},
		{1, 2, 3, 4},
		{1.2, 2.3, 3.4, 5.6},
		{"\"this\" is a \"test\"", "\"test\"", "test", "test"},
	}

	var buffer bytes.Buffer
	assert.Nil(t, gocsv.NewEncoder(&buffer).Encode(&decoded))
	assert.Equal(t, input, buffer.String())

	t.Logf("Encoded document: %v", buffer.String())
}

func TestEncodeWithNewLineInQuotedString(t *testing.T) {
	input := "a,b,c,\"d\ne\"\n1,2,3,4\n"

	decoded := gocsv.Document{}
	*decoded.Data() = []gocsv.Record{
		{"a", "b", "c", "d\ne"},
		{1, 2, 3, 4},
	}

	var buffer bytes.Buffer
	assert.Nil(t, gocsv.NewEncoder(&buffer).Encode(&decoded))
	assert.Equal(t, input, buffer.String())

	t.Logf("Encoded document: %v", buffer.String())
}
