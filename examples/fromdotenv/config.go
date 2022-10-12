package main

import "time"

// Config is an example struct usable by the cfgbuild package.
type Config struct {
	MyInt      int           `json:"myInt" envvar:"MY_INT"`
	MyFloat    float32       `json:"myFloat" envvar:"MY_FLOAT"`
	MyDuration time.Duration `json:"myDuration" envvar:"MY_DURATION"`
	MyTime     time.Time     `json:"myTime" envvar:"MY_TIME"`
	MyBytes    []byte        `json:"myBytes" envvar:"MY_BYTES"`
	MyString   string        `json:"myString" envvar:"MY_STRING"`
	MyBool     bool          `json:"myBool" envvar:"MY_BOOL"`
}

// Init sets some default values in the config.
func (cfg *Config) Init() error {
	cfg.MyInt = 8081
	cfg.MyTime = time.Date(2000, time.March, 17, 0, 13, 37, 0, time.UTC)
	return nil
}