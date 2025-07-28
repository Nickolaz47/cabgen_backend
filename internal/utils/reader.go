package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadJSONFile[T any](filepath string) ([]T, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %s: %v", filepath, err)
	}

	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s: %v", filepath, err)
	}

	return items, nil
}
