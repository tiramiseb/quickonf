package modules

import (
	"errors"
    "fmt"

	"github.com/Masterminds/semver/v3"

	quickonfErrors "github.com/tiramiseb/quickonf/errors"
	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("stop-if-older", StopIfOlder)
}

func StopIfOlder(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Comparing versions")
	data, err := input.MapStringString(in, store)
	if err != nil {
		return err
	}
	curStr, ok := data["current"]
	if !ok {
		return errors.New("Missing current version")
	}
    if curStr == "" {
        out.Info("Current version is empty")
        return nil
    }
	candStr, ok := data["candidate"]
	if !ok {
		return errors.New("Missing candidate version")
	}
	curVersion, err := semver.NewVersion(curStr)
	if err != nil {
		return fmt.Errorf("With current as \"%s\": %w", curStr, err)
	}
	candVersion, err := semver.NewVersion(candStr)
	if err != nil {
		return fmt.Errorf("With candidate as \"%s\": %w", curStr, err)
	}
	diff := candVersion.Compare(curVersion)
	switch diff {
	case -1:
		out.Info("Candidate (" + candStr + ") is older than current version (" + curStr + ")")
		return quickonfErrors.NoError
	case 0:
		out.Info("Candidate (" + candStr + ") is the same as current version (" + curStr + ")")
		return quickonfErrors.NoError
	case 1:
		out.Info("Candidate (" + candStr + ") is newer than current version (" + curStr + ")")
	}
	return nil
}
