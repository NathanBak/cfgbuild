package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NathanBak/cfgbuild/examples/bootstrap/config"
)

// This main function uses a config that bootstraps itself to initialize the fields.
//
// Expected output:
//
// {"myInt":8081,"myFloat":3.14159,"myString":"Smoke me a kipper, I'll be back for breakfast!","myBool":false}
// cfg.MyInt = 8081
// cfg.MyFloat = 3.14159
// cfg.MyString = Smoke me a kipper, I'll be back for breakfast!
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
