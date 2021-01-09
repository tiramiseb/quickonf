// Package helper provides helpers for modules, in order to simplify their development
package helper

import "strings"

// ResultStatus is returned by helper to inform the module of the result status
type ResultStatus int

const (
	// ResultError means there is an error
	ResultError ResultStatus = iota
	// ResultSuccess means the requested action is successful
	ResultSuccess
	// ResultDryrun means the action is not done because of dry-run mode
	ResultDryrun
	// ResultAlready means the action was already done
	ResultAlready
)

// Dryrun allows running instructions without system modification
var Dryrun bool

func init() {
	cmdout, err := Exec(nil, "", "lsb_release", "--codename", "--short")
	if err != nil {
		panic(err)
	}
	Store("oscodename", strings.TrimSpace(string(cmdout)))
}
