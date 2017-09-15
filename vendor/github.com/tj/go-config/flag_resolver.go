package config

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

// FlagResolver resolves configuration from command-line flags.
//
// For example LogLevel field would become -log-level=error.
type FlagResolver struct {
	Args []string // Args from the command-line
	set  *flag.FlagSet
}

// Name implementation.
func (f *FlagResolver) Name() string {
	return "flag"
}

// Setup implementation setting up the flag set.
func (f *FlagResolver) Setup() error {
	f.set = flag.NewFlagSet(f.Args[0], flag.ExitOnError)
	return nil
}

// Field implementation populating the flag set.
func (f *FlagResolver) Field(field Field) error {
	name := field.Name()
	help := field.Tag("help")
	addr := field.Interface()

	// TODO: need to fix the case where an undefined flag
	// is given, as it will output incomplete help
	// (or disable the help output)
	if !f.has(name) && !f.hasHelp() {
		return ErrFieldNotFound
	}

	switch addr.(type) {
	case *bool:
		v := addr.(*bool)
		f.set.BoolVar(v, name, *v, help)
	case *string:
		v := addr.(*string)
		f.set.StringVar(v, name, *v, help)
	case *int:
		v := addr.(*int)
		f.set.IntVar(v, name, *v, help)
	case *uint:
		v := addr.(*uint)
		f.set.UintVar(v, name, *v, help)
	case *float64:
		v := addr.(*float64)
		f.set.Float64Var(v, name, *v, help)
	case *time.Duration:
		v := addr.(*time.Duration)
		f.set.DurationVar(v, name, *v, help)
	case Value, *[]string:
		f.set.Var(field.Value(), name, help)
	}

	return nil
}

// Resolve implementation parsing the flag set.
func (f *FlagResolver) Resolve() error {
	err := f.set.Parse(f.Args[1:])

	if err == flag.ErrHelp {
		return nil
	}

	return err
}

// HasHelp checks for the presence of the --help or -h flags.
func (f *FlagResolver) hasHelp() bool {
	return f.has("help") || f.has("h")
}

// Has checks for a flag with the given `name`.
func (f *FlagResolver) has(name string) bool {
	// TODO: hacky, need to actually support -- etc,
	// would be nice if we could grab this from the flag
	// package's map but it seems inaccessible
	for _, arg := range f.Args[1:] {
		switch {
		case arg == fmt.Sprintf("--%s", name):
			return true
		case arg == fmt.Sprintf("-%s", name):
			return true
		case strings.HasPrefix(arg, fmt.Sprintf("--%s=", name)):
			return true
		case strings.HasPrefix(arg, fmt.Sprintf("-%s=", name)):
			return true
		}
	}

	return false
}
