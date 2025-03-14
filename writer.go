package gocsv

import (
	"bufio"
	"io"
	"strings"
)

// This type represents a reader for CSV files
type CSVWriter struct {
	w        *bufio.Writer
	carriage bool
	sep      rune
}

// This function returns a CSVWriter
func NewCSVWriter(w io.Writer) CSVWriter {
	return CSVWriter{
		w:        bufio.NewWriter(w),
		carriage: false,
		sep:      Separator,
	}
}

// This function sets the writer separator
func (w *CSVWriter) Separator(sep rune) {
	w.sep = sep
}

// This function writes a line
func (w *CSVWriter) Write(line []string) error {
	w.w.WriteString(w.composeLine(line))
	return w.w.Flush()
}

// This function writes a slice of lines
func (w *CSVWriter) WriteAll(lines [][]string) error {
	for _, line := range lines {
		w.w.WriteString(w.composeLine(line))
	}
	return w.w.Flush()
}

// This function composes an array of CSV fields into a line
func (w *CSVWriter) composeLine(fields []string) string {
	values := make([]string, 0, len(fields))
	for _, v := range fields {
		// Replace every double quote with two double quotes
		v = strings.ReplaceAll(v, StringWrapper, "\"\"")
		// If the line contains the separator or double quotes
		// wrap it in double quotes
		if strings.Contains(v, string(w.sep)) ||
			strings.Contains(v, "\"\"") ||
			strings.Contains(v, NewLine) {
			v = StringWrapper + v + StringWrapper
		}
		values = append(values, v)
	}

	s := strings.Join(values, string(w.sep))
	if w.carriage {
		s += CarriageReturn
	}
	s += NewLine

	return s
}
