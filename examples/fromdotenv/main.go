package main

import (
	"fmt"
	"log"

	"github.com/NathanBak/cfgbuild"

	// The import below will automatically find and load the .env file.  Refer to
	// https://github.com/joho/godotenv/blob/main/autoload/autoload.go for details.
	_ "github.com/joho/godotenv/autoload"
)

// This main function shows how to use a Builder to create a config with the values initialized.
func main() {
	// create a new config Builder
	builder := cfgbuild.Builder[*Config]{}
	// build the new config setting the values from the env vars
	cfg, err := builder.Build()
	if err != nil {
		log.Fatal(err)
	}

	// print out the different values in the config
	fmt.Printf("cfg.MyInt = %v\ncfg.MyFloat = %v\ncfg.MyDuration = %v\n",
		cfg.MyInt, cfg.MyFloat, cfg.MyDuration)
	fmt.Printf("cfg.MyTime = %v\ncfg.MyBytes = %v\ncfg.MyString = %v\n",
		cfg.MyTime, cfg.MyBytes, cfg.MyString)
	fmt.Printf("cfg.MyBool = %v\n", cfg.MyBool)
}
