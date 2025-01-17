# goconfloader

 GoConfLoader is a lightweight configuration loader for Go applications. It supports loading configuration from environment variables and .env files with help of github.com/joho/godotenv. This library is designed to simplify the configuration management process in Go applications by providing a unified way to handle configuration data.

 Now supports Int, String, Uint, Float and any of types which have those types as underlying types.

 Bool, Slices, Map, Struct and !!!TESTS!!! will be added in future.

 ## Installation
 ```shell
 go get github.com/LSD2409/goconfloader
 ```
 ## Usage
 As joho/godotenv add app configuration to .env file 
 ```shell
 ENV_VAR1=SOMEVALUE
 ENV_VAR2=ANOTHERVALUE
 ```
 After all done, create a struct with fields you want to parse from .env and variable to store config. 
 ```go
 package main

 import (
    "fmt"

    "github.com/LSD2409/goconfloader"
 )
 
 // struct with config fields
 type struct Config {
    EnvVar1 string `ConfLoader:"ENV_VAR1"`  // Alias
    EnvVar2 string `ConfLoader: "ENV_VAR2"` // Alias
 }

 // variable where you want to store config
 var AppConfig Config

 func main() {
    // provide pointer to you conf variable and pathes to .env files, also you can skip path to env if you env variables already loaded
    err := goconfloader.LoadConfig(&AppConfig, "pathToEnv")
    if err != nil {
        panic(err.Error())
    }

    fmt.Println(AppConfig.EnvVar1)
    fmt.Println(AppConfig.EnvVar2)
 }
 ```

 The result is 
 ```shell
 SOMEVALUE
 ANOTHERVALUE
 ```

 ## Default values
 You can provide defualt values with struct tag ConfLoader, if you env don't have any variables for some reason.
 ```go
package main

 import (
    "fmt"

    "github.com/LSD2409/goconfloader"
 )
 
 // add default values with ConfLoader tag
 type struct Config {
    EnvVar1 string `ConfLoader:"ENV_VAR1,defaultValue"` // alias and default value
    EnvVar2 string `ConfLoader:"ENV_VAR2,oneMoreDefaultValue"` // alias and default value
 }

 var AppConfig Config

 func main() {
    err := goconfloader.LoadConfig(&AppConfig, "pathToEnv")
    if err != nil {
        panic(err.Error())
    }

    fmt.Println(AppConfig.EnvVar1)
    fmt.Println(AppConfig.EnvVar2)
 }
 ```
 ```shell
 defaultValue
 oneMoreDefaultValue
 ```

 ## Underlying types
 You can use any type which underlying type is supported.
 ```go
 package main

import (
	"fmt"

	"github.com/LSD2409/goconfloader"
)

type T string

type Config struct {
	EnvVar1 T `ConfLoader:"ENV_VAR1"`
	EnvVar2 T `ConfLoader:"ENV_VAR2"`
}

var AppConfig Config

func main() {
	err := goconfloader.LoadConfig(&AppConfig, "pathToEnv")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(AppConfig.EnvVar1)
	fmt.Println(AppConfig.EnvVar2)
}
 ```
