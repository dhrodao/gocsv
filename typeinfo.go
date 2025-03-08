package gocsv

import (
	"fmt"
	"reflect"
)

// The thag that the elements of the struct may contain
var tagName string = "csv"

// This structure will contain every field from a type
type typeInfo struct {
	parentType reflect.Type
	fields     []fieldInfo
}

// This type will contain the information of a given type
type fieldInfo struct {
	index []int
	fName string
	fTag  string
}

// This variable holds the Unmarshaler interface type
var unmarshalerType reflect.Type = reflect.TypeOf((*Unmarshaler)(nil)).Elem()

func getTypeInfo(t reflect.Type) (*typeInfo, error) {
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %s (%s) is not a struct", t.String(), t.Kind())
	}

	if t.NumField() == 0 {
		return nil, fmt.Errorf("type %s has no fields", t.String())
	}

	tInfo := typeInfo{parentType: t, fields: make([]fieldInfo, 0, t.NumField())}
	for i := range t.NumField() {
		tField := t.Field(i)
		// If non exported field ignore
		if (tField.PkgPath != "" && !tField.Anonymous) || tField.Tag.Get(tagName) == "-" {
			continue
		}
		fKind := tField.Type.Kind()
		// If embedded struct extract its fields
		if fKind == reflect.Struct {
			// Check if the struct implements the Unmarshaler interface
			if reflect.New(t.FieldByIndex(tField.Index).Type).Type().Implements(unmarshalerType) {
				goto INSERT
			}
			embeddedInfo, err := getTypeInfo(tField.Type)
			if err != nil {
				return nil, err
			}
			for _, newFieldInfo := range embeddedInfo.fields {
				newFieldInfo.index = append([]int{i}, newFieldInfo.index...)
				if err := addFieldInfo(tField.Type, &tInfo, &newFieldInfo); err != nil {
					return nil, err
				}
			}
			continue
		}
	INSERT:
		fInfo, err := getStructFieldInfo(tField)
		if err != nil {
			return nil, err
		}
		tInfo.fields = append(tInfo.fields, *fInfo)
	}

	return &tInfo, nil
}

func addFieldInfo(t reflect.Type, tInfo *typeInfo, newField *fieldInfo) error {
	for _, field := range tInfo.fields {
		// Return the first error
		if field.fName == newField.fName {
			return fmt.Errorf("field %s (tag: %s) conflicts with %s (%s)",
				field.fName, t.FieldByIndex(field.index).Tag.Get(tagName),
				newField.fName, t.FieldByIndex(newField.index).Tag.Get(tagName))
		}
	}

	tInfo.fields = append(tInfo.fields, *newField)
	return nil
}

func getStructFieldInfo(f reflect.StructField) (*fieldInfo, error) {
	fInfo := &fieldInfo{
		index: f.Index, fName: f.Name, fTag: f.Tag.Get(tagName)}

	return fInfo, nil
}
