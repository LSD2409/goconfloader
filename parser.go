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
	parse(val *reflect.Value, typ *reflect.StructField) error
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

func parseValue(value string, typ *reflect.StructField) (reflect.Value, error) {
	switch typ.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf(typeMismatchError, typ.Name, err)
		}

		return reflect.ValueOf(intVal).Convert(reflect.TypeOf(int(0))), nil
	case reflect.String:
		return reflect.ValueOf(value).Convert(reflect.TypeOf("")), nil
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf(typeMismatchError, typ.Name, err)
		}

		return reflect.ValueOf(floatVal).Convert(reflect.TypeOf(float64(0))), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf(typeMismatchError, typ.Name, err)
		}

		return reflect.ValueOf(uintVal).Convert(reflect.TypeOf(uint(0))), nil
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return reflect.Value{}, fmt.Errorf(typeMismatchError, typ.Name, err)
		}

		return reflect.ValueOf(boolVal).Convert(reflect.TypeOf(true)), nil
	default:
		return reflect.Value{}, fmt.Errorf("unsupported type")
	}
}

func parseField(val *reflect.Value, typ *reflect.StructField) error {
	value, err := getValue(typ)
	if err != nil {
		return err
	}

	parsedValue, err := parseValue(value, typ)
	if err != nil {
		return err
	}
	val.Set(parsedValue)

	return nil
}
