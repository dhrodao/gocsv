package gocsv

import (
	"errors"
	"regexp"
)

const (
	Comment        = '#'
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
var ErrHeaderEmpty = errors.New("empty header")

// This is the structure that holds the CSV Document data
// If the CSV file has a header, the 'headerValues' map will
// contain the header values as keys and the values of each
// column as a slice of any.
// The 'values' slice will contain the data of the CSV file
type Document struct {
	headerValues []string
	values       []Record
}

// This function returns the header values of the CSV Document
func (d *Document) Header() *[]string {
	return &d.headerValues
}

// This function returns the values of the CSV Document
func (d *Document) Data() *[]Record {
	return &d.values
}
