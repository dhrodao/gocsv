# gocsv

This program provides go services to Decode/Encode files in CSV format.

Features:
* [RFC 4180](https://datatracker.ietf.org/doc/html/rfc4180) compliant Decoder/Encoder
* mapping to strings, integers, floats and boolean values
* buffered Decoder/Encoder

Decoding snippet:
```Go
if err := gocsv.NewDecoder(reader).Decode(&records); err != nil {
    panic("error decoding csv")
}
```

Encoding snippet:
```Go
var buffer bytes.Buffer
if err := gocsv.NewEncoder(&buffer).Encode(&decoded); err != nil {
    panic("error encoding csv")
}
```
