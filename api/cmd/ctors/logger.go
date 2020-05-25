package ctors

import (
	"log"
	"os"
)

// NewLogger constructs a logger. It's just a regular Go function, without any
// special relationship to Fx.
//
// Since it returns a *log.Logger, Fx will treat NewLogger as the constructor
// function for the standard library's logger. (We'll see how to integrate
// NewLogger into an Fx application in the main function.) Since NewLogger
// doesn't have any parameters, Fx will infer that loggers don't depend on any
// other types - we can create them from thin air.
//
// Fx calls constructors lazily, so NewLogger will only be called only if some
// other function needs a logger. Once instantiated, the logger is cached and
// reused - within the application, it's effectively a singleton.
//
// By default, Fx applications only allow one constructor for each type. See
// the documentation of the In and Out types for ways around this restriction.
func NewLogger() *log.Logger {
	logger := log.New(os.Stdout, "" /* prefix */, 0 /* flags */)
	return logger
}
