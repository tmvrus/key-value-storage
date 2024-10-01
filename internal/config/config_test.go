package config

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func Test_ParseSize(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		in   string
		out  uint
		err  bool
	}{
		{
			name: "valid 4KB",
			in:   "4KB",
			out:  4 * 1024,
		},
		{
			name: "valid 10MB",
			in:   "10MB",
			out:  10 * 1024 * 1024,
		},
		{
			name: "valid 100B",
			in:   "10B",
			out:  10,
		},
		{
			name: "invalid 0",
			in:   "0M",
			err:  true,
		},
		{
			name: "invalid size class",
			in:   "100XXX",
			err:  true,
		},
		{
			name: "invalid number class",
			in:   "1a0b0MB",
			err:  true,
		},
	}

	for i := range tt {
		size, err := parseBytes(tt[i].in)
		if tt[i].err {
			require.Errorf(t, err, tt[i].name)
		} else {
			require.Equalf(t, tt[i].out, size, tt[i].name, err)
		}
	}
}

func Test_Config(t *testing.T) {
	t.Parallel()

	t.Run("check defaults", func(t *testing.T) {
		t.Parallel()

		defCfg := *NewConfigWithDefaults()
		newCfg := Config{}

		require.NotEqual(t, defCfg, newCfg)
	})

	t.Run("check yaml unmarshalling", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}
		err := yaml.Unmarshal(yamlData, &cfg)

		require.NoError(t, err)
		require.Equal(t, "in_memory", cfg.Engine.Type)
		require.Equal(t, "info", cfg.Logging.Level)
		require.Equal(t, uint(4*1024), cfg.Network.MaxMessageSize.Uint())
	})
}

var yamlData = []byte(`
engine:
  type: "in_memory"
network:
  address: "127.0.0.1:3223"
  max_connections: 100
  max_message_size: "4KB"
  idle_timeout: 5m
logging:
  level: "info"
  output: "/log/output.log"
`)
