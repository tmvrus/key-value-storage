package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandType_Valid(t *testing.T) {
	t.Parallel()

	valid := []CommandType{CommandDelete, CommandSet, CommandGet}
	for _, v := range valid {
		require.True(t, v.Valid())
	}

	notValid := []CommandType{"", "zzz"}
	for _, v := range notValid {
		require.False(t, v.Valid())
	}
}
