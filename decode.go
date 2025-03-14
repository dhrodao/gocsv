package gocsv

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

var ErrHeaderEmpty = errors.New("empty header")

// This type will be used to Unmarshal custom types
// from the CSV
type Unmarshaler interface {
	UnmarshalCSV(str string) error
}

// The Decoder type used to decode a *.csv file
type Decoder struct {
	scanner        CSVReader
	sep            rune
	headerValues   []string
	containsHeader bool
	currentLine    int
	err            error
}

// This function creates a CSV Decoder and returns it
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		scanner:      NewCSVReader(r),
		sep:          Separator,
		headerValues: nil,
		currentLine:  0,
		err:          nil,
	}
}

// This function sets a separator string for the CSV Decoder
func (d *Decoder) Separator(s rune) {
	d.sep = s
}

func (d *Decoder) Comment(comment rune) {
	d.scanner.comment = comment
}

// This function toggles the header parsing for a CSV document
func (d *Decoder) ContainsHeader(v bool) {
	d.containsHeader = v
}

// This function returns the error of the CSV Decoder
func (d *Decoder) Error() error {
	return d.err
}

// This function decodes a CSV into the passed structure which
// should be a pointer to a slice of records
func (d *Decoder) Decode(out any) error {
	outVal, outType := getOutValueAndType(out)
	if err := ensureOutType(outType); err != nil {
		return err
	}

	wasInnerPointer, outInnerType := getOutInnerType(outType)
	if err := ensureOutInnerType(outInnerType); err != nil {
		return err
	}

	typeInfo, err := getTypeInfo(outInnerType)
	if err != nil {
		return err
	}

	if len(typeInfo.fields) == 0 {
		return errors.New("decode: expected fields to decode")
	}

	// Decode the header from the struct tags
	if !d.containsHeader {
		d.headerValues = make([]string, 0, len(typeInfo.fields))
		for _, v := range typeInfo.fields {
			d.headerValues = append(d.headerValues, v.fTag)
		}
	} else {
		// Decode header from the input
		if err := d.decodeHeader(); err != nil {
			return err
		}
	}

	lines, err := d.scanner.ReadAll()
	if err != nil {
		return err
	}

	if len(lines) == 0 {
		return errors.New("decode: empty CSV file")
	}

	if len(lines[0]) != len(d.headerValues) {
		return fmt.Errorf("decode: header len (%d) is not equal to content len (%d)", len(lines[0]), len(d.headerValues))
	}

	if err := ensureOutCapacity(outVal, len(lines)); err != nil {
		return err
	}

	for i, line := range lines {
		outInnerValue := getNewOutInnerValue(wasInnerPointer, outInnerType)
		for j, value := range line {
			oi := outInnerValue
			if wasInnerPointer {
				oi = outInnerValue.Elem()
			}
			if err := setValue(oi.FieldByIndex(typeInfo.fields[j].index), value); err != nil {
				return err
			}
		}
		outVal.Index(i).Set(outInnerValue)
	}

	return nil
}

// This function creates a new inner type value
func getNewOutInnerValue(wasInnerPointer bool, typ reflect.Type) reflect.Value {
	if wasInnerPointer {
		return reflect.New(typ)
	}
	return reflect.New(typ).Elem()
}

// This function returns the out value and type. It is
// expected to receive a slice or a pointer to a slice
func getOutValueAndType(out any) (reflect.Value, reflect.Type) {
	val := reflect.ValueOf(out)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	return val, val.Type()
}

// This function check the out type has the correct type
func ensureOutType(typ reflect.Type) error {
	switch typ.Kind() {
	case reflect.Slice:
		return nil
	default:
		return fmt.Errorf("decode: unexpected out type (%s)", typ.String())
	}
}

// This function extracts the inner type from the out type.
// The out type should be a slice of records
func getOutInnerType(out reflect.Type) (wasInnerPointer bool, innerType reflect.Type) {
	innerType = out.Elem()

	if innerType.Kind() == reflect.Pointer {
		wasInnerPointer = true
		innerType = innerType.Elem()
	}

	return wasInnerPointer, innerType
}

// This function ensures that the out inner type is the expected
func ensureOutInnerType(outInnerType reflect.Type) error {
	switch outInnerType.Kind() {
	case reflect.Struct:
		return nil
	default:
		return fmt.Errorf("decode: expected inner type to be struct (%s)", outInnerType.String())
	}
}

// This function ensures that the out value has enough capacity to
// fit every CSV record
func ensureOutCapacity(out reflect.Value, lenght int) error {
	switch out.Kind() {
	case reflect.Slice:
		if !out.CanAddr() && out.Len() < lenght {
			return fmt.Errorf("decode: out value is not addressable and it has not enough lenght (%d)", out.Len())
		} else {
			out.Set(reflect.MakeSlice(out.Type(), lenght, lenght))
		}
	}

	return nil
}

// This function decodes a header of a CSV document
func (d *Decoder) decodeHeader() error {
	if d.headerValues == nil {
		d.headerValues = make([]string, 0)
	}

	line, err := d.scanner.Read()
	if err != nil {
		return err
	}

	d.headerValues = line
	if len(d.headerValues) == 0 {
		return ErrHeaderEmpty
	}
	return nil
}
