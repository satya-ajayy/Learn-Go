package config

import (
	// Local Packages
	"learn-go/errors"
)

var DefaultConfig = []byte(`
application: "learn-go"

logger:
  level: "debug"

listen: ":8888"

prefix: "/learn-go"

is_prod_mode: false

mongo:
  uri: "mongodb://localhost:27017"

redis:
  uri: "localhost:6379"
  password: ""
`)

type Config struct {
	Application string `koanf:"application"`
	Logger      Logger `koanf:"logger"`
	Listen      string `koanf:"listen"`
	Prefix      string `koanf:"prefix"`
	IsProdMode  bool   `koanf:"is_prod_mode"`
	Mongo       Mongo  `koanf:"mongo"`
	Redis       Redis  `koanf:"redis"`
}

type Logger struct {
	Level string `koanf:"level"`
}

type Mongo struct {
	URI string `koanf:"uri"`
}

type Redis struct {
	URI      string `koanf:"uri"`
	Password string `koanf:"password"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
	ve := errors.ValidationErrs()

	if c.Application == "" {
		ve.Add("application", "cannot be empty")
	}
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
