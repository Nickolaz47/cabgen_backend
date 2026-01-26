package utils

import (
	"os"
	"path/filepath"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
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
	if config.AppRoot != "" {
		return config.AppRoot, nil
	}

	if cwd, err := os.Getwd(); err == nil {
		current := cwd
		for {
			if _, err := os.Stat(
				filepath.Join(current, "go.mod")); err == nil {
				return current, nil
			}

			parent := filepath.Dir(current)
			if parent == current {
				break
			}
			current = parent
		}
	}

	return os.Getwd()
}
