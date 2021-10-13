package modules

import (
	"strconv"
	"strings"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("dconf", Dconf)
}

func interfaceToDconfValue(src interface{}) string {
	if asBool, ok := src.(bool); ok {
		return strconv.FormatBool(asBool)
	}
	if asInt, ok := src.(int); ok {
		return strconv.Itoa(asInt)
	}
	if asFloat, ok := src.(float64); ok {
		return strconv.FormatFloat(asFloat, 'f', -1, 64)
	}
	if asString, ok := src.(string); ok {
		if strings.HasPrefix(asString, "uint32:") {
			return asString[7:]
		}
		return "\"" + strings.ReplaceAll(asString, "\"", "\\\"") + "\""
	}
	if asSlice, ok := src.([]interface{}); ok {
		if len(asSlice) == 0 {
			return "@as []"
		}
		dest := make([]string, len(asSlice))
		for i, val := range asSlice {
			dest[i] = interfaceToDconfValue(val)
		}
		return "[" + strings.Join(dest, ", ") + "]"
	}
	return ""
}

// Dconf sets a parameter in the dconf database
func Dconf(in interface{}, out output.Output) error {
	out.InstructionTitle("Dconf database")
	data, err := helper.MapStringInterface(in)
	if err != nil {
		return err
	}
	for k, v := range data {
		if Dryrun {
			out.Infof("Would set %s to %s", k, v)
			continue
		}
		val := interfaceToDconfValue(v)
		if val == "" {
			out.Infof("Could not understand %#v", v)
			continue
		}
		if _, err := helper.Exec(nil, "", "dconf", "write", k, val); err != nil {
			return err
		}
		out.Successf("Set %s to %s", k, val)
	}
	return nil
}
