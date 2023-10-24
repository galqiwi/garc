package shell

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCommandExists(t *testing.T) {
	exists, err := CommandExists("which")
	require.NoError(t, err)
	require.True(t, exists)

	commandUuid, err := uuid.NewUUID()
	require.NoError(t, err)

	exists, err = CommandExists(commandUuid.String())
	require.NoError(t, err)
	require.False(t, exists)
}
