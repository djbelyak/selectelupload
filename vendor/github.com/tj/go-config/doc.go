/*
Package config provides an API for resolving configuration
values from structs, with extensible resolvers, type coercion and validation.

Out of the box FlagResolver and EnvResolver are provided,
however you may provide your own by implementing the
Resolver interface.

Each field may have a "name" tag, which is otherwise derived
from field, a "help" tag used to describe the field, and a "validate"
tag which utilizes https://gopkg.in/validator.v2 under the hood for validation.

Defaults are provided by the initialized struct.
*/
package config
