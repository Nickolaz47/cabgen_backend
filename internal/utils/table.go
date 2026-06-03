package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"reflect"
	"time"
)

func GenerateDynamicTSV(items any) ([]byte, error) {
	sliceVal := reflect.ValueOf(items)
	if sliceVal.Kind() != reflect.Slice {
		return nil, fmt.Errorf(
			"Invalid items format: expected a slice, got %T", items)
	}

	if sliceVal.Len() == 0 {
		return []byte{}, nil
	}

	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)
	writer.Comma = '\t'

	firstItem := sliceVal.Index(0)
	// If the slice contains pointers dereference it to get the actual struct
	if firstItem.Kind() == reflect.Pointer {
		firstItem = firstItem.Elem()
	}
	typ := firstItem.Type()

	var headers []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		headers = append(headers, tag)
	}

	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// Iterate over each item in the slice to generate data rows
	for i := 0; i < sliceVal.Len(); i++ {
		item := sliceVal.Index(i)

		// Dereference pointer if necessary to access struct fields safely
		if item.Kind() == reflect.Pointer {
			item = item.Elem()
		}

		var row []string
		// Iterate over the fields of the current struct
		for j := 0; j < item.NumField(); j++ {
			field := typ.Field(j)
			tag := field.Tag.Get("json")

			// Must apply the same filtering rule used for headers to maintain
			// column alignment
			if tag == "" || tag == "-" {
				continue
			}

			// Extract the actual field value and format it into a string
			val := item.Field(j)
			row = append(row, formatValue(val))
		}

		// Write the formatted row to the buffer
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	// Flush any buffered data to the underlying io.Writer
	writer.Flush()

	return buffer.Bytes(), writer.Error()
}

func formatValue(v reflect.Value) string {
	// Handle pointer values safely to prevent panics
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return "" // Return an empty string for nil pointers
		}
		// Dereference the pointer to get the actual underlying value
		v = v.Elem()
	}

	switch val := v.Interface().(type) {
	case time.Time:
		// Format dates to a readable standard instead of the default Go time 
		// format
		return val.Format("02-01-2006 15:04:05")
	case []byte:
		// Specifically target datatypes.JSON or raw JSON bytes
		// If the byte array is empty or represents a JSON null, return an 
		// empty JSON object
		if len(val) == 0 || string(val) == "null" {
			return ""
		}
		return string(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
