package tool

import "os"

// PathCreate creates a path if it does not exist.
func PathCreate(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return false
		}
	}
	return true
}
