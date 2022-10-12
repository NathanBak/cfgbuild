# cfgbuild
A Lightweight Golang Library for loading app configs

## Introduction
The purpose of cfgbuild is to provide an easy, lightweight package for loading application configuration settings.  It is able to build a struct and initialize the fields with associated environment variable values (see examples for loading from a .env file).  The main package does not have any external dependencies (but the tests and examples do reference external projects).  This library is published under the [BSD 2-Clause License](LICENSE) which provides a lot of flexiblity for usage.

## Example Config

A Config is just a struct with fields for different application settings.  Fields can be of most types that have a string representation.  Each field should have an `envvar` tag specifying the environment variable that will provide the value. 

```golang
// Config struct defines fields for application settings.
// It can be called "Config" or anything else.
type Config struct {
    // Field names can be anything and many different types
    // are supported.  The cfgBuild.Builder will use the 
    // envvar tag to know which environment variable to use.
	MyInt    int     `envvar:"MY_INT"`
	MyFloat  float64 `envvar:"MY_FLOAT"`
	MyString string  `envvar:"MY_STRING"`
	MyBool   bool    `envvar:"MY_BOOL"`
}

// Init is used to set default values.  This function should
// only be called by the cfgBuild.Builder.
func (cfg *Config) Init() error {
	cfg.MyInt = 1234
	return nil
}
```

## Usage
Building a config is pretty simple.  Basically, you just need to create a new builder (providing the Config type) and then run the `builder.Build()` function.  Here's some example code:
```golang
package main

import "github.com/NathanBak/cfgbuild"

func main() {
	builder := cfgbuild.Builder[*Config]{}
	cfg, err := builder.Build()
	// ...
    // Handle errors, use cfg, ... , profit
    //...
}

```

## Examples
The [examples](examples/) directory includes:
- [simple](examples/simple/) which shows a simple use case of loading a config from environment variables
- [fromdotenv](examples/fromdotenv/) which shows how to load a config from a `.env` file
- [bootstrap](examples/bootstrap/) which shows how to wrap a Builder into a Config constructor

## FAQ
 Q - Can this library read configuration information from .env files?<br>
A - The cfgbuild package does not know how to read .env packages, but can easily be paired with [godotenv](https://github.com/joho/godotenv).  The examples show two different ways to use godotenv.  Note that godetenv (created by John Barton) uses an [MIT License](https://github.com/joho/godotenv/blob/main/LICENCE).

Q - How does cfgbuild compare with [Viper](https://github.com/spf13/viper)?
<br>
A - Viper has a lot more whistles and bells and is overall more flexible and more powerful, but is also more complicated and pulls in numerous dependencies.  If you're looking for something to implement a complex CLI, then Viper is your friend.  But if you have a microservice that will be running in a container and needs an easy way to get configuration information, it's worth considering cfgbuild.