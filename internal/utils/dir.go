package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

/*
GetProjectRoot searches for the root directory of the Go project by traversing upwards
from the executable's location until it finds a directory containing a 'go.mod' file.

Returns:
  - string: The absolute path to the project root directory if found.
  - error:  An error if the project root cannot be located (typically meaning
    the executable is not running from within a Go module project).
*/
func GetProjectRoot() (string, error) {
	goMod := "go.mod"

	// Get the current file path
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get caller information")
	}

	rootDir := filepath.Dir(filename)
	for rootDir != "/" {
		files, err := os.ReadDir(rootDir)
		if err != nil {
			return "", err
		}

		for _, file := range files {
			if file.Name() == goMod {
				return rootDir, nil
			}
		}

		// Remove the last path
		rootDir = filepath.Dir(rootDir)
	}

	return rootDir, nil
}
