package tools

import "encoding/csv"

// SecureCSVWriter is utility class to make it easy to securely write CSV files.
type SecureCSVWriter struct {
	file      *SecureFile
	csvWriter *csv.Writer
}

// NewSecureCSVWriter creates a new SecureCSVWriter to secure write CSV files.
func NewSecureCSVWriter(path string) (*SecureCSVWriter, error) {
	file, err := OpenSecureFile(path)
	if err != nil {
		return nil, err
	}

	csvWriter := csv.NewWriter(file)
	writer := &SecureCSVWriter{
		file:      file,
		csvWriter: csvWriter,
	}
	return writer, nil
}

func (wr *SecureCSVWriter) Write(record []string) error {
	err := wr.csvWriter.Write(record)
	if err != nil {
		wr.file.Abort()
	}
	return err
}

func (wr *SecureCSVWriter) Abort() {
	wr.file.Abort()
}

func (wr *SecureCSVWriter) Close() error {
	wr.csvWriter.Flush()
	if err := wr.csvWriter.Error(); err != nil {
		wr.file.Abort()
		return err
	}
	return wr.file.Close()
}
