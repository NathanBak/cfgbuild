# cfgbuild
A Lightweight Golang Library for loading app configs

## Introduction
The purpose of cfgbuild is to provide an easy, lightweight package for loading application configuration settings.  It is able to build a struct and initialize the fields with associated environment variable values (see examples for loading from a .env file).  The main package does not have any external dependencies (but the tests and examples do require external projects).  This library is published under the [BSD 2-Clause License](LICENSE) which provides a lot of flexiblity for usage.

## Example Config

A Config is just a struct with fields for different application settings.  Fields can be of most types that have a string representation.  Each field should have an `envvar` tag specifying the environment variable that will provide the value. 

```golang
// Config struct defines fields for application settings.
// It can be called "Config" or anything else.
type Config struct {
    // Field names can be anything and many different types
    // are supported.  The cfgBuild.Builder will use the 
    // envvar tag to know which environment variable to use.
	// If the tag also contains the "required" flag, then 
	// calls to Build() will fail if the value has not been
	// set.  If the tag contains the "default" flag then the
	// provided value will be applied as the default value.
	MyInt    int     `envvar:"MY_INT,required"`
	MyFloat  float64 `envvar:"MY_FLOAT"`
	MyString string  `envvar:"MY_STRING"`
	MyBool   bool    `envvar:"MY_BOOL,default=TRUE"`
}
```

## Usage
Building a config is pretty simple.  Basically, you just need to create a new builder (providing the Config type) and then run the `Builder.Build()` function.  Here's some example code:
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

## Tags

### Name
Tags typically follow the format of 
```golang
MyString string `envvar:"ENV_VAR_NAME"`
```
where `envvar` is a keyword and "ENV_VAR_NAME" is the name of the environment variable containing the value.

### Required
If a value is required (must be set), the `required` flag can be added.
```golang
MyString string `envvar:"ENV_VAR_NAME,required"`
```
When the required flag is included in the `envvar` tag, the cfgbuild.Builder.Build() function will return an error if the value is not set.

### Default
If there is a default value, it can be set using the `default` flag.
```golang
MyNumber int `envvar:"ENV_VAR_NAME,default=1234"`
```
If the environment variable is not set then the default value specified after the equals sign will be used instead.  There is no compile time validation of the default value so if an integer has something like `default=abc` specified then the cfgbuild.Builder.Build() function will always fail.

## Functions
Additional flexibility and customization can be achieved by adding implementations of specific functions to the Config struct.

### CfgBuildInit()
The CfgBuildInit() function can be used to perform any special initialization logic.  This can include things such as specifying complex default values or initializing special fields.  The function should have a signature like `func (cfg *Config) CfgBuildInit() error`.  It will be invoked during the Build() right after the new instance is created.

### CfgBuildValidate()
The CfgBuildValidate() function can be used to perform special validation the config.  This can include things such as verifying that set values are within certain ranges.  The function should have a signature like `func (cfg *Config) CfgBuildInit() error`.  It will be invoked as the final step during the Build().

## Examples
The [examples](examples/) directory includes:
- [simple](examples/simple/) which shows a simple use case of loading a config from environment variables
- [fromdotenv](examples/fromdotenv/) which shows how to load a config from a `.env` file
- [bootstrap](examples/bootstrap/) which shows how to wrap a Builder into a Config constructor
- [enumparse](examples/enumparse/) which shows how a config field can be an enum

## FAQ
 Q - Can this library read configuration information from .env files?<br>
A - The cfgbuild package does not know how to read .env packages, but can easily be paired with [godotenv](https://github.com/joho/godotenv).  The examples show two different ways to use godotenv.  Note that godetenv (created by John Barton) uses an [MIT License](https://github.com/joho/godotenv/blob/main/LICENCE).

Q - How does cfgbuild compare with [Viper](https://github.com/spf13/viper)?
<br>
A - Viper has a lot more whistles and bells and is overall more flexible and more powerful, but is also more complicated and pulls in numerous dependencies.  If you're looking for something to implement a complex CLI, then Viper is your friend.  But if you have a microservice that will be running in a container and needs an easy way to get configuration information, it's worth considering cfgbuild.

Q - Does cfgbuild support enums?
<br>
A - A config can have an enum field if the enum implements the [TextUnmarshaler interface](https://pkg.go.dev/encoding#TextUnmarshaler).  See the [Color](examples/enumparse/color.go) enum for an example.

