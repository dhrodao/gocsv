package gocsv_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dhrodao/gocsv"
	"github.com/stretchr/testify/assert"
)

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

// This test validates a basic CSV decoding
func TestWithHeaderAndCarriageReturn(t *testing.T) {
	input := "name,age\r\n\"John, Michael\",25\r\nJane,23\r\n"
	reader := strings.NewReader(input)

	decoder := gocsv.NewDecoder(reader)
	decoder.Header(true)
	var decoded gocsv.Document
	assert.Nil(t, decoder.Decode(&decoded))

	// Check document header
	assert.Equal(t, len(decoded.Header()), 2)
	assert.NotNil(t, decoded.Header()["name"])
	assert.NotNil(t, decoded.Header()["age"])

	// Check document contents
	expected := [][]any{
		{"John, Michael", 25},
		{"Jane", 23},
	}
	for i, row := range decoded.Data() {
		assert.Equal(t, len(row), 2)
		for j, col := range row {
			assert.Equal(t, col, expected[i][j])
		}
	}

	t.Logf("Decoded document: %v", decoded)

	var buffer bytes.Buffer
	encoder := gocsv.NewEncoder(&buffer)
	encoder.CarriageReturn()
	assert.Nil(t, encoder.Encode(&decoded))
	assert.Equal(t, input, buffer.String())

	t.Logf("Encoded document: %v", buffer.String())
}

func TestWithQuotedString(t *testing.T) {
	input := "a,b,c,d\n1,2,3,4\n1.2,2.3,3.4,5.6\n\"\"\"this\"\" is a \"\"test\"\"\",\"\"\"test\"\"\",test,test\n"
	reader := strings.NewReader(input)

	decoder := gocsv.NewDecoder(reader)
	var decoded gocsv.Document
	assert.Nil(t, decoder.Decode(&decoded))

	// Check document header
	assert.Nil(t, decoded.Header())

	// Check document contents
	expected := [][]any{
		{"a", "b", "c", "d"},
		{1, 2, 3, 4},
		{1.2, 2.3, 3.4, 5.6},
		{"\"this\" is a \"test\"", "\"test\"", "test", "test"},
	}
	for i, row := range decoded.Data() {
		assert.Equal(t, len(row), len(expected[i]))
		for j, col := range row {
			assert.Equal(t, col, expected[i][j])
		}
	}

	t.Logf("Decoded document: %v", decoded)

	var buffer bytes.Buffer
	assert.Nil(t, gocsv.NewEncoder(&buffer).Encode(&decoded))
	assert.Equal(t, input, buffer.String())

	t.Logf("Encoded document: %v", buffer.String())
}
