package goconfloader

import (
	"reflect"

	"github.com/joho/godotenv"
)

type ParsingPolicy byte

const (
	EscapeErrors = ParsingPolicy(iota)
	ErrorNoValue
)

var configPolicy = ErrorNoValue

func SetParsingPolicy(policy ParsingPolicy) {
	configPolicy = policy
}

func LoadConfig(config interface{}, envPath ...string) error {
	if len(envPath) > 0 {
		err := godotenv.Load(envPath...)
		if err != nil {
			return err
		}
	}

	val := reflect.ValueOf(config).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		structVal := val.Field(i)
		structTyp := typ.Field(i)

		parser := GetParser(&structVal)

		err := parser.Parse(&structVal, &structTyp)
		if err != nil && configPolicy == ErrorNoValue {
			return err
		}
	}

	return nil
}
