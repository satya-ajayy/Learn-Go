package config

import (
	"learn-go/errors"
)

var DefaultConfig = []byte(`
logger:
  level: "info"

listen: ":8888"

prefix: "/ajay-verse"

is_prod_mode: false

mongo:
  uri: "mongodb://localhost:27017"

redis:
  uri: "localhost:6379"
`)

type Config struct {
	Logger     Logger `koanf:"logger"`
	Listen     string `koanf:"listen"`
	Prefix     string `koanf:"prefix"`
	IsProdMode bool   `koanf:"is_prod_mode"`
	Mongo      Mongo  `koanf:"mongo"`
	Redis      Redis  `koanf:"redis"`
}

type Logger struct {
	Level string `koanf:"level"`
}

type Mongo struct {
	URI string `koanf:"uri"`
}

type Redis struct {
	URI string `koanf:"uri"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
	ve := errors.ValidationErrs()

	if c.Listen == "" {
		ve.Add("listen", "cannot be empty")
	}
	if c.Logger.Level == "" {
		ve.Add("logger.level", "cannot be empty")
	}
	if c.Mongo.URI == "" {
		ve.Add("mongo.uri", "cannot be empty")
	}
	if c.Redis.URI == "" {
		ve.Add("redis.uri", "cannot be empty")
	}

	return ve.Err()
}
