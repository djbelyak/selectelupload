package config

// Resolver represents a struct which "resolves" configuration fields,
// for example flag, environment variables, key/value stores and so on.
//
// Warning: this interface will be nicer in the near future, the current
// interface is limited due to a limitation with the stdlib flag package.
type Resolver interface {
	// Name of the resolver such as "env", or "flag".
	Name() string

	// Field attempts to reoslve a field; this method should replace the
	// field's value by using its pointer via Field.Interface(), or return
	// ErrFieldNotFound.
	Field(Field) error

	// Setup (temporary method, don't get attached ;D).
	Setup() error

	// Resolve (temporary method, don't get attached ;D).
	Resolve() error
}
