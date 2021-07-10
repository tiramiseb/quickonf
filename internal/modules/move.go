package modules

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

type moveDestinationExists int

const (
	moveErase moveDestinationExists = iota
	moveFail
	movePass
)

var moveMigrationSource = ""

func init() {
	Register("move", Move)
	Register("force-move", ForceMove)
	Register("migration-source", MigrationSource)
	Register("force-migrate", ForceMigrate)
	Register("migrate", Migrate)
	Register("do-not-migrate", DoNotMigrate)
}

// Move moves files or directories, or does nothing if the source does not exist
func Move(in interface{}, out output.Output) error {
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	out.InstructionTitle("Move file or directory")
	return move(data, out, moveFail)
}

// ForceMove moves files or directories, removing the destination if it exists
func ForceMove(in interface{}, out output.Output) error {
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	out.InstructionTitle("Move file or directory, crushing destination if necessary")
	return move(data, out, moveErase)
}

// MigrationSource sets the migration source path
func MigrationSource(in interface{}, out output.Output) error {
	out.InstructionTitle("Set migration source")
	data, err := helper.String(in)
	if err != nil {
		return err
	}
	moveMigrationSource = helper.Path(data)
	out.Successf("Migration source is %s", moveMigrationSource)
	return nil
}

// Migrate migrates a file or directory from the previous home
func Migrate(in interface{}, out output.Output) error {
	return migrate(in, out, movePass)
}

// ForceMigrate migrates a file or directory from the previous home, removing the destination if it exists
func ForceMigrate(in interface{}, out output.Output) error {
	return migrate(in, out, moveErase)
}

// Migrate migrates a file or directory from the previous home
func migrate(in interface{}, out output.Output, ifExists moveDestinationExists) error {
	out.InstructionTitle("Migrate file or directory")
	if moveMigrationSource == "" {
		return errors.New("migration source is not defined")
	}
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	translated := map[string]string{}
	for _, path := range data {
		translated[filepath.Join(moveMigrationSource, path)] = path
	}
	return move(translated, out, ifExists)
}

// DoNotMigrate removes a file fro the previous home, without migrating it
func DoNotMigrate(in interface{}, out output.Output) error {
	if moveMigrationSource == "" {
		return errors.New("migration source is not defined")
	}
	out.InstructionTitle("Remove do-no-migrate files")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	for _, fpath := range data {
		fpath = filepath.Join(moveMigrationSource, fpath)
		if Dryrun {
			out.Infof("Would make sure %s does not exist", fpath)
			continue
		}
		out.Infof("Making sure %s does not exist", fpath)
		if err := os.RemoveAll(fpath); err != nil {
			return err
		}
	}
	return nil
}

func move(data map[string]string, out output.Output, destExists moveDestinationExists) error {
	for from, to := range data {
		from = helper.Path(from)
		to = helper.Path(to)
		if _, err := os.Stat(from); err != nil {
			if os.IsNotExist(err) {
				out.Infof("Source %s does not exist", from)
				continue
			}
			return err
		}
		_, err := os.Stat(to)
		if err == nil {
			if destExists == moveFail {
				return fmt.Errorf("%s already exists", to)
			}
			if destExists == movePass {
				out.Infof("%s already exists", to)
				continue
			}
			if Dryrun {
				out.Info("Would remove " + to)
			} else {
				if err := os.RemoveAll(to); err != nil {
					return err
				}
			}
		} else {
			if !os.IsNotExist(err) {
				return err
			}
		}
		if Dryrun {
			out.Infof("Would move %s to %s", from, to)
			continue
		}
		if err = os.Rename(from, to); err != nil {
			return err
		}
		out.Successf("Moved %s to %s", from, to)
	}
	return nil
}
