package gocsv_test

import (
	"strings"
	"testing"
	"time"

	"github.com/dhrodao/gocsv"
	"github.com/stretchr/testify/assert"
)

type A struct {
	Name string `csv:"Name"`
	Age  int64  `csv:"Age"`
}

func TestDecode(t *testing.T) {
	input := "John,25\nMichael,50"
	reader := strings.NewReader(input)

	var records []*A

	assert.Nil(t, gocsv.NewDecoder(reader).Decode(&records))

	expected := []A{
		{"John", 25},
		{"Michael", 50},
	}

	for i, v := range records {
		assert.EqualValues(t, expected[i], *v)
	}
}

func TestDecodeWithHeaderAndCarriageReturn(t *testing.T) {
	input := "name,age\r\n\"John, Michael\",25\r\nJane,23\r\n"
	reader := strings.NewReader(input)

	decoder := gocsv.NewDecoder(reader)
	decoder.ContainsHeader(true)
	var records []A
	assert.Nil(t, decoder.Decode(&records))

	expected := []A{
		{"John, Michael", 25},
		{"Jane", 23},
	}

	for i, v := range records {
		assert.EqualValues(t, expected[i], v)
	}
}

type B struct {
	A string `csv:"a"`
	B string `csv:"b"`
	C string `csv:"c"`
	D string `csv:"d"`
}

func TestDecodeWithQuotedString(t *testing.T) {
	input := "a,b,c,d\n1,2,3,4\n1.2,2.3,3.4,5.6\n\"\"\"this\"\" is a \"\"test\"\"\",\"\"\"test\"\"\",test,test\n"
	reader := strings.NewReader(input)

	var records []B
	assert.Nil(t, gocsv.NewDecoder(reader).Decode(&records))

	expected := []B{
		{"a", "b", "c", "d"},
		{"1", "2", "3", "4"},
		{"1.2", "2.3", "3.4", "5.6"},
		{"\"this\" is a \"test\"", "\"test\"", "test", "test"},
	}

	for i, v := range records {
		assert.EqualValues(t, expected[i], v)
	}
}

func TestDecodeWithNewLineInQuotedString(t *testing.T) {
	input := "a,b,c,\"d\ne\"\n1,2,3,4\n"
	reader := strings.NewReader(input)

	var records []B
	assert.Nil(t, gocsv.NewDecoder(reader).Decode(&records))

	expected := []B{
		{"a", "b", "c", "d\ne"},
		{"1", "2", "3", "4"},
	}

	for i, v := range records {
		assert.EqualValues(t, expected[i], v)
	}
}

type BirthDate struct {
	time.Time
}

func (d *BirthDate) UnmarshalCSV(str string) (err error) {
	d.Time, err = time.Parse("20060201", str)
	return err
}

func (d *BirthDate) MarshalCSV() (string, error) {
	return d.Time.Format("20060201"), nil
}

type C struct {
	A         A
	BirthDate BirthDate `csv:"birthdate"`
	B         B         `csv:"-"` // This field will be ignored
}

func TestDecodeWithCustomType(t *testing.T) {
	input := "John,25,19990112\nMichael,50,19750101"
	reader := strings.NewReader(input)

	var records []*C

	assert.Nil(t, gocsv.NewDecoder(reader).Decode(&records))

	expected := []C{
		{A{"John", 25}, BirthDate{time.Date(1999, 12, 1, 0, 0, 0, 0, time.UTC)}, B{}},
		{A{"Michael", 50}, BirthDate{time.Date(1975, 1, 1, 0, 0, 0, 0, time.UTC)}, B{}},
	}

	for i, v := range records {
		assert.EqualValues(t, expected[i], *v)
	}
}
