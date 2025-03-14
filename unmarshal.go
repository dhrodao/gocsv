package gocsv

import "strconv"

// This type defines an interface which implements the unmarshaling function
type UnMarshaler interface {
	UnMarshalCSV() (Record, error)
}

// This type represents a raw CSV record
type RawRecord []string

func (r RawRecord) UnMarshalCSV() (Record, error) {
	// Match each token to its type
	lineData := make([]any, 0, len(r))
	for _, cell := range r {
		switch {
		case reInteger.MatchString(cell):
			v, err := strconv.Atoi(cell)
			if err != nil {
				return nil, err
			}
			lineData = append(lineData, v)
		case reFloat.MatchString(cell):
			v, err := strconv.ParseFloat(cell, 64)
			if err != nil {
				return nil, err
			}
			lineData = append(lineData, v)
		case reBool.MatchString(cell):
			v, err := strconv.ParseBool(cell)
			if err != nil {
				return nil, err
			}
			lineData = append(lineData, v)
		default:
			// By default consider it a string
			lineData = append(lineData, cell)
		}
	}
	return lineData, nil
}
