package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONMap map[string]string

// Saves data in database (map -> JSON -> []byte).
func (j *JSONMap) Scan(value any) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}

	var bytes []byte

	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("invalid type for JSONMap: %T", value)
	}

	return json.Unmarshal(bytes, j)
}

// Reads data from database ([]byte -> JSON -> map).
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}

	return json.Marshal(j)
}
