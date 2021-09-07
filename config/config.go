package config

import (
	"fmt"
	"os"
	"path"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type Config struct{
	path string
	googleCredentials []byte
	oauthConfig *oauth2.Config
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
	gconf, err := google.ConfigFromJSON(creds, calendar.CalendarScope)
	if err != nil {
		return nil, err
		}
	c.oauthConfig = gconf
	return c, nil
}

func (c *Config) relativPath(p string) string {
	return path.Join(c.path, p)
}

func (c *Config) GoogleCredentials() []byte {
	return c.googleCredentials
}

func (c *Config) Path() string {
	return c.path
}

func (c *Config) OAuthConfig() *oauth2.Config {
	return c.oauthConfig
}
