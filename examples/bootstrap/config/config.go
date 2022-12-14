package config

import (
	"fmt"
	"sync"

	"github.com/NathanBak/cfgbuild"
	"github.com/joho/godotenv"
)

// Config is an example struct usable by the cfgbuild package.
type Config struct {
	MyInt    int     `json:"myInt" envvar:"MY_INT"`
	MyFloat  float32 `json:"myFloat" envvar:"MY_FLOAT"`
	MyString string  `json:"myString" envvar:"MY_STRING"`
	MyBool   bool    `json:"myBool" envvar:"MY_BOOL"`
	once     sync.Once
}

// New will load the .env file env vars (if not already set) and then create and return a new Config
// based on the env vars.
func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, nil
	}

	// create a new config Builder
	builder := cfgbuild.Builder[*Config]{}
	// build the new config setting the values from the env vars
	return builder.Build()
}

// CfgBuildInit sets some default values in the config.  It is called by cfgbuild.Builder.Build().
func (cfg *Config) CfgBuildInit() error {
	// only set defaults once--this prevents users from overwriting set values
	cfg.once.Do(func() {
		cfg.MyInt = 8081
	})
	return nil
}

// CfgBuildValidate can check that the certain set values are valid.  It is called by
// cfgbuild.Builder.Build().
func (cfg *Config) CfgBuildValidate() error {
	if cfg.MyInt < 8080 || cfg.MyInt > 9999 {
		return fmt.Errorf("MY_INT value %d is out of range", cfg.MyInt)
	}
	return nil
}
