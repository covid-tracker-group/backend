package keystorage

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"simplon.biz/corona/pkg/tools"
)

const filedirBatchSize = 1000

var errInvalidKeyStoragePath = errors.New("Invalid key storage path")
var errFailedToCreateKeyFile = errors.New("Failed to create a key file")
var errFailedToUpdateKeyFile = errors.New("Failed to update a key file")
var errListingKeys = errors.New("Error listing keys")

// DiskKeyStorage is a KeyStorage which stores all data on disk
type DiskKeyStorage struct {
	path string
}

func NewDiskKeyStorage(path string) (*DiskKeyStorage, error) {
	if err := tools.VerifyDirectoryExists(path); err != nil {
		return nil, fmt.Errorf("%w: %v", errInvalidKeyStoragePath, err)
	}

	return &DiskKeyStorage{
		path: path,
	}, nil
}

func (dks *DiskKeyStorage) pathForAuthorisationCode(authorisationCode string) string {
	return filepath.Join(dks.path, authorisationCode+".csv")
}

func (dks *DiskKeyStorage) AddKeyRecords(authorisationCode string, records []KeyRecord) error {
	// TODO: for extra safety we could copy the old file to a temporary file, append data to it and
	// then move the new file into place.
	path := dks.pathForAuthorisationCode(authorisationCode)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("%w: %v", errFailedToUpdateKeyFile, err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	for _, record := range records {
		err = writer.Write([]string{
			strconv.FormatInt(record.ProcessedAt.Unix(), 10),
			strconv.Itoa(record.DayNumber),
			record.DailyTracingKey,
		})
		if err != nil {
			return fmt.Errorf("%w: %v", errFailedToUpdateKeyFile, err)
		}
	}
	writer.Flush()
	if err = writer.Error(); err != nil {
		return fmt.Errorf("%w: %v", errFailedToUpdateKeyFile, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("%w: %v", errFailedToUpdateKeyFile, err)
	}

	return nil
}

func (dks *DiskKeyStorage) PurgeRecords(authorisationCode string) error {
	path := dks.pathForAuthorisationCode(authorisationCode)
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (dks *DiskKeyStorage) ListRecords(recordChan chan RawKeyRecord, errChan chan interface{}) {
	dir, err := os.Open(dks.path)
	if err != nil {
		errChan <- fmt.Errorf("%w: can not open directory: %v", errListingKeys, err)
		close(recordChan)
		return
	}
	defer dir.Close()

	var path string

	for {
		entries, err := dir.Readdir(filedirBatchSize)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			errChan <- fmt.Errorf("%w: can not read directory: %v", errListingKeys, err)
			return
		}

		for _, entry := range entries {
			if entry.Mode().IsRegular() && filepath.Ext(entry.Name()) == ".csv" {
				path = filepath.Join(dks.path, entry.Name())
				err = dks.streamFromFile(path, recordChan)
				if err != nil {
					errChan <- fmt.Errorf("%w: error streaming from key file %s: %v", errListingKeys, path, err)
					return
				}
			}
		}

		if len(entries) < filedirBatchSize {
			break
		}
	}

	close(recordChan)
}

func (dks *DiskKeyStorage) streamFromFile(path string, records chan RawKeyRecord) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for _, row := range rows {
		records <- RawKeyRecord{
			ProcessedAt:     row[0],
			DayNumber:       row[1],
			DailyTracingKey: row[2],
		}
	}
	return nil
}
