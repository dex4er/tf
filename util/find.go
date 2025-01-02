package util

import (
	"os"
	"path/filepath"
)

const (
	OpentofuVersionFile  = ".opentofu-version"
	TerraformVersionFile = ".terraform-version"
)

func checkVersionFilesInDir(dir string) string {
	if _, err := os.Stat(filepath.Join(dir, OpentofuVersionFile)); err == nil {
		return OpentofuVersionFile
	}
	if _, err := os.Stat(filepath.Join(dir, TerraformVersionFile)); err == nil {
		return TerraformVersionFile
	}
	return ""
}

func FindDotVersionFile() string {
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		if versionFile := checkVersionFilesInDir(currentDir); versionFile != "" {
			return versionFile
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return ""
}
