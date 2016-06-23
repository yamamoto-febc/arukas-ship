package ship

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Arukas *ArukasConfig
	Serve  *ServeConfig
}

type ArukasConfig struct {
	Token  string
	Secret string
}

type ServeConfig struct {
	Token string
	Port  int64
}

func InitializeConfig() (*Config, error) {

	config := &Config{
		Arukas: &ArukasConfig{
			Token:  os.Getenv("ARUKAS_JSON_API_TOKEN"),
			Secret: os.Getenv("ARUKAS_JSON_API_SECRET"),
		},
		Serve: &ServeConfig{
			Token: os.Getenv("SHIP_TOKEN"),
			Port:  -1,
		},
	}

	if port, ok := os.LookupEnv("SHIP_PORT"); ok {
		res, err := strconv.ParseInt(port, 10, 64)
		if err != nil {
			return nil, err
		}
		config.Serve.Port = res
	}

	err := config.validate()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.Arukas.Token == "" {
		return fmt.Errorf("Missing %s", "ARUKAS_JSON_API_TOKEN")
	}
	if c.Arukas.Secret == "" {
		return fmt.Errorf("Missing %s", "ARUKAS_JSON_API_SECRET")
	}
	if c.Serve.Token == "" {
		return fmt.Errorf("Missing %s", "SHIP_TOKEN")
	}
	if c.Serve.Port == -1 {
		return fmt.Errorf("Missing %s", "SHIP_PORT")
	}

	if !(0 < c.Serve.Port && c.Serve.Port <= 65535) {
		return fmt.Errorf("%s must be between 1 and 65535", "SHIP_PORT")

	}

	return nil
}
