package ls

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuccessfulCalculation(t *testing.T) {
	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)

	testFiles := map[string]string{
		"file1.txt":     "content1",
		"dir/file2.txt": "content2",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	directory = tempDir
	err := lsCmd()
	require.NoError(t, err)

	w.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)
	result := buf.String()

	var output []map[string]string
	err = yaml.Unmarshal([]byte(result), &output)
	require.NoError(t, err)

	// Convert slice of maps to single map for easier testing
	flatOutput := make(map[string]string)
	for _, entry := range output {
		for k, v := range entry {
			flatOutput[k] = v
		}
	}

	// Verify exact hashes
	h1 := sha256.New()
	h1.Write([]byte("content1"))
	expectedHash1 := hex.EncodeToString(h1.Sum(nil))

	h2 := sha256.New()
	h2.Write([]byte("content2"))
	expectedHash2 := hex.EncodeToString(h2.Sum(nil))

	assert.Equal(t, expectedHash1, flatOutput["file1.txt"])
	assert.Equal(t, expectedHash2, flatOutput["dir/file2.txt"])
}

func TestNonExistentDirectory(t *testing.T) {
	directory = "non-existent-dir"
	err := lsCmd()
	assert.Error(t, err)
}

func TestEmptyDirectory(t *testing.T) {
	directory = ""
	err := lsCmd()
	assert.Error(t, err)
	assert.Equal(t, "directory is required", err.Error())
}

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
