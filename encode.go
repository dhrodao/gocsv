package gocsv

import (
	"io"
	"strings"
)

// This is the structure that holds the CSV Encoder data
type Encoder struct {
	writer CSVWriter
	sep    rune
	err    error
}

// This function encodes a 'Document' into a CSV file
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		writer: NewCSVWriter(w),
		sep:    Separator,
		err:    nil,
	}
}

// This function sets a separator string for the CSV Encoder
func (e *Encoder) Separator(s rune) {
	e.sep = s
}

// This function sets a carriage return for the CSV Encoder
func (e *Encoder) CarriageReturn() {
	e.writer.carriage = true
}

// This function returns the error of the CSV Encoder
func (e *Encoder) Error() error {
	return e.err
}

// This function writes a 'Document' to the CSV file
func (e *Encoder) Encode(doc *Document) error {
	if err := e.writeHeader(doc); err != nil {
		e.err = err
		return e.err
	}
	if err := e.writeRows(doc); err != nil {
		e.err = err
		return e.err
	}
	return nil
}

// This function writes the header of a CSV file
func (e *Encoder) writeHeader(doc *Document) error {
	if len(doc.headerValues) == 0 {
		return nil
	}

	header := make(Record, 0, len(doc.headerValues))
	for _, v := range doc.headerValues {
		header = append(header, v)
	}
	return e.marshall(header)
}

// This function writes the rows of a CSV file
func (e *Encoder) writeRows(doc *Document) error {
	for _, row := range doc.values {
		if err := e.marshall(row); err != nil {
			return err
		}
	}
	return nil
}

// This function marshalls a slice of 'any' into a CSV line
func (e *Encoder) marshall(record Marshaler) error {
	fields, _ := record.MarshallCSV()
	values := make([]string, 0, len(fields))
	for _, v := range fields {
		// Replace every double quote with two double quotes
		v = strings.ReplaceAll(v, StringWrapper, "\"\"")
		// If the line contains the separator or double quotes
		// wrap it in double quotes
		if strings.Contains(v, string(e.sep)) ||
			strings.Contains(v, "\"\"") ||
			strings.Contains(v, NewLine) {
			v = StringWrapper + v + StringWrapper
		}
		values = append(values, v)
	}

	// Write the line to the file
	line := strings.Join(values, string(e.sep))

	return e.writer.WriteString(line)
}
