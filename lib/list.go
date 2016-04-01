package launchy

import "path/filepath"

// DoesFileHavePlistExtension checks whether a string is the default value. If it is, return
// true else turn false.
func DoesFileHavePlistExtension(path string) bool {
	if filepath.Ext(path) == ".plist" {
		return true
	}
	return false
}

// IsTextProvided checks whether a string is the default value. If it is, return
// true else turn false.
func IsTextProvided(providedText string) bool {
	if providedText != "" {
		return true
	}
	return false
}
