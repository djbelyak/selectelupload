package config

import (
	"log"
	"os"
)

// DefaultResolvers used by Resolve() and MustResolve().
var DefaultResolvers = []Resolver{
	&FlagResolver{Args: os.Args},
	&EnvResolver{},
}

// Resolve `options` using the built-in flag and env resolvers.
func Resolve(options interface{}) error {
	c := Config{
		Options:   options,
		Resolvers: DefaultResolvers,
	}

	return c.Resolve()
}

// MustResolve `options` using the built-in flag and env resolvers.
func MustResolve(options interface{}) {
	if err := Resolve(options); err != nil {
		log.Fatalf("error resolving configuration: %s", err)
	}
}
