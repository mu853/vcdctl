package module

import (
	"encoding/base64"
	"fmt"
	"log"
)

type Config struct {
	CurrentSite string `json:"current-site" mapstructure:"current-site"`
	Sites       []Site `json:"sites"`
}

type Site struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (s *Site) GetCredential() string {
	passwordText, err := base64.StdEncoding.DecodeString(s.Password)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString([]byte(s.User + ":" + string(passwordText)))
}

func (c *Config) GetCurrentSite() (Site, error) {
	for _, s := range c.Sites {
		if s.Name == c.CurrentSite {
			return s, nil
		}
	}
	return Site{}, fmt.Errorf("site '%s' not found", c.CurrentSite)
}

func (c *Config) GetSite(name string) (Site, error) {
	for _, s := range c.Sites {
		if s.Name == name {
			return s, nil
		}
	}
	return Site{}, fmt.Errorf("site '%s' not found", name)
}

func (t *Site) SetPassword(password string) {
	t.Password = base64.StdEncoding.EncodeToString([]byte(password))
}
