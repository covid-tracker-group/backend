package tokens

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"simplon.biz/corona/pkg/tools"
)

const filedirBatchSize = 1000

var errInvalidTokenManagerPath = errors.New("Invalid token manager path")

// DiskTokenManager is a TokenManager which stores all tokens on disk
type DiskTokenManager struct {
	path string
	ttl  time.Duration
}

func NewDiskTokenManager(path string, ttl time.Duration) (*DiskTokenManager, error) {
	if err := tools.VerifyDirectoryExists(path); err != nil {
		return nil, fmt.Errorf("%w: %v", errInvalidTokenManagerPath, err)
	}

	return &DiskTokenManager{
		path: path,
		ttl:  ttl,
	}, nil
}

func (dtm *DiskTokenManager) pathForToken(token string) string {
	return filepath.Join(dtm.path, token)
}

func (dtm *DiskTokenManager) StoreToken(token Token) error {
	path := dtm.pathForToken(token.GetCode())
	f, err := tools.OpenSecureFile(path)
	if err != nil {
		return err
	}
	defer f.Abort()

	encoder := json.NewEncoder(f)
	encoder.SetEscapeHTML(false)
	if err = encoder.Encode(token); err != nil {
		return err
	}

	return f.Close()
}

func (dtm *DiskTokenManager) GetToken(code string, token Token) error {
	valid, err := dtm.VerifyToken(code)
	if !valid {
		return errors.New("Code is not valid")
	}
	if err != nil {
		return err
	}

	path := dtm.pathForToken(code)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(token)
	if err != nil {
		return err
	}
	token.SetCode(code)
	return nil
}

func (dtm *DiskTokenManager) VerifyToken(token string) (bool, error) {
	path := dtm.pathForToken(token)
	st, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if !st.Mode().IsRegular() {
		return false, nil
	}

	if st.ModTime().Before(time.Now().Add(-dtm.ttl)) {
		return false, nil
	}

	return st.Mode().IsRegular(), nil
}

func (dtm *DiskTokenManager) RetractToken(token string) error {
	path := dtm.pathForToken(token)
	err := os.Remove(path)
	if os.IsNotExist((err)) {
		return nil
	}
	return err
}

func (dtm *DiskTokenManager) Expire(onExpire func(string)) error {
	dir, err := os.Open(dtm.path)
	if err != nil {
		return err
	}
	defer dir.Close()

	deleteBefore := time.Now().Add(-dtm.ttl)

	for {
		entries, err := dir.Readdir(filedirBatchSize)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.Mode().IsRegular() && entry.ModTime().Before(deleteBefore) {
				if err = os.Remove(filepath.Join(dtm.path, entry.Name())); err != nil {
					return nil
				}
				onExpire(entry.Name())
			}
		}

		if len(entries) < filedirBatchSize {
			break
		}
	}
	return nil
}
