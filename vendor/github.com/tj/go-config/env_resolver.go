package config

import (
	"os"
	"strings"
)

// EnvResolver resolves configuration from environment variables.
//
// For example LogLevel field would become LOG_LEVEL=error.
type EnvResolver struct {
	// Prefix optionally applied to each lookup. Omit the
	// trailing "_", this is applied automatically.
	Prefix string
}

// Name implementation.
func (e *EnvResolver) Name() string {
	return "env"
}

// Setup implementation (temporary noop).
func (e *EnvResolver) Setup() error {
	return nil
}

// Field implementation normalizing the field name
// and performing coercion to the field type.
func (e *EnvResolver) Field(field Field) error {
	name := field.Name()
	s := os.Getenv(e.envize(name))

	if s == "" {
		return ErrFieldNotFound
	}

	return field.Value().Set(s)
}

// Resolve implementation (temporary noop).
func (*EnvResolver) Resolve() error {
	return nil
}

// Normalize `name` with prefix support.
func (e *EnvResolver) envize(name string) string {
	if e.Prefix != "" {
		return e.normalize(e.Prefix) + "_" + e.normalize(name)
	}

	return e.normalize(name)
}

// Normalize `name`.
func (*EnvResolver) normalize(name string) string {
	return strings.ToUpper(strings.Replace(name, "-", "_", -1))
}
