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
	defaultValue string
	omitEmpty    bool
}

func (tag confTag) getDefaultValue() (string, bool) {
	if len(tag.defaultValue) > 0 {
		return tag.defaultValue, true
	}

	return "", false
}

func parseTag(tag reflect.StructTag) (confTag, error) {
	val, exists := tag.Lookup(structTag)
	if !exists {
		return confTag{}, fmt.Errorf("no tag")
	}

	tags := strings.Split(val, ",")

	switch len(tags) {
	case 0:
		return confTag{}, nil
	case 1:
		return confTag{defaultValue: tags[0], omitEmpty: false}, nil
	default:
		return confTag{defaultValue: tags[0], omitEmpty: tags[1] == "omitempty"}, nil
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

func (p intParser) Parse(val *reflect.Value, typ *reflect.StructField) error {
	value, exists := os.LookupEnv(typ.Name)
	if exists {
		if typedVal, err := strconv.ParseInt(value, 10, 64); err != nil {
			return fmt.Errorf(typeMismatchError, typ.Name, err)
		} else {
			val.SetInt(typedVal)
		}

	} else {
		tag, err := parseTag(typ.Tag)
		if err != nil {
			return fmt.Errorf(noValueForFieldError, typ.Name)
		}

		value, exists := tag.getDefaultValue()
		if exists {
			if typedVal, err := strconv.ParseInt(value, 10, 64); err == nil {
				val.SetInt(typedVal)
			} else {
				return fmt.Errorf(typeMismatchError, typ.Name, err)
			}

		} else if !exists && !tag.omitEmpty {
			return fmt.Errorf(noValueForFieldError, typ.Name)
		}
	}

	return nil
}

type stringParser struct{}

func (p stringParser) Parse(val *reflect.Value, typ *reflect.StructField) error {
	value, exists := os.LookupEnv(typ.Name)
	if exists {
		val.SetString(value)

	} else {
		tag, err := parseTag(typ.Tag)
		if err != nil {
			return fmt.Errorf(noValueForFieldError, typ.Name)
		}

		value, exists := tag.getDefaultValue()
		if exists {
			val.SetString(value)

		} else if !exists && !tag.omitEmpty {
			return fmt.Errorf("no value for %s config field", typ.Name)
		}
	}

	return nil
}

type uintParser struct{}

func (p uintParser) Parse(val *reflect.Value, typ *reflect.StructField) error {
	value, exist := os.LookupEnv(typ.Name)
	if exist {
		if typedVal, err := strconv.ParseUint(value, 10, 64); err != nil {
			return fmt.Errorf(typeMismatchError, typ.Name, err)
		} else {
			val.SetUint(typedVal)
		}

	} else {
		tag, err := parseTag(typ.Tag)
		if err != nil {
			return fmt.Errorf(noValueForFieldError, typ.Name)
		}

		value, exists := tag.getDefaultValue()
		if exists {
			if typedVal, err := strconv.ParseUint(value, 10, 64); err != nil {
				return fmt.Errorf(typeMismatchError, typ.Name, err)
			} else {
				val.SetUint(typedVal)
			}
		} else if !exists && !tag.omitEmpty {
			return fmt.Errorf("no value for %s config field", typ.Name)
		}
	}

	return nil
}

type floatParser struct{}

func (p floatParser) Parse(val *reflect.Value, typ *reflect.StructField) error {
	value, exists := os.LookupEnv(typ.Name)
	if exists {
		if typedVal, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf(typeMismatchError, typ.Name, err)
		} else {
			val.SetFloat(typedVal)
		}
	} else {
		tag, err := parseTag(typ.Tag)
		if err != nil {
			return fmt.Errorf(noValueForFieldError, typ.Name)
		}

		value, exists := tag.getDefaultValue()
		if exists {
			if typedVal, err := strconv.ParseFloat(value, 64); err != nil {
				return fmt.Errorf(typeMismatchError, typ.Name, err)
			} else {
				val.SetFloat(typedVal)
			}

		} else if !exists && !tag.omitEmpty {
			return fmt.Errorf(noValueForFieldError, typ.Name)
		}
	}

	return nil
}
