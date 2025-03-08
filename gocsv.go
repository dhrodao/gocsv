package gocsv

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
)

const (
	Separator      = ','
	StringWrapper  = "\""
	NewLine        = "\n"
	CarriageReturn = "\r"
)

// Regular expression to match integers
var reInteger = regexp.MustCompile(`^[+-]?(\d)+$`)

// Regular expression to match floats
var reFloat = regexp.MustCompile(`^[+-]?([0-9]*[.])?[0-9]+$`)

// Regular expression to match booleans
var reBool = regexp.MustCompile(`^(true|false)$`)

var ErrHeaderRowLenMismatch = errors.New("header and row length mismatch")
var ErrLineEmpty = errors.New("empty line")
var ErrEof = errors.New("EOF")
var ErrUnsuportedType = errors.New("unsupported type")

// This is the structure that holds the CSV Document data
// If the CSV file has a header, the 'headerValues' map will
// contain the header values as keys and the values of each
// column as a slice of any.
// The 'values' slice will contain the data of the CSV file
type Document struct {
	headerValues map[string][]any
	values       [][]any
}

// This function returns the header values of the CSV Document
func (d *Document) Header() map[string][]any {
	return d.headerValues
}

// This function returns the values of the CSV Document
func (d *Document) Data() [][]any {
	return d.values
}

// The Decoder type used to parse a *.csv file
type Decoder struct {
	scanner      bufio.Scanner
	sep          rune
	headerValues []string
	currentLine  int
	err          error
}

// This function creates a CSV Decoder and returns it
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		scanner:      *bufio.NewScanner(r),
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

// This function toggles the header parsing for a CSV document
func (d *Decoder) Header(v bool) {
	if v {
		d.headerValues = make([]string, 0)
	} else {
		d.headerValues = nil
	}
}

// This function returns the error of the CSV Decoder
func (d *Decoder) Error() error {
	return d.err
}

// This function reads from the buffered input and writes
// the decoded data to the CSV Document
func (d *Decoder) Decode(csv *Document) error {
	// Parse the header if needed
	if d.headerValues != nil {
		d.parseHeader()
	}
	for {
		line, err := d.readLine()
		d.err = err
		if err == ErrEof {
			break
		}
		if err == ErrLineEmpty {
			continue
		}
		// TODO handle: RFC 4180: If \n between double quotes, it is part of the field
		d.unmarshal(line, csv)
	}
	return nil
}

// This function reads a line if possible and returns it
func (d *Decoder) readLine() (string, error) {
	if d.scanner.Scan() {
		d.currentLine++
		line := d.scanner.Text()
		// ignore empty lines or
		if len(line) == 0 || line == NewLine {
			return "", ErrLineEmpty
		}
		// Remove every carriage return (avoid \r\n)
		line = strings.ReplaceAll(line, CarriageReturn, "")
		return line, nil
	}
	return "", ErrEof
}

// This function parses a header of a CSV document
func (d *Decoder) parseHeader() error {
	line, err := d.readLine()
	if err != nil {
		return err
	}

	d.headerValues = strings.Split(line, string(d.sep))

	return nil
}

// This function unmarshals a line into the CSV Document
func (d *Decoder) unmarshal(line string, doc *Document) error {
	// Split the line into tokens
	tokens := strings.Split(line, string(d.sep))

	// Combine the tokens between double quotes into a single token
	combinedTokens := make([]string, 0, len(tokens))
	var merged string
	for _, v := range tokens {
		switch true {
		case len(v) == 1 && strings.HasPrefix(v, StringWrapper):
			// (1) .. ,",", .. (2) .. ," text,", ..
			if merged != "" {
				merged += string(d.sep)
				combinedTokens = append(combinedTokens, merged)
				merged = ""
			} else {
				merged += string(d.sep)
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
			merged = strings.Join([]string{merged, v[:len(v)-1]}, string(d.sep))
			combinedTokens = append(combinedTokens, merged)
			merged = ""
		default:
			// (1) .. " , text , text , " .. (middle part)
			if merged != "" {
				merged = strings.Join([]string{merged, v}, string(d.sep))
			} else {
				combinedTokens = append(combinedTokens, v)
			}
		}
	}

	tokens = combinedTokens

	// Match each token to its type
	lineData := make([]any, 0, len(tokens))
	for _, cell := range tokens {
		switch {
		case reInteger.MatchString(cell):
			v, _ := strconv.Atoi(cell)
			lineData = append(lineData, v)
		case reFloat.MatchString(cell):
			v, _ := strconv.ParseFloat(cell, 64)
			lineData = append(lineData, v)
		case reBool.MatchString(cell):
			v, _ := strconv.ParseBool(cell)
			lineData = append(lineData, v)
		default:
			// By default consider it a string
			lineData = append(lineData, cell)
		}
	}

	// If the document has a header, add the values to the map
	if d.headerValues != nil {
		if len(lineData) != len(d.headerValues) {
			d.err = ErrHeaderRowLenMismatch
			return d.err
		}

		if doc.headerValues == nil {
			doc.headerValues = make(map[string][]any)
		}

		for i, v := range d.headerValues {
			if doc.headerValues[v] == nil {
				doc.headerValues[v] = make([]any, 0)
			}
			doc.headerValues[v] = append(doc.headerValues[v], lineData[i])
		}
	}

	doc.values = append(doc.values, lineData)

	return nil
}

// This is the structure that holds the CSV Encoder data
type Encoder struct {
	writer   bufio.Writer
	sep      rune
	carriage bool
	err      error
}

// This function encodes a 'Document' into a CSV file
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		writer: *bufio.NewWriter(w),
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
	e.carriage = true
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

	header := make([]any, 0, len(doc.values[0]))
	for k := range doc.headerValues {
		header = append(header, k)
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
func (e *Encoder) marshall(records []any) error {
	values := make([]string, 0, len(records))
	for _, record := range records {
		switch record := record.(type) {
		case int:
			values = append(values, strconv.Itoa(record))
		case float64:
			values = append(values, strconv.FormatFloat(record, 'f', -1, 64))
		case bool:
			values = append(values, strconv.FormatBool(record))
		case string:
			// Replace every double quote with two double quotes
			record = strings.ReplaceAll(record, StringWrapper, "\"\"")
			// If the line contains the separator or double quotes
			// wrap it in double quotes
			if strings.Contains(record, string(e.sep)) ||
				strings.Contains(record, "\"\"") {
				record = StringWrapper + record + StringWrapper
			}
			values = append(values, record)
		default:
			return ErrUnsuportedType
		}
	}

	// Write the line to the file
	line := strings.Join(values, string(e.sep))
	// Add a carriage return to the end of the line
	if e.carriage {
		line += CarriageReturn
	}
	// Add a new line to the end of the line
	line += NewLine

	if _, err := e.writer.WriteString(line); err != nil {
		return err
	}
	return e.writer.Flush()
}
