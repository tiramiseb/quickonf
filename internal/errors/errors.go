package errors

import "errors"

// NoError means the step must be stopped but it is not an error
var NoError = errors.New("no-error")
