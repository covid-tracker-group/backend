package tools

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var errCanNotWriteToAbortedFile = errors.New("Can not write to aborted secure file")

// SecureFile is a simple file that writes data to a temporary file, which is atomically renamed,
// possibly replacing an existing file, on close. It implements the WriteCloser interface.
type SecureFile struct {
	file   *os.File
	path   string
	closed bool
}

// OpenSecureFile opens a new SecureFile.
func OpenSecureFile(path string) (*SecureFile, error) {
	f, err := ioutil.TempFile(filepath.Dir(path), "update*")
	if err != nil {
		return nil, err
	}
	return &SecureFile{
		file:   f,
		path:   path,
		closed: false,
	}, nil
}

func (sf *SecureFile) Write(p []byte) (n int, err error) {
	if sf.closed {
		return 0, errCanNotWriteToAbortedFile
	}
	return sf.file.Write(p)
}

func (sf *SecureFile) Abort() {
	if sf.closed {
		return
	}
	_ = sf.file.Close()
	_ = os.Remove(sf.file.Name())
	sf.closed = true
}

func (sf *SecureFile) Close() error {
	if sf.closed {
		return nil
	}

	if err := sf.file.Close(); err != nil {
		_ = os.Remove(sf.file.Name())
		return err
	}

	if err := os.Rename(sf.file.Name(), sf.path); err != nil {
		return fmt.Errorf("Error moving temporary file into place: %w", err)
	}

	sf.closed = true
	return nil
}
