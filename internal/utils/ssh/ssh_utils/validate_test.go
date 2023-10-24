package ssh_utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidatePath(t *testing.T) {
	require.Error(t, validatePath(";"))
	require.Error(t, validatePath("a;"))
	require.Error(t, validatePath("a;a"))
	require.Error(t, validatePath(";a"))
	require.Error(t, validatePath("\\a\\dasd"))
	require.Error(t, validatePath("localfile"))
	require.NoError(t, validatePath("/file"))
	require.NoError(t, validatePath("/dir/"))
	require.NoError(t, validatePath("/dir/file"))
	require.NoError(t, validatePath("/dir/dir/"))
	require.NoError(t, validatePath("/dir/.file"))
}
