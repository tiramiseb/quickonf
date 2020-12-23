package modules

import (
	"errors"
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"

	quickonfErrors "github.com/tiramiseb/quickonf/internal/errors"
	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	Register("stop-if-exist", StopIfExist)
	Register("stop-if-older", StopIfOlder)
}

// StopIfExist stops the step if all given files exist
func StopIfExist(in interface{}, out output.Output) error {
	out.InstructionTitle("Stop if exist")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, f := range data {
		f = helper.Path(f)
		_, err := os.Stat(f)
		if err != nil {
			if os.IsNotExist(err) {
				out.Info(fmt.Sprintf("File %s does not exist", f))
				return nil
			}
			return err
		}
	}
	out.Info(fmt.Sprintf("All listed files exist"))
	return quickonfErrors.NoError
}

// StopIfOlder compares versions and stops the step if older
func StopIfOlder(in interface{}, out output.Output) error {
	out.InstructionTitle("Stop if older")
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
