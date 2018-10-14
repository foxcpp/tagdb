package main

import "path/filepath"

// Absolute minimal path with all symbolic links resolved.
func canonicalPath(path string) (string, error) {
	target, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	// EvalSymlinks already does filepath.Clean.
	return filepath.Abs(target)
}
