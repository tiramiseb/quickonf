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

const (
	fstabPath  = "/etc/fstab"
	mountsPath = "/proc/mounts"
)

func init() {
	register(mount)
}

var mount = &Command{
	"mount",
	"Add a mountpoint it to the fstab file and immediately mount the partition",
	[]string{
		"Absolute path to the partition",
		"Mount point absolute path",
		"Filesystem",
		"Options",
	},
	nil,
	"Mount data disk\n  file.directory /home/alice/data\n  mount /dev/sdb1 /home/alice/data ext4 defaults",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		partition := args[0]
		mountpoint := args[1]
		filesystem := args[2]
		options := args[3]
		if !filepath.IsAbs(partition) {
			return nil, fmt.Sprintf("%s is not an absolute path", partition), nil, StatusError, "", ""
		}
		if !filepath.IsAbs(mountpoint) {
			return nil, fmt.Sprintf("%s is not an absolute path", mountpoint), nil, StatusError, "", ""
		}
		if _, err := os.Lstat(partition); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Sprintf("%s does not exist", partition), nil, StatusError, "", ""
			}
			return nil, err.Error(), nil, StatusError, "", ""
		}
		finfo, err := os.Lstat(mountpoint)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Sprintf("%s does not exist", mountpoint), nil, StatusError, "", ""
			}
			return nil, err.Error(), nil, StatusError, "", ""
		}
		if !finfo.IsDir() {
			return nil, fmt.Sprintf("%s is not a directory", mountpoint), nil, StatusError, "", ""
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
			if fields[0] == partition && fields[1] == mountpoint {
				if fields[2] == filesystem && fields[3] == options {
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
			if fields[0] == partition && fields[1] == mountpoint {
				if fields[2] == filesystem && fields[3] == options {
					mounted = true
				} else {
					mountedAsExpected = true
				}
				break
			}
		}

		if inFstabAsExpected && mountedAsExpected {
			return nil, fmt.Sprintf("%s is mounted on %s", partition, mountpoint), nil, StatusSuccess, "", ""
		}

		newLine := fmt.Sprintf("%s %s %s %s", partition, mountpoint, filesystem, options)
		msgParts := []string{}
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
						if string(fields[0]) == partition && string(fields[1]) == mountpoint {
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
			out.Successf("Mounted %s on %s", partition, mountpoint)
			return true
		}

		return nil, "Need to " + strings.Join(msgParts, " and "), apply, StatusInfo, inFstabDifferent, newLine
	},
	nil,
}
