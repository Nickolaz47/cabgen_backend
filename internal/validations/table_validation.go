package validations

import (
	"errors"
	"io"
	"strings"
)

const MaxTableSize = int64(50) << 20

var ErrTableTooLarge = errors.New("table exceeds size limit")

func IsAllowedTableFile(fileName string) bool {
	if len(fileName) > 255 {
		return false
	}
	return strings.HasSuffix(strings.ToLower(fileName), ".xlsx")
}

func ReadTablePart(part io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(part, MaxTableSize+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > MaxTableSize {
		return nil, ErrTableTooLarge
	}
	return data, nil
}
