// Package helper provides helpers for modules, in order to simplify their development
package helper

import "strings"

// Dryrun allows running instructions without system modification
var Dryrun = false

func init() {
	cmdout, err := Exec(nil, "lsb_release", "--codename", "--short")
	if err != nil {
		panic(err)
	}
	Store("oscodename", strings.TrimSpace(string(cmdout)))
}
