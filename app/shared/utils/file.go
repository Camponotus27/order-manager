package utils

import (
	"path/filepath"
)

func MergeArrayToPath(array []string) string {
	path := ""
	for _, name := range array {
		path = filepath.Join(name, path)
	}
	return path
}
