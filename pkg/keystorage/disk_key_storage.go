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

// FormatTagVersion1 is the version tag used to identify the current file version.
const FormatTagVersion1 = "v1"
const filedirBatchSize = 1000

var errInvalidKeyStoragePath = errors.New("Invalid key storage path")
var errListingKeys = errors.New("Error listing keys")
var errUnsupportedKeyVersion = errors.New("Unsupported key file version")

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
	path := dks.pathForAuthorisationCode(authorisationCode)

	newFile, err := tools.NewSecureCSVWriter(path)
	if err != nil {
		return err
	}
	defer newFile.Abort()

	if err = newFile.Write([]string{FormatTagVersion1, "", ""}); err != nil {
		return err
	}

	existingRecords, err := dks.readFromFile((path))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	recordMap := make(map[string]bool)
	for _, record := range existingRecords {
		recordMap[record.DailyTracingKey] = true
	}

	for _, record := range records {
		_, exists := recordMap[record.DailyTracingKey]
		if !exists {
			err = newFile.Write([]string{
				strconv.FormatInt(record.ProcessedAt.Unix(), 10),
				strconv.Itoa(record.DayNumber),
				record.DailyTracingKey,
			})
			if err != nil {
				return err
			}
		}
	}

	return newFile.Close()
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
				records, err := dks.readFromFile(path)
				if err != nil {
					errChan <- fmt.Errorf("%w: error reading from key file %s: %v", errListingKeys, path, err)
					return
				}
				for _, record := range records {
					recordChan <- record
				}
			}
		}

		if len(entries) < filedirBatchSize {
			break
		}
	}

	close(recordChan)
}

func (dks *DiskKeyStorage) readFromFile(path string) ([]RawKeyRecord, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if rows[0][0] != FormatTagVersion1 {
		return nil, fmt.Errorf("%w: %v", errUnsupportedKeyVersion, rows[0][0])
	}
	rows = rows[1:]

	records := make([]RawKeyRecord, len(rows))
	for ix, row := range rows {
		records[ix] = RawKeyRecord{
			ProcessedAt:     row[0],
			DayNumber:       row[1],
			DailyTracingKey: row[2],
		}
	}
	return records, nil
}
