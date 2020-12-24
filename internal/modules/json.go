package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("json-build", JSONBuild)
	Register("json-get", JSONGet)
}

// JSONBuild builds a (simple) json structure from flat information
func JSONBuild(in interface{}, out output.Output) error {
	out.InstructionTitle("Build a JSON")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	storeKeys := []string{}
	source := map[string]interface{}{}
	for _, instr := range data {
		details := strings.SplitN(instr, "=", 2)
		if len(details) < 2 || details[0] == "" {
			// Ignore any line which does not contains a "=" or with an empty key
			continue
		}
		if details[0] == "store" {
			// If key is "store", use it as a store key
			storeKeys = append(storeKeys, details[1])
			continue
		}
		source[details[0]] = details[1]
	}
	result, err := json.Marshal(source)
	if err != nil {
		return err
	}
	resultS := string(result)
	out.Info(fmt.Sprintf("Result is %s", resultS))
	for _, k := range storeKeys {
		helper.Store(k, resultS)
	}
	return nil

}

// JSONGet returns a single value from a JSON object
func JSONGet(in interface{}, out output.Output) error {
	out.InstructionTitle("Get a JSON value")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	from, ok := data["from"]
	if !ok {
		return errors.New("Missing from")
	}
	key, ok := data["key"]
	if !ok {
		return errors.New("Missing key")
	}

	result := gjson.Get(from, key).String()
	out.Info(fmt.Sprintf("Got %s", result))

	store, ok := data["store"]
	if ok {
		helper.Store(store, result)
	}
	return nil
}
