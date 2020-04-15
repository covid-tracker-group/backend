package tokens

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func tempDir(t *testing.T) (path string, cleanup func()) {
	t.Helper()

	tmpDir, err := ioutil.TempDir("", "*")
	if err != nil {
		t.Fatalf("could not create temporary directory: %s", err)
	}

	return tmpDir, func() {
		err = os.RemoveAll(tmpDir)
		if err != nil {
			t.Errorf("could not remove temporary directory: %s", err)
		}
	}
}

func TestNewDiskTokenManager(t *testing.T) {
	if _, err := NewDiskTokenManager("missing", time.Hour); err == nil {
		t.Error("No error returned for missing directory ")
	}

	tempDir, cleanup := tempDir(t)
	defer cleanup()
	if _, err := NewDiskTokenManager(tempDir, time.Hour); err != nil {
		t.Fatalf("Error DiskTokenManager: %v", err)
	}
}

func TestRoundtripToken(t *testing.T) {
	tempDir, cleanup := tempDir(t)
	defer cleanup()

	dtm, err := NewDiskTokenManager(tempDir, time.Hour)
	if err != nil {
		t.Fatalf("Error DiskTokenManager: %v", err)
	}

	token := NewTracingAuthenticationToken()
	if err = dtm.StoreToken(token); err != nil {
		t.Fatalf("Error storing token: %v", err)
	}
	var readToken TracingAuthenticationToken
	if err = dtm.GetToken(token.GetCode(), &readToken); err != nil {
		t.Fatalf("Error reading token: %v", err)
	}

	if token.GetCode() != readToken.GetCode() {
		t.Errorf("Code is changed during roundtrip: %s -> %s", token.GetCode(), readToken.GetCode())
	}

	if err = dtm.RetractToken(token.GetCode()); err != nil {
		t.Fatalf("Error tracking token: %v", err)
	}

	if err = dtm.GetToken(token.GetCode(), &readToken); err == nil {
		t.Fatal("GetToken succeeds after RetractToken")
	}

	if !os.IsNotExist(err) {
		t.Fatalf("GetToken does not return ENOENT after RetractToken, but %v", err)
	}
}
