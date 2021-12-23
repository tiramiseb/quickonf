package instructions

import (
	"bytes"
	"os"
	"strings"

	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

type Variables map[string]string

var globalVars = Variables{}

func init() {
	// Initialize global variables, required by some instructions
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	globalVars.define("hostname", hostname)

	var buf bytes.Buffer
	wait, err := helper.Exec(nil, &buf, "lsb_release", "--codename", "--short")
	if err != nil {
		panic(err)
	}
	wait()
	codename := strings.TrimSpace(buf.String())
	globalVars.define("oscodename", codename)
}

func NewVariablesSet() Variables {
	v := Variables{}
	for key, val := range globalVars {
		v[key] = val
	}
	return v
}

func (v Variables) define(key, val string) {
	v["<"+key+">"] = val
}

func (v Variables) translateVariables(src string) string {
	for key, val := range v {
		src = strings.ReplaceAll(src, key, val)
	}
	return src
}
