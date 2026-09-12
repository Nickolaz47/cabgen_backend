package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func IsGzip(r io.Reader) (io.Reader, bool) {
	magic := make([]byte, 2)
	if _, err := io.ReadFull(r, magic); err != nil {
		return io.MultiReader(bytes.NewReader(magic), r), false
	}
	isGz := magic[0] == 0x1f && magic[1] == 0x8b
	return io.MultiReader(bytes.NewReader(magic), r), isGz
}
