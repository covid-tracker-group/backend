package tokens

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"simplon.biz/corona/pkg/config"
	"simplon.biz/corona/pkg/tools"
)

const filedirBatchSize = 1000

var errInvalidTokenManagerPath = errors.New("Invalid token manager path")

// DiskTokenManager is a TokenManager which stores all tokens on disk
type DiskTokenManager struct {
	path string
}

func NewDiskTokenManager(path string) (*DiskTokenManager, error) {
	if err := tools.VerifyDirectoryExists(path); err != nil {
		return nil, fmt.Errorf("%w: %v", errInvalidTokenManagerPath, err)
	}

	return &DiskTokenManager{
		path: path,
	}, nil
}

func (dtm *DiskTokenManager) pathForToken(token string) string {
	return filepath.Join(dtm.path, token)
}

func (dtm *DiskTokenManager) CreateToken() (string, error) {
	token := uuid.New().String()
	path := dtm.pathForToken(token)
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	now := time.Now().Unix()
	_, err = fmt.Fprintf(f, "%d", now)
	if err != nil {
		_ = os.Remove(path)
		return "", err
	}
	return token, nil
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

	if st.ModTime().Before(time.Now().Add(-config.ExpireDailyTracingTokensAfter)) {
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

	deleteBefore := time.Now().Add(-config.ExpireDailyTracingTokensAfter)

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
