# gocsv

This program provides go services to Decode/Encode files in CSV format.

Features:
* [RFC 4180](https://datatracker.ietf.org/doc/html/rfc4180) compliant Decoder/Encoder
* mapping to strings, integers, floats and boolean values
* buffered Decoder/Encoder
* support Marshal/Unmarshal custom structures

Decoding snippet:
```Go
if err := gocsv.NewDecoder(reader).Decode(&records); err != nil {
	...
}
```

Encoding snippet:
```Go
var buffer bytes.Buffer
if err := gocsv.NewEncoder(&buffer).Encode(&decoded); err != nil {
	...
}
```

Custom Marshal/Unmarshal Structures:
```Go
// Custom Structure to Marshal/Unmarshal
type BirthDate struct {
	time.Time
}

// Implement Marshaler interface
func (d *BirthDate) MarshalCSV() (string, error) {
	return d.Time.Format("20060201"), nil
}

// Implement Unmarshaler interface
func (d *BirthDate) UnmarshalCSV(str string) (err error) {
	d.Time, err = time.Parse("20060201", str)
	return err
}

type C struct {
	Name      string    `csv:"Name"`
	Age       string    `csv:"Age`
	BirthDate BirthDate `csv:"birthdate"`
}
```
