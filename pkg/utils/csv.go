package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"reflect"
)

// GenerateCSV converts a slice of structs to CSV bytes using reflection.
// It uses "json" tags for headers.
func GenerateCSV(data interface{}) ([]byte, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("GenerateCSV: expected a slice, got %s", v.Kind())
	}

	if v.Len() == 0 {
		return []byte(""), nil // empty CSV
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	firstElem := v.Index(0)
	if firstElem.Kind() == reflect.Ptr {
		firstElem = firstElem.Elem()
	}

	if firstElem.Kind() != reflect.Struct {
		return nil, fmt.Errorf("GenerateCSV: expected slice of structs, got slice of %s", firstElem.Kind())
	}

	t := firstElem.Type()
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		header := field.Tag.Get("json")
		if header == "" {
			header = field.Name
		}
		headers = append(headers, header)
	}

	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		var row []string
		for j := 0; j < elem.NumField(); j++ {
			fieldValue := elem.Field(j).Interface()
			row = append(row, fmt.Sprintf("%v", fieldValue))
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
