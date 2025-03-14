package gocsv

import (
	"errors"
	"io"
	"strings"
)

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

// This function reads from the buffered input and writes
// the decoded data to the CSV Document
func (d *Decoder) Decode(csv *Document) error {
	// Decode the header if needed
	if d.containsHeader {
		if d.err = d.decodeHeader(); d.err != nil {
			return d.err
		}
	}

	for {
		var line string
		line, d.err = d.line()
		if errors.Is(d.err, ErrEof) {
			break
		}
		if errors.Is(d.err, ErrLineEmpty) {
			continue
		}
		if d.err = d.unmarshal(line, csv); d.err != nil {
			return d.err
		}
	}
	return nil
}

// This function reads a line if possible and returns it
func (d *Decoder) line() (string, error) {
	line, err := d.scanner.ReadLine()
	if err == nil {
		d.currentLine++
		return line, nil
	}
	return line, err
}

// This function decodes a header of a CSV document
func (d *Decoder) decodeHeader() error {
	if d.headerValues == nil {
		d.headerValues = make([]string, 0)
	}

	line, err := d.line()
	if err != nil {
		return err
	}

	d.headerValues = strings.Split(line, string(d.sep))
	if len(d.headerValues) == 0 {
		return ErrHeaderEmpty
	}
	return nil
}

// This function unmarshals a line into the CSV Document
func (d *Decoder) unmarshal(line string, doc *Document) error {
	// Split the line into tokens
	tokens := strings.Split(line, string(d.sep))

	// Combine the tokens between double quotes into a single token
	combinedTokens := make(RawRecord, 0, len(tokens))
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
			// This token is part of the previous line token
			if merged == "" {
				prevToken := doc.values[len(doc.values)-1][len(doc.values[len(doc.values)-1])-1]
				doc.values[len(doc.values)-1][len(doc.values[len(doc.values)-1])-1] = (prevToken.(string) + v[:len(v)-1])
			} else {
				merged = strings.Join([]string{merged, v[:len(v)-1]}, string(d.sep))
				combinedTokens = append(combinedTokens, merged)
				merged = ""
			}
		default:
			// (1) .. " , text , text , " .. (middle part)
			if merged != "" {
				merged = strings.Join([]string{merged, v}, string(d.sep))
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

	// Match each token to its type
	lineData, err := combinedTokens.UnMarshalCSV()
	if err != nil {
		return err
	}

	// If the document has a header, set the header slice
	if d.headerValues != nil {
		if len(lineData) != len(d.headerValues) {
			d.err = ErrHeaderRowLenMismatch
			return d.err
		}
		if doc.headerValues == nil {
			doc.headerValues = d.headerValues
		}
	}

	if len(lineData) > 0 {
		doc.values = append(doc.values, lineData)
	}

	return nil
}
