package tarball

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateTarball(t *testing.T) {
	defer goleak.VerifyNone(t)

	dirname, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer func() {
		err := os.RemoveAll(dirname)
		require.NoError(t, err)
	}()

	srcPath := filepath.Join(dirname, "src")
	dstPath := filepath.Join(dirname, "dst")

	err = os.Mkdir(srcPath, 0700)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(srcPath, "file"), []byte("test"), 0777)
	require.NoError(t, err)

	buf := &bytes.Buffer{}

	err = CreateTarball(srcPath, "", buf, []string{})
	require.NoError(t, err)

	err = ExtractTarball(buf, dstPath)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(dstPath, "file"))
	require.NoError(t, err)
	require.Equal(t, []byte("test"), content)
}
