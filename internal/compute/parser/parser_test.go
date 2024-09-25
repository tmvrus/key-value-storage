package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tmvrus/key-value-storage/internal/domain"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tt := []struct {
		in  string
		out domain.Command
		err bool
	}{
		{
			in: "GET key",
			out: domain.Command{
				Type: domain.CommandGet,
				Key:  "key",
			},
		},
		{
			in:  "GET ",
			err: true,
		},
		{
			in:  "GET key value",
			err: true,
		},
		{
			in: "DELETE key",
			out: domain.Command{
				Type: domain.CommandDelete,
				Key:  "key",
			},
		},
		{
			in:  "Delete key",
			err: true,
		},
		{
			in:  "set key value",
			err: true,
		},
		{
			in: "SET key value",
			out: domain.Command{
				Type:  domain.CommandSet,
				Key:   "key",
				Value: "value",
			},
		},
		{
			in:  "SET key ",
			err: true,
		},
	}

	for i, c := range tt {
		cmd, err := Parse(c.in)
		if c.err {
			require.Errorf(t, err, "iter %d", i)
		} else {
			require.Equal(t, c.out, cmd)
		}
	}
}
