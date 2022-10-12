package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NathanBak/cfgbuild/examples/bootstrap/config"
)

// This main function uses a config that bootstraps itself to initialize the fields.
func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	// print cfg as JSON string
	buf, err := json.Marshal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))

	// print out the different values in the config
	fmt.Printf("cfg.MyInt = %v\ncfg.MyFloat = %v\ncfg.MyString = %v\n",
		cfg.MyInt, cfg.MyFloat, cfg.MyString)
}
