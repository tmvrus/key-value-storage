package inmemory

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tmvrus/key-value-storage/internal/domain"
)

func TestEngine_DeleteSetGet(t *testing.T) {
	t.Parallel()

	storage := New()
	err := storage.Delete(nil, "key")
	require.True(t, errors.Is(err, domain.ErrNotFound))

	err = storage.Set(nil, "key", "value")
	require.NoError(t, err)

	val, err := storage.Get(nil, "key")
	require.NoError(t, err)
	require.Equal(t, "value", val)

	err = storage.Delete(nil, "key")
	require.NoError(t, err)

	val, err = storage.Get(nil, "key")
	require.True(t, errors.Is(err, domain.ErrNotFound))
	require.Empty(t, val)
}
