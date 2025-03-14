package gocsv

import "strconv"

// This type defines an interface which implements the marshaling function
type Marshaler interface {
	MarshallCSV() ([]string, error)
}

// This type represents a CSV record
type Record []any

// This function marshals a record
func (r Record) MarshallCSV() ([]string, error) {
	values := make([]string, 0, len(r))
	for _, field := range r {
		switch v := field.(type) {
		case int:
			values = append(values, strconv.Itoa(v))
		case float64:
			values = append(values, strconv.FormatFloat(v, 'f', -1, 64))
		case bool:
			values = append(values, strconv.FormatBool(v))
		case string:
			values = append(values, v)
		default:
			return nil, ErrUnsuportedType
		}
	}
	return values, nil
}
