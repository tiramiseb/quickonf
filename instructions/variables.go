package instructions

import (
	"bytes"
	"os"
	"strings"

	"github.com/tiramiseb/quickonf/commands/helper"
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
	if err := helper.Exec(nil, &buf, "lsb_release", "--codename", "--short"); err != nil {
		panic(err)
	}
	codename := strings.TrimSpace(buf.String())
	globalVars.define("oscodename", codename)
}

func NewGlobalVar(key, value string) {
	globalVars.define(key, value)
}

func newVariablesSet() Variables {
	return globalVars.clone()
}

func (v Variables) clone() Variables {
	newVars := map[string]string{}
	for key, value := range v {
		newVars[key] = value
	}
	return newVars
}

func (v Variables) define(key, val string) {
	v["<"+key+">"] = val
}

func (v Variables) TranslateVariables(src string) string {
	for key, val := range v {
		src = strings.ReplaceAll(src, key, val)
	}
	return src
}
