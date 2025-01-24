package os

import "os"

func WriteFile(path, data string) {
	// create or override file
	destFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer destFile.Close()
	// write data to dest file
	if _, err = destFile.WriteString(data); err != nil {
		panic(err)
	}
}
