package config

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		Logging Logging
		CouchDB CouchDB
	}

	Logging struct {
		Debug  bool `envconfig:"LOGS_DEBUG"`
		Trace  bool `envconfig:"LOGS_TRACE"`
		Color  bool `envconfig:"LOGS_COLOR"`
		Pretty bool `envconfig:"LOGS_PRETTY"`
		Text   bool `envconfig:"LOGS_TEXT"`
	}

	CouchDB struct {
		Host                string `envconfig:"COUCHDB_HOST"`
		User                string `envconfig:"COUCHDB_USER"`
		Password            string `envconfig:"COUCHDB_PASSWORD"`
		DBName              string `envconfig:"COUCHDB_DBNAME"`
		DesignViewImportDir string `envconfig:"COUCHDB_DESIGNVIEWIMPORTDIR" default:"design"`
	}
)

// Environ returns the settings from the environment.
func Environ() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	return cfg, err
}

func (c *Config) String() string {
	out, _ := yaml.Marshal(c)
	return string(out)
}
