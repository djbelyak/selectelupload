package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

// Value represents a custom type which may be represented as text,
// primarily used within resolvers which are string-based such as
// environment variables and flags.
//
// For example the Value interface may be used to implement a
// custom comma-delimited string slice.
type Value interface {
	String() string
	Set(string) error
}

// Bool value.
type boolValue bool

// Set value from string.
func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)
	return err
}

// String representatio
func (b *boolValue) String() string {
	return fmt.Sprint(*b)
}

// Int value.
type intValue int

// Set value from string.
func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = intValue(v)
	return err
}

// String representation.
func (i *intValue) String() string {
	return fmt.Sprint(*i)
}

// Uint value.
type uintValue uint

// Set value from string.
func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	*i = uintValue(v)
	return err
}

// String representation.
func (i *uintValue) String() string {
	return fmt.Sprint(*i)
}

// Float value.
type floatValue float64

// Set value from string.
func (f *floatValue) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = floatValue(v)
	return err
}

// String representation.
func (f *floatValue) String() string {
	return fmt.Sprint(*f)
}

// String value.
type stringValue string

// Set value from string.
func (s *stringValue) Set(v string) error {
	*s = stringValue(v)
	return nil
}

// String representation.
func (s *stringValue) String() string {
	return fmt.Sprint(*s)
}

// Duration value.
type durationValue time.Duration

// Set value from string.
func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(s)
	*d = durationValue(v)
	return err
}

// String representation.
func (d *durationValue) String() string {
	return fmt.Sprint(*d)
}

// Strings value..
type stringsValue []string

// Set value from string.
func (s *stringsValue) Set(v string) error {
	for _, part := range strings.Split(v, ",") {
		*s = append(*s, part)
	}
	return nil
}

// String representation.
func (s *stringsValue) String() string {
	return strings.Join(*s, ",")
}

// Bytes value.
type Bytes uint64

// Set value from string.
func (b *Bytes) Set(s string) error {
	v, err := humanize.ParseBytes(s)
	*b = Bytes(v)
	return err
}

// String representation.
func (b *Bytes) String() string {
	return humanize.Bytes(uint64(*b))
}

// ParseBytes is a utility function to parse initial Bytes values.
//
// This function panics if parsing fails.
func ParseBytes(s string) (b Bytes) {
	err := b.Set(s)
	if err != nil {
		panic(err)
	}

	return
}

// ValueOf returns the Value for `v` or panics.
func valueOf(v interface{}) Value {
	switch v.(type) {
	case *bool:
		return (*boolValue)(v.(*bool))
	case *int:
		return (*intValue)(v.(*int))
	case *uint:
		return (*uintValue)(v.(*uint))
	case *float64:
		return (*floatValue)(v.(*float64))
	case *string:
		return (*stringValue)(v.(*string))
	case *time.Duration:
		return (*durationValue)(v.(*time.Duration))
	case *[]string:
		return (*stringsValue)(v.(*[]string))
	case Value:
		return v.(Value)
	default:
		panic(fmt.Errorf("cannot convert %#v to a Value", v))
	}
}
