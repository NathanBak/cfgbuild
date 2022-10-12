package config

import (
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

// Init sets some default values in the config.
func (cfg *Config) Init() error {
	// only set defaults once--this prevents users from overwriting set values
	cfg.once.Do(func() {
		cfg.MyInt = 8081
	})
	return nil
}
