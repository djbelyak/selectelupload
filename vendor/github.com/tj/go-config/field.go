package config

import (
	"fmt"
	"reflect"
)

// Field represents a struct field which is passed to resolvers for lookup.
type Field interface {
	// Name returns the field's name. The name is derived either from
	// the "name" tag, or via reflection. Nested structs inherit the
	// parent field's name.
	Name() string

	// Interface of the field's pointer.
	Interface() interface{}

	// Value representation of the field's value.
	Value() Value

	// Tag returns the field's tag via its `name`.
	Tag(name string) string
}

// Field implementation.
type field struct {
	value  reflect.Value
	field  reflect.StructField
	parent *field
}

// Interface implementation.
func (f *field) Interface() interface{} {
	return f.value.Addr().Interface()
}

// Name implementation.
func (f *field) Name() string {
	s := f.Tag("name")

	if s == "" {
		s = normalizeName(f.field.Name)
	}

	if f.parent != nil {
		s = f.parent.prefix(s)
	}

	return s
}

// Tag implementation.
func (f *field) Tag(name string) string {
	return f.field.Tag.Get(name)
}

// Value implementation.
func (f *field) Value() Value {
	return valueOf(f.Interface())
}

// Prefixed name.
func (f *field) prefix(name string) string {
	return fmt.Sprintf("%s-%s", f.Name(), name)
}
