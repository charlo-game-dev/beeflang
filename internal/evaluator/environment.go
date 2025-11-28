package evaluator

import "github.com/elitwilson/beeflang/internal/object"

// Re-export Environment types from object package for backward compatibility
type Environment = object.Environment

var NewEnvironment = object.NewEnvironment
var NewEnclosedEnvironment = object.NewEnclosedEnvironment
