package gocsv

import (
	"bufio"
	"io"
)

// This type represents a reader for CSV files
type CSVWriter struct {
	w        *bufio.Writer
	carriage bool
}

// This function returns a CSVWriter
func NewCSVWriter(w io.Writer) CSVWriter {
	return CSVWriter{
		w:        bufio.NewWriter(w),
		carriage: false,
	}
}

// This function writes the given string
func (w *CSVWriter) WriteString(s string) error {
	if w.carriage {
		s += CarriageReturn
	}
	s += NewLine

	if _, err := w.w.WriteString(s); err != nil {
		return err
	}
	return w.w.Flush()
}
