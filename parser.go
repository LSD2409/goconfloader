package goconfloader

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	typeMismatchError    = "type mismatch while parsing %s field %w"
	noValueForFieldError = "no value for %s config field"
	structTag            = "ConfLoader"
)

type ConfigParser interface {
	Parse(val *reflect.Value, typ *reflect.StructField) error
}

type confTag struct {
	alias        string
	defaultValue string
}

func getTag(tag reflect.StructTag) (confTag, bool) {
	val, exists := tag.Lookup(structTag)
	if !exists {
		return confTag{}, false
	}

	tags := strings.Split(val, ",")

	switch len(tags) {
	case 0:
		return confTag{}, true
	case 1:
		return confTag{alias: tags[0], defaultValue: ""}, true
	default:
		return confTag{alias: tags[0], defaultValue: tags[1]}, true
	}
}

func GetParser(val *reflect.Value) ConfigParser {
	switch kind := val.Kind(); kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intParser{}
	case reflect.String:
		return stringParser{}
	case reflect.Float32, reflect.Float64:
		return floatParser{}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintParser{}
	default:
		return stringParser{}
	}
}

type intParser struct{}

func getValue(typ *reflect.StructField) (string, error) {
	var envName string

	tag, tagExists := getTag(typ.Tag)
	if tagExists && tag.alias != "" {
		envName = tag.alias
	} else {
		envName = typ.Name
	}

	if value, exists := os.LookupEnv(envName); exists {
		return value, nil
	} else {
		if tagExists && len(tag.defaultValue) > 0 {
			return tag.defaultValue, nil
		} else {
			return "", fmt.Errorf(noValueForFieldError, envName)
		}
	}
}

func (p intParser) Parse(val *reflect.Value, typ *reflect.StructField) error {
	value, err := getValue(typ)
	if err != nil {
		return err
	}

	typedVal, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf(typeMismatchError, typ.Name, err)
	} else {
		val.SetInt(typedVal)
	}

	return nil
}

type stringParser struct{}

func (p stringParser) Parse(val *reflect.Value, typ *reflect.StructField) error {
	value, err := getValue(typ)
	if err != nil {
		return err
	}

	val.SetString(value)

	return nil
}

type uintParser struct{}

func (p uintParser) Parse(val *reflect.Value, typ *reflect.StructField) error {
	value, err := getValue(typ)
	if err != nil {
		return err
	}

	typedVal, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fmt.Errorf(typeMismatchError, typ.Name, err)
	}

	val.SetUint(typedVal)

	return nil
}

type floatParser struct{}

func (p floatParser) Parse(val *reflect.Value, typ *reflect.StructField) error {
	value, err := getValue(typ)
	if err != nil {
		return err
	}

	typedVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf(typeMismatchError, typ.Name, err)
	}

	val.SetFloat(typedVal)
	return nil
}
