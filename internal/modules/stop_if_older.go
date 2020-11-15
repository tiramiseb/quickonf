package modules

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	quickonfErrors "github.com/tiramiseb/quickonf/internal/errors"
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("stop-if-older", StopIfOlder)
}

// StopIfOlder compares versions
func StopIfOlder(in interface{}, out output.Output) error {
	out.InstructionTitle("Comparing versions")
	data, err := helper.MapStringString(in)
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
		return fmt.Errorf("With candidate as \"%s\": %w", candStr, err)
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
