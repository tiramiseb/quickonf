package errors

import "errors"

// NoError means the step must be stopped but it is not an error
var NoError = errors.New("no-error")

// SkipNext means the next instruction must be skipped, but the rest of the step must still be executed, but it is not an error
var SkipNext = errors.New("skip-next")
