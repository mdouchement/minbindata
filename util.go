package main

import (
	"path/filepath"
)

func getOutputFilename(base, input, output string) (string, error) {
	rel, err := filepath.Rel(base, input)
	if err != nil {
		return "", err
	}
	return filepath.Clean(filepath.ToSlash(filepath.Join(output, filepath.Base(base), rel))), nil
}
