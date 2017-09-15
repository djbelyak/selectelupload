package config

import (
	"errors"
	"reflect"
	"strings"

	"github.com/segmentio/go-snakecase"
	"gopkg.in/validator.v2"
)

// Errors used during resolution.
var (
	ErrFieldNotFound = errors.New("field not found")
)

// Config resolves options from the provided struct
// using one or more Resolvers.
type Config struct {
	// Options struct.
	Options interface{}

	// Resolvers list; the ordering is significant,
	// as it defines precedence. The first resolver
	// is used for all fields unless the "from" tag
	// of a field indictates otherwise.
	Resolvers []Resolver
}

// Resolve the configuration.
func (c *Config) Resolve() error {
	if c.Options == nil {
		return errors.New("Config.Options required")
	}

	if len(c.Resolvers) == 0 {
		return errors.New("Config.Resolvers required")
	}

	if reflect.ValueOf(c.Options).Kind() != reflect.Ptr {
		return errors.New("Config.Options must be a pointer")
	}

	for _, resolver := range c.Resolvers {
		err := resolver.Setup()
		if err != nil {
			return err
		}
	}

	val := reflect.ValueOf(c.Options).Elem()
	err := c.resolveStruct(val, nil)
	if err != nil {
		return err
	}

	for _, resolver := range c.Resolvers {
		err := resolver.Resolve()
		if err != nil {
			return err
		}
	}

	return validator.Validate(val.Interface())
}

// Resolve fields in the given struct.
func (c *Config) resolveStruct(val reflect.Value, parent *field) error {
	val = reflect.Indirect(val)
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		f := &field{
			value:  reflect.Indirect(val.Field(i)),
			field:  typ.Field(i),
			parent: parent,
		}

		from := f.Tag("from")

		for j, resolver := range c.Resolvers {
			if from == "*" || listed(from, resolver.Name()) || (j == 0 && from == "") {
				if f.value.Kind() == reflect.Struct {
					err := c.resolveStruct(f.value, f)
					if err != nil {
						return err
					}
				}

				err := resolver.Field(f)

				if err == ErrFieldNotFound {
					continue
				}

				if err != nil {
					return err
				}

				break
			}
		}
	}

	return nil
}

// Listed checks if `name` is within the comma-delimited `list` string.
func listed(list, name string) bool {
	for _, s := range strings.Split(list, ",") {
		if s == name {
			return true
		}
	}
	return false
}

// Normalize field name string, for example transforming LogLevel to log-level.
func normalizeName(s string) string {
	return strings.Replace(snakecase.Snakecase(s), "_", "-", -1)
}
