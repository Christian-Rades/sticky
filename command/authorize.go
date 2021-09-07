package command

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"c.rades/sticky/config"
	"c.rades/sticky/runner"
	"golang.org/x/oauth2"
)

type Authorize struct {
	config config.Config
}

func (a *Authorize) Run(ctx context.Context, _ *runner.CommandStack) error {
	authConfig := a.config.OAuthConfig()
	authURL := authConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return fmt.Errorf("Unable to read authorization code: %v", err)
	}

	token, err := authConfig.Exchange(ctx, authCode)
	if err != nil {
		return fmt.Errorf("Unable to retrieve token from web: %w", err)
	}

	path := path.Join(a.config.Path(), "token.json")
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Unable to cache oauth token: %w", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}
