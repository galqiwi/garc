package lifetime_verifier

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func touchFile(path string, modTime time.Time) error {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	if err := os.Chtimes(path, modTime, modTime); err != nil {
		return err
	}

	return nil
}

func checkLifetimes(t *testing.T, maxLifetimeS int64, lifetimes []time.Duration) bool {
	dirname, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer func() {
		err := os.RemoveAll(dirname)
		require.NoError(t, err)
	}()

	for fileIdx, lifetime := range lifetimes {
		err = touchFile(
			filepath.Join(dirname, fmt.Sprintf("%v.txt", fileIdx)),
			time.Now().Add(-lifetime),
		)
		require.NoError(t, err)
	}

	verifier := NewLifetimeVerifier(&LifetimeVerifierConfig{
		ErrorLifetimeS: maxLifetimeS,
		Dirs:           []string{dirname},
	})

	assert.Equal(t, "lifetime_verifier", verifier.Name())

	return len(verifier.Verify()) > 0
}

func TestName(t *testing.T) {
	var maxLifetimeS int64 = 3600
	require.False(t, checkLifetimes(t, maxLifetimeS, nil))
	require.False(t, checkLifetimes(t, maxLifetimeS, []time.Duration{time.Minute}))
	require.True(t, checkLifetimes(t, maxLifetimeS, []time.Duration{time.Hour * 2}))
	require.True(t, checkLifetimes(t, maxLifetimeS, []time.Duration{time.Minute, time.Hour * 2}))
}
