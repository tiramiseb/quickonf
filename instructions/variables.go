package instructions

import (
	"bytes"
	"os"
	"strings"

	"github.com/tiramiseb/quickonf/commands/helper"
)

type Variables struct {
	names   []string
	mapping map[string]string
}

var globalVars = &Variables{
	mapping: map[string]string{},
}

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

func GlobalVars() []string {
	return globalVars.names
}

func newVariablesSet() *Variables {
	return globalVars.clone()
}

func (v *Variables) clone() *Variables {
	newVars := &Variables{
		names:   make([]string, len(v.names)),
		mapping: map[string]string{},
	}
	copy(newVars.names, v.names)
	for k, v := range v.mapping {
		newVars.mapping[k] = v
	}
	return newVars
}

func (v *Variables) define(key, val string) {
	v.mapping["<"+key+">"] = val
	v.names = append(v.names, key)
}

func (v *Variables) TranslateVariables(src string) string {
	for key, val := range v.mapping {
		src = strings.ReplaceAll(src, key, val)
	}
	return src
}
