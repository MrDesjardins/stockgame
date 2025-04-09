package util

import (
	"path/filepath"
	"runtime"
)

// GetProjectRoot returns the absolute path to the project root directory
func GetProjectRoot() string {
	// Get the file path of the current file (db.go)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Cannot get current file path")
	}

	// Navigate up from internal/database to project root
	dir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	return dir
}
