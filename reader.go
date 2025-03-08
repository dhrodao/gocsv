package gocsv

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

var ErrEof = errors.New("EOF")

// This type represents the CSVReader
type CSVReader struct {
	s       *bufio.Scanner
	comment rune
	sep     rune
}

// This function returns a CSVReader
func NewCSVReader(r io.Reader) CSVReader {
	return CSVReader{
		s:       bufio.NewScanner(r),
		comment: Comment,
		sep:     Separator,
	}
}

// This function implements the logic to scan a file
func (s *CSVReader) Read() ([]string, error) {
	for s.s.Scan() {
		line := s.s.Text()

		// Ignore empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, string(s.comment)) {
			continue
		}

		// Remove every carriage return (avoid \r\n)
		line = strings.ReplaceAll(line, CarriageReturn, "")

		// TODO: check if token includes NewLine, this way it will be needed to read
		// another line

		return strings.Split(s.s.Text(), string(s.sep)), nil
	}
	return nil, ErrEof
}

// This function reads every line ignoring empty lines
// and comments (# ...). It also removes carriage returns
// at end of lines if exist
func (s *CSVReader) ReadAll() ([][]string, error) {
	lines := make([][]string, 0)
	for s.s.Scan() {
		line := s.s.Text()

		// Ignore empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, string(s.comment)) {
			continue
		}
		// Remove every carriage return (avoid \r\n)
		line = strings.ReplaceAll(line, CarriageReturn, "")

		// Split the line into tokens
		tokens := strings.Split(line, string(s.sep))

		// Combine the tokens between double quotes into a single token
		combinedTokens := make([]string, 0, len(tokens))
		var merged string
		for _, v := range tokens {
			switch true {
			case len(v) == 1 && strings.HasPrefix(v, StringWrapper):
				// (1) .. ,",", .. (2) .. ," text,", ..
				if merged != "" {
					merged += string(s.sep)
					combinedTokens = append(combinedTokens, merged)
					merged = ""
				} else {
					merged += string(s.sep)
				}
			case len(v) >= 2 && strings.HasPrefix(v, StringWrapper) && strings.HasSuffix(v, StringWrapper):
				// (1) .. ," text ", .. (2) .. ,"", ..
				combinedTokens = append(combinedTokens, strings.ReplaceAll(v[1:len(v)-1], "\"\"", StringWrapper))
				merged = ""
			case strings.HasPrefix(v, StringWrapper):
				// (1) .. , " text , text " , .. (1st part)
				merged = v[1:]
			case strings.HasSuffix(v, StringWrapper):
				// (1) .. , " text , text" , .. (2nd part)
				// This token is part of the previous line token
				if merged == "" {
					prevToken := lines[len(lines)-1][len(lines[0])-1]
					lines[len(lines)-1][len(lines[0])-1] = (prevToken + v[:len(v)-1])
				} else {
					merged = strings.Join([]string{merged, v[:len(v)-1]}, string(s.sep))
					combinedTokens = append(combinedTokens, merged)
					merged = ""
				}
			default:
				// (1) .. " , text , text , " .. (middle part)
				if merged != "" {
					merged = strings.Join([]string{merged, v}, string(s.sep))
				} else {
					combinedTokens = append(combinedTokens, v)
				}
			}
		}

		// (1) .. "text\ntext", ..
		// This token is incomplete, so we need to wait for the next line
		if merged != "" {
			combinedTokens = append(combinedTokens, merged+NewLine)
		}

		if len(combinedTokens) > 0 {
			lines = append(lines, combinedTokens)
		}
	}
	return lines, nil
}
