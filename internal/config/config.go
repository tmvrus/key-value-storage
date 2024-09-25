package config

import "time"

const (
	EngineTypeInMemory = "in-memory"
	LogLevelDebug      = "debug"
)

type Config struct {
	Engine struct {
		Type string `yaml:"type"`
	} `yaml:"engine"`

	Network struct {
		Address        string        `yaml:"address"`
		MaxConnections uint          `yaml:"max_connections"`
		MaxMessageSize uint          `yaml:"max_message_size"`
		IdleTimeout    time.Duration `yaml:"idle_timeout"`
	} `yaml:"network"`

	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
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
