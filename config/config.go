package config

import (
	"fmt"
	"os"
	"path"
)

type Config struct{
	path string
	googleCredentials []byte
}

type CredentialsNotFound struct {
	Err error
}

func (e CredentialsNotFound) Error() string {
	return fmt.Sprintf("credentials not found: %s", e.Err)
}

func OpenConfig(configPath string) (*Config, error) {
	_, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(configPath, 0777)
		} else {
			return nil, err
		}
	}
	c := &Config{
		path: configPath,
	}
	creds, err := os.ReadFile(c.relativPath("credentials.json"))
	if err != nil {
		return nil, CredentialsNotFound{err}
	}
	c.googleCredentials = creds
	return c, nil
}

func (c *Config) relativPath(p string) string {
	return path.Join(c.path, p)
}

func (c *Config) GoogleCredentials() []byte {
	return c.googleCredentials
}
