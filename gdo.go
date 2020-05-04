// Package gdo implements the garage door controller
package gdo

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config is the gdo config struct.  Its members capture config for sections of the program
type Config struct {
	Oauth2  Oauth2Cfg    `yaml:"oauth2"`
	Server  WebserverCfg `yaml:"webserver"`
	Door    DoorCfg      `yaml:"door"`
	Testing bool
}

// Oauth2Cfg catptures Oauth configs
type Oauth2Cfg struct {
	ClientSecret string        `yaml:"clientSecret"`
	ClientID     string        `yaml:"clientID"`
	AppName      string        `yaml:"appName"`
	Redirect     string        `yaml:"redirect"`
	Scopes       []string      `yaml:"scopes"`
	Access       []AccessEntry `yaml:"access"`
}

// Timeout captures server timeout settings
type Timeout struct {
	Server int `yaml:"server"`
	Write  int `yaml:"write"`
	Read   int `yaml:"read"`
	Idle   int `yaml:"idle"`
}

// WebserverCfg captures webserver configs
type WebserverCfg struct {
	Host     string  `yaml:"host"`
	Port     int     `yaml:"port"`
	TLSPort  int     `yaml:"tlsPort"`
	Addr     string  `yaml:"addr"`
	Addr6    string  `yaml:"addr6"`
	Cachedir string  `yaml:"cachedir"`
	Timeout  Timeout `yaml:"timeout"`
}

//AccessEntry populates the access array and allow google id holders to run the services
type AccessEntry struct {
	Username string `yaml:"email,omitempty"`
	GoogleID string `yaml:"googleId,omitempty"`
}

// DoorCfg configures the door parameters
type DoorCfg struct {
	ClosedPin   int `yaml:"closedPin"`
	ClosedValue int `yaml:"closedValue"`
	ClickRelay  int `yaml:"clickRelay"`
}

var cfg *Config
var cfgFileName = "/etc/gdo.conf"

// NewConfig gets the program config from the specified file
func NewConfig(configFileName string) (*Config, error) {
	cfgFileName = configFileName
	cfg = nil
	return GetConfig()
}

// GetConfig gets the program config
func GetConfig() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}
	f, err := os.Open(cfgFileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(f)
	cfg = &Config{}
	// Start YAML decoding from file
	if err := d.Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
