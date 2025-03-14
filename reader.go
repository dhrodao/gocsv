package gocsv

import (
	"bufio"
	"io"
	"strings"
)

// This type represents the CSVReader
type CSVReader struct {
	s       *bufio.Scanner
	comment rune
}

// This function returns a CSVReader
func NewCSVReader(r io.Reader) CSVReader {
	return CSVReader{
		s:       bufio.NewScanner(r),
		comment: Comment,
	}
}

// This function implements the logic to scan a file
func (s *CSVReader) ReadLine() (string, error) {
	for s.s.Scan() {
		line := s.s.Text()
		// Ignore empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, string(s.comment)) {
			continue
		}
		// Remove every carriage return (avoid \r\n)
		line = strings.ReplaceAll(line, CarriageReturn, "")
		return s.s.Text(), nil
	}
	return "", ErrEof
}
