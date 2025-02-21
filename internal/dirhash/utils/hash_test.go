package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFileHash(t *testing.T) {
	content := []byte("test content")
	tempFile := filepath.Join(t.TempDir(), "file")
	err := os.WriteFile(tempFile, content, 0644)
	require.NoError(t, err)
	defer os.Remove(tempFile)

	hash, err := getFileHash(tempFile)
	require.NoError(t, err)

	h := sha256.New()
	h.Write(content)
	expectedHash := hex.EncodeToString(h.Sum(nil))
	assert.Equal(t, expectedHash, hash)

	_, err = getFileHash("non-existent-file")
	assert.Error(t, err)
}
