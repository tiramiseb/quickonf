package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tiramiseb/quickonf/commands/helper"
)

func init() {
	register(mountBind)
}

var mountBind = &Command{
	"mount.bind",
	"Add a bind mountpoint it to the fstab file and immediately mount the partition",
	[]string{
		"Absolute path to the target",
		"Mount point absolute path",
	},
	nil,
	"Mount directory elsewhere\n  mount.bind /tmp /home/alice/temp",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		target := args[0]
		mountpoint := args[1]
		if !filepath.IsAbs(target) {
			return nil, fmt.Sprintf("%s is not an absolute path", target), nil, StatusError, "", ""
		}
		if !filepath.IsAbs(mountpoint) {
			return nil, fmt.Sprintf("%s is not an absolute path", mountpoint), nil, StatusError, "", ""
		}
		var mustCreateTarget bool
		if _, err := os.Lstat(target); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				mustCreateTarget = true
			} else {
				return nil, err.Error(), nil, StatusError, "", ""
			}
		}

		var mustCreateMountpoint bool
		finfo, err := os.Lstat(mountpoint)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				mustCreateMountpoint = true
			} else {
				return nil, err.Error(), nil, StatusError, "", ""
			}
		} else {
			if !finfo.IsDir() {
				return nil, fmt.Sprintf("%s is not a directory", mountpoint), nil, StatusError, "", ""
			}
		}

		var inFstabAsExpected bool
		var inFstabDifferent string

		f, err := os.Open(fstabPath)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if strings.HasPrefix(fields[0], "#") {
				continue
			}
			if fields[0] == target && fields[1] == mountpoint {
				if fields[2] == "bind" && fields[3] == "defaults" {
					inFstabAsExpected = true
				} else {
					inFstabDifferent = line
				}
				break
			}
		}
		f.Close()
		if err := scanner.Err(); err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}

		var mounted bool
		var mountedAsExpected bool

		f, err = os.Open(mountsPath)
		if err != nil {
			return nil, err.Error(), nil, StatusError, "", ""
		}
		scanner = bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if fields[0] == target && fields[1] == mountpoint {
				if fields[2] == "bind" && fields[3] == "defaults" {
					mounted = true
				} else {
					mountedAsExpected = true
				}
				break
			}
		}

		if inFstabAsExpected && mountedAsExpected {
			return nil, fmt.Sprintf("%s is mounted on %s", target, mountpoint), nil, StatusSuccess, "", ""
		}

		newLine := fmt.Sprintf("%s %s bind defaults", target, mountpoint)
		msgParts := []string{}
		if mustCreateTarget {
			msgParts = append(msgParts, "create target "+target)
		}
		if mustCreateMountpoint {
			msgParts = append(msgParts, "create mountpoint "+mountpoint)
		}
		if inFstabDifferent != "" {
			msgParts = append(msgParts, "change the content of "+fstabPath)
		} else if !inFstabAsExpected {
			msgParts = append(msgParts, "add an entry to "+fstabPath)
		}
		if mounted && !mountedAsExpected {
			msgParts = append(msgParts, "unmount and remount "+mountpoint)
		} else if !mounted {
			msgParts = append(msgParts, "mount "+mountpoint)
		}

		apply = func(out Output) bool {
			if mustCreateTarget {
				out.Runningf("Creating target %s", target)
				if err := os.MkdirAll(target, 0o755); err != nil {
					out.Errorf("Could not create target %s: %s", target, err)
					return false
				}
			}
			if mustCreateMountpoint {
				out.Runningf("Creating mountpoint %s", mountpoint)
				if err := os.MkdirAll(mountpoint, 0o755); err != nil {
					out.Errorf("Could not create mountpoint %s: %s", mountpoint, err)
					return false
				}
			}
			if !inFstabAsExpected {
				if inFstabDifferent == "" {
					out.Runningf("Adding entry to %s", fstabPath)
					f, err := os.OpenFile(fstabPath, os.O_APPEND, 0o644)
					if err != nil {
						out.Errorf("Could not open %s: %s", fstabPath, err)
						return false
					}
					f.WriteString(newLine)
					f.Write([]byte{'\n'})
					f.Close()
				} else {
					out.Runningf("Modifying entry in %s", fstabPath)
					src, err := os.Open(fstabPath)
					if err != nil {
						out.Errorf("Could not open %s: %s", fstabPath, err)
						return false
					}
					var buf bytes.Buffer
					wroteIt := false
					scanner := bufio.NewScanner(src)
					for scanner.Scan() {
						line := scanner.Bytes()
						fields := bytes.Fields(line)
						if string(fields[0]) == target && string(fields[1]) == mountpoint {
							if !wroteIt {
								buf.WriteString(newLine)
								buf.WriteByte('\n')
								wroteIt = true
							}
						} else {
							buf.Write(line)
							buf.WriteByte('\n')
						}
					}
					src.Close()
					if err := scanner.Err(); err != nil {
						out.Errorf("Could not read %s: %s", fstabPath, err)
						return false
					}
					if err := os.WriteFile(fstabPath, buf.Bytes(), finfo.Mode()); err != nil {
						out.Errorf("Could not write to %s: %s", fstabPath, err)
						return false
					}
				}

				if !mountedAsExpected {
					if err := helper.Exec(nil, nil, "umount", mountpoint); err != nil {
						out.Errorf("Could not unmount %s: %s", mountpoint, helper.ExecErr(err))
						return false
					}
				}
				if !mounted {
					if err := helper.Exec(nil, nil, "mount", mountpoint); err != nil {
						out.Errorf("Could not mount %s: %s", mountpoint, helper.ExecErr(err))
						return false
					}
				}
			}
			out.Successf("Mounted %s on %s", target, mountpoint)
			return true
		}

		return nil, "Need to " + strings.Join(msgParts, " and "), apply, StatusInfo, inFstabDifferent, newLine
	},
	nil,
}
