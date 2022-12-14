package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NathanBak/cfgbuild"
)

type config struct {
	MyString string `envvar:"MY_STRING"`
	// Color implements the TextUnmarshaler interface
	MyColor Color `envvar:"MY_COLOR"`
}

// This main function shows how to use a Builder to create a config from env vars.
// In this case, on of the config values is an enum.  The enum can be parsed because it
// implements the TextUnmarshaler interface.
//
// Expected output:
//
// cfg.MyString = Certainly not. We run.
// cfg.MyColor = Green
// cfg.MyColor enum value = 3
func main() {

	os.Setenv("MY_STRING", "Certainly not. We run.")
	os.Setenv("MY_COLOR", "grEEn")

	// create a new config Builder
	builder := cfgbuild.Builder[*config]{}
	// build the new config setting the values from the env vars
	cfg, err := builder.Build()
	if err != nil {
		log.Fatal(err)
	}

	// print out the different values in the config
	fmt.Printf("cfg.MyString = %v\ncfg.MyColor = %v\n",
		cfg.MyString, cfg.MyColor)

	// print out enum value of color (green is 3)
	fmt.Printf("cfg.MyColor enum value = %d", cfg.MyColor)
}
