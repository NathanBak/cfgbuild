package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NathanBak/cfgbuild"
)

type config struct {
	MyInt    int     `envvar:"MY_INT"`
	MyFloat  float64 `envvar:"MY_FLOAT"`
	MyString string  `envvar:"MY_STRING"`
	MyBool   bool    `envvar:"MY_BOOL"`
}

// This main function shows how to use a Builder to create a config from env vars.
//
// Expected output:
//
// cfg.MyInt = 1234
// cfg.MyFloat = 1.41421356237
// cfg.MyString = I have a cunning plan
// cfg.MyBool = true
func main() {

	os.Setenv("MY_INT", "1234")
	os.Setenv("MY_FLOAT", "1.41421356237")
	os.Setenv("MY_STRING", "I have a cunning plan")
	os.Setenv("MY_BOOL", "true")

	// create a new config and initialize from the env vars
	cfg, err := cfgbuild.NewConfig[*config]()
	if err != nil {
		log.Fatal(err)
	}

	// print out the different values in the config
	fmt.Println("cfg.MyInt =", cfg.MyInt)
	fmt.Println("cfg.MyFloat =", cfg.MyFloat)
	fmt.Println("cfg.MyString =", cfg.MyString)
	fmt.Println("cfg.MyBool =", cfg.MyBool)
}
