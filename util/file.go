package util

import (
	"os"
	"regexp"
)

func ReplacePatternInFile(filePath string, pattern string, replacement string) error {
	re := regexp.MustCompile(pattern)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	modifiedContent := re.ReplaceAllString(string(content), replacement)

	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	permissions := info.Mode()

	err = os.WriteFile(filePath, []byte(modifiedContent), permissions)
	if err != nil {
		return err
	}

	return nil
}
