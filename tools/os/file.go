package os

import (
	"os"
	"path/filepath"
)

func WriteFile(workDir, shortname, data string) string {
	destPath := filepath.Join(workDir, shortname)
	// create or override file
	destFile, err := os.Create(destPath)
	if err != nil {
		panic(err)
	}
	defer destFile.Close()
	// write data to dest file
	if _, err = destFile.WriteString(data); err != nil {
		return ""
	}
	return destPath
}
