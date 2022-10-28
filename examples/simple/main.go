package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NathanBak/cfgbuild"
)

type config struct {
	cfgbuild.BaseConfig         // adds default implementations of required functions
	MyInt               int     `envvar:"MY_INT"`
	MyFloat             float64 `envvar:"MY_FLOAT"`
	MyString            string  `envvar:"MY_STRING"`
	MyBool              bool    `envvar:"MY_BOOL"`
}

// This main function shows how to use a Builder to create a config from env vars.
func main() {

	os.Setenv("MY_INT", "1234")
	os.Setenv("MY_FLOAT", "1.41421356237")
	os.Setenv("MY_STRING", "I have a cunning plan")
	os.Setenv("MY_BOOL", "true")

	// create a new config Builder
	builder := cfgbuild.Builder[*config]{}
	// build the new config setting the values from the env vars
	cfg, err := builder.Build()
	if err != nil {
		log.Fatal(err)
	}

	// print out the different values in the config
	fmt.Printf("cfg.MyInt = %v\ncfg.MyFloat = %v\ncfg.MyString = %v\ncfg.MyBool = %v",
		cfg.MyInt, cfg.MyFloat, cfg.MyString, cfg.MyBool)
}
