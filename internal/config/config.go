package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	EngineTypeInMemory = "in-memory"
	LogLevelDebug      = "debug"
)

type SizeBytes int

func (m *SizeBytes) UnmarshalYAML(node *yaml.Node) error {
	size, err := parseBytes(node.Value)
	if err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	*m = SizeBytes(size)
	return nil
}

func (m *SizeBytes) Int() int {
	return int(*m)
}

func parseBytes(s string) (int, error) {
	const bytesInKB = 1024

	kind := []struct {
		kind  string
		ratio int
	}{
		{"KB", bytesInKB}, {"MB", bytesInKB * bytesInKB}, {"B", 1},
	}

	for _, v := range kind {
		if !strings.HasSuffix(s, v.kind) {
			continue
		}

		sInt := strings.Replace(s, v.kind, "", 1)
		size, err := strconv.ParseInt(sInt, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("faile to prase size %q: %w", s, err)
		}
		if size <= 0 {
			return 0, fmt.Errorf("size must be grather than zero")
		}

		return int(size) * v.ratio, nil
	}

	return 0, fmt.Errorf("invalid size value: %q", s)
}

type Config struct {
	Engine struct {
		Type string `yaml:"type"`
	} `yaml:"engine"`

	Network struct {
		Address        string        `yaml:"address"`
		MaxConnections uint          `yaml:"max_connections"`
		MaxMessageSize SizeBytes     `yaml:"max_message_size"`
		IdleTimeout    time.Duration `yaml:"idle_timeout"`
	} `yaml:"network"`

	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`

	Wal struct {
		FlushingBatchSize    int           `yaml:"flushing_batch_size"`
		FlushingBatchTimeout time.Duration `yaml:"flushing_batch_timeout"`
		MaxSegmentSize       SizeBytes     `yaml:"max_segment_size"`
		DataDirectory        string        `yaml:"data_directory"`
	} `yaml:"wal"`
}

func NewConfigWithDefaults() *Config {
	cfg := &Config{}
	cfg.Engine.Type = EngineTypeInMemory
	cfg.Network.Address = "127.0.0.1:3223"
	cfg.Network.MaxConnections = 20
	cfg.Network.IdleTimeout = time.Minute
	cfg.Network.MaxMessageSize = 1024
	cfg.Logging.Output = "./output.log"
	cfg.Logging.Level = LogLevelDebug
	return cfg
}

func FillWithFile(cfg *Config, fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	return nil
}
