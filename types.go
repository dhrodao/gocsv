package gocsv

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	Comment        = '#'
	Separator      = ','
	StringWrapper  = "\""
	NewLine        = "\n"
	CarriageReturn = "\r"
)

func setValue(value reflect.Value, valStr string) error {
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	switch value.Interface().(type) {
	case string:
		val, err := toString(valStr)
		if err != nil {
			return err
		}
		value.SetString(val)
	case bool:
		val, err := toBool(valStr)
		if err != nil {
			return err
		}
		value.SetBool(val)
	case int8, int16, int32, int64:
		val, err := toInt(valStr)
		if err != nil {
			return err
		}
		value.SetInt(val)
	case uint8, uint16, uint32, uint64:
		val, err := toUint(valStr)
		if err != nil {
			return err
		}
		value.SetUint(val)
	case float32, float64:
		val, err := toFloat(valStr)
		if err != nil {
			return err
		}
		value.SetFloat(val)
	default:
		// Check if interface of Unmarshaler and call Unmarshal method
		if interfaceVal, ok := reflect.New(value.Type()).Interface().(Unmarshaler); ok {
			if err := interfaceVal.(Unmarshaler).UnmarshalCSV(valStr); err != nil {
				return err
			}
			value.Set(reflect.ValueOf(interfaceVal).Elem())
		} else {
			// TODO: check if type is subtype of a primitive type
		}
	}

	return nil
}

func toString(val any) (outStr string, err error) {
	inVal := reflect.ValueOf(val)

	switch inVal.Kind() {
	case reflect.String:
		return inVal.String(), nil
	case reflect.Bool:
		b := "false"
		if inVal.Bool() {
			b = "true"
		}
		return b, nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%v", inVal.Int()), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%v", inVal.Uint()), nil
	case reflect.Float32:
		return strconv.FormatFloat(inVal.Float(), byte('f'), -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(inVal.Float(), byte('f'), -1, 64), nil
	default:
		newVal := reflect.New(inVal.Type())
		if interfaceVal, ok := newVal.Interface().(Marshaler); ok {
			newVal.Elem().Set(inVal)
			if outStr, err = interfaceVal.(Marshaler).MarshalCSV(); err != nil {
				return "", err
			}
			return outStr, err
		} else {
			// TODO: check if type if subtype of primitive type
		}
	}

	return "", fmt.Errorf("unknown conversion from %s to string", inVal.Kind())
}

func toBool(valStr string) (bool, error) {
	return strconv.ParseBool(valStr)
}

func toInt(val any) (int64, error) {
	inVal := reflect.ValueOf(val)

	switch inVal.Kind() {
	case reflect.String:
		str := strings.TrimSpace(inVal.String())
		if str == "" {
			return 0, nil
		}
		splitted := strings.SplitN(str, ".", 2)
		return strconv.ParseInt(splitted[0], 0, 64)
	case reflect.Bool:
		if inVal.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return inVal.Int(), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(inVal.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(inVal.Float()), nil
	}

	return 0, fmt.Errorf("unknown conversion from %v to int", inVal.Kind())
}

func toUint(val any) (uint64, error) {
	inVal := reflect.ValueOf(val)

	switch inVal.Kind() {
	case reflect.String:
		str := strings.TrimSpace(inVal.String())
		if str == "" {
			return 0, nil
		}
		splitted := strings.SplitN(str, ".", 2)
		return strconv.ParseUint(splitted[0], 0, 64)
	case reflect.Bool:
		if inVal.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(inVal.Uint()), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return inVal.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return uint64(inVal.Float()), nil
	}

	return 0, fmt.Errorf("unknown conversion from %v to uint", inVal.Kind())
}

func toFloat(val any) (float64, error) {
	inVal := reflect.ValueOf(val)

	switch inVal.Kind() {
	case reflect.String:
		str := strings.TrimSpace(inVal.String())
		if str == "" {
			return 0, nil
		}
		return strconv.ParseFloat(str, 64)
	case reflect.Bool:
		if inVal.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(inVal.Uint()), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(inVal.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return inVal.Float(), nil
	}

	return 0, fmt.Errorf("unknown conversion from %v to float", inVal.Kind())
}
