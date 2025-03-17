package gocsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
)

type CSVWriter interface {
	Write(record []string) error
	WriteAll(records [][]string) error
}

// This interface defines a contract to define custom
// types that could be marshaled
type Marshaler interface {
	MarshalCSV() (string, error)
}

// This is the structure that holds the CSV Encoder data
type Encoder struct {
	writer CSVWriter
	err    error
	header []string
}

// This function encodes a 'Document' into a CSV file
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		writer: csv.NewWriter(w),
		err:    nil,
	}
}

// This function allows the user to customise the CSV writer
func (e *Encoder) SetWriter(create func() CSVWriter) {
	e.writer = create()
}

// This function returns the error of the CSV Encoder
func (e *Encoder) Error() error {
	return e.err
}

// This function writes a 'Document' to the CSV file
func (e *Encoder) Encode(records any) error {
	inValue, inType := getInValueAndType(records)
	if ensureInType(inType) != nil {
		return fmt.Errorf("encode: wrong type received (%s), expected to receive a slice", inType.Kind())
	}

	if inValue.Len() == 0 {
		return fmt.Errorf("encode: received an empty slice")
	}

	inInnerType := getInInnerType(inType)
	if err := ensureInInnerType(inInnerType); err != nil {
		return err
	}

	typeInfo, err := getTypeInfo(inInnerType)
	if err != nil {
		return err
	}

	e.header = make([]string, 0, len(typeInfo.fields))
	for _, field := range typeInfo.fields {
		e.header = append(e.header, field.fTag)
	}

	lines := make([][]string, 0, inValue.Len())
	for i := range inValue.Len() {
		record := inValue.Index(i)
		line := make([]string, 0, len(typeInfo.fields))
		for _, fieldInfo := range typeInfo.fields {
			val, err := toString(record.FieldByIndex(fieldInfo.index).Interface())
			if err != nil {
				return err
			}
			line = append(line, val)
		}
		lines = append(lines, line)
	}

	if err := e.encodeHeader(); err != nil {
		return err
	}

	if err := e.encodeContent(lines); err != nil {
		return err
	}

	return nil
}

func (e *Encoder) encodeHeader() error {
	return e.writer.Write(e.header)
}

func (e *Encoder) encodeContent(lines [][]string) error {
	return e.writer.WriteAll(lines)
}

// This function returns the in data structure value and type
func getInValueAndType(in any) (reflect.Value, reflect.Type) {
	val := reflect.ValueOf(in)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	return val, val.Type()
}

// This function returns 'true' if the in data structure
// has the correct type
func ensureInType(in reflect.Type) error {
	switch in.Kind() {
	case reflect.Slice:
		return nil
	default:
		return fmt.Errorf("decode: unexpected in type: %s", in.Kind())
	}
}

// This function returns the inner data structure type and value
func getInInnerType(in reflect.Type) reflect.Type {
	return in.Elem()
}

// This function checks if the inner type is correct
func ensureInInnerType(inner reflect.Type) error {
	switch inner.Kind() {
	case reflect.Struct:
		return nil
	default:
		return fmt.Errorf("decode: unexpected inner type: %s", inner.Kind())
	}
}
