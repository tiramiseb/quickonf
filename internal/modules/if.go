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
	Register("stop-if-equal", StopIfEqual)
	Register("skip-next-if-exist", SkipNextIfExist)
	Register("skip-next-if-older", SkipNextIfOlder)
	Register("skip-next-if-equal", SkipNextIfEqual)
}

func exists(in interface{}, out output.Output) (bool, error) {
	data, err := helper.SliceString(in)
	if err != nil {
		return false, err
	}
	for _, f := range data {
		f = helper.Path(f)
		_, err := os.Stat(f)
		if err != nil {
			if os.IsNotExist(err) {
				out.Infof("File %s does not exist", f)
				return false, nil
			}
			return false, err
		}
	}
	out.Successf("All listed files exist")
	return true, err

}

// StopIfExist stops the step if all given files exist
func StopIfExist(in interface{}, out output.Output) error {
	out.InstructionTitle("Stop if exist")
	ok, err := exists(in, out)
	if err != nil {
		return err
	}
	if ok {
		return quickonfErrors.NoError
	}
	return nil
}

// SkipNextIfExist skips the next instruction if all given files exist
func SkipNextIfExist(in interface{}, out output.Output) error {
	out.InstructionTitle("Skip next if exist")
	ok, err := exists(in, out)
	if err != nil {
		return err
	}
	if ok {
		return quickonfErrors.SkipNext
	}
	return nil
}

func older(in interface{}, out output.Output) (bool, error) {
	data, err := helper.MapStringString(in)
	if err != nil {
		return false, err
	}
	curStr, ok := data["current"]
	if !ok {
		return false, errors.New("missing current version")
	}
	if curStr == "" {
		out.Info("Current version is empty")
		return false, nil
	}
	candStr, ok := data["candidate"]
	if !ok {
		return false, errors.New("missing candidate version")
	}
	curVersion, err := semver.NewVersion(curStr)
	if err != nil {
		return false, fmt.Errorf(`with current as "%s": %w`, curStr, err)
	}
	candVersion, err := semver.NewVersion(candStr)
	if err != nil {
		return false, fmt.Errorf(`with candidate as "%s": %w`, candStr, err)
	}
	diff := candVersion.Compare(curVersion)
	switch diff {
	case -1:
		out.Infof("Candidate (%s) is older than current version (%s)", candStr, curStr)
		return true, nil
	case 0:
		out.Infof("Candidate (%s) is the same as current version (%s)", candStr, curStr)
		return true, nil
	case 1:
		out.Infof("Candidate (%s) is newer than current version (%s)", candStr, curStr)
	}
	return false, nil
}

// StopIfOlder compares versions and stops the step if older
func StopIfOlder(in interface{}, out output.Output) error {
	out.InstructionTitle("Stop if older")
	ok, err := older(in, out)
	if err != nil {
		return err
	}
	if ok {
		return quickonfErrors.NoError
	}
	return nil
}

// SkipNextIfOlder compares versions and skips the next instruction if older
func SkipNextIfOlder(in interface{}, out output.Output) error {
	out.InstructionTitle("Skip next if older")
	ok, err := older(in, out)
	if err != nil {
		return err
	}
	if ok {
		return quickonfErrors.SkipNext
	}
	return nil
}

func equal(in interface{}, out output.Output) (bool, error) {
	data, err := helper.SliceString(in)
	if err != nil {
		return false, err
	}
	for i := 1; i < len(data); i++ {
		if data[0] != data[i] {
			out.Info("Values are not equal")
			return false, nil
		}
	}
	out.Info("Values are equal")
	return true, nil
}

// StopIfEqual compares values and stops the step if all of them are equal
func StopIfEqual(in interface{}, out output.Output) error {
	out.InstructionTitle("Stop if equal")
	ok, err := equal(in, out)
	if err != nil {
		return err
	}
	if ok {
		return quickonfErrors.NoError
	}
	return nil
}

// SkipNextIfEqual compares values and skip the next instruction if all of them are equal
func SkipNextIfEqual(in interface{}, out output.Output) error {
	out.InstructionTitle("Skip next if equal")
	ok, err := equal(in, out)
	if err != nil {
		return err
	}
	if ok {
		return quickonfErrors.SkipNext
	}
	return nil
}
