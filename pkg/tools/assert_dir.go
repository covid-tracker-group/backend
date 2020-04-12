package tools

import (
	"fmt"
	"os"
)

func VerifyDirectoryExists(path string) error {
	st, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !st.Mode().IsDir() {
		return fmt.Errorf("%s it not a directory", path)
	}

	return nil
}
