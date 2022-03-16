package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const aptKeysBase = "/etc/apt/trusted.gpg.d/"

func init() {
	register(aptKey)
}

var aptKey = Command{
	"apt.key",
	"Add apt key",
	[]string{
		"Local name of the key (simple short text)",
		"Key",
	},
	nil,
	"NextDNS\n  key = http.get.var https://repo.nextdns.io/nextdns.gpg\n  apt.key nextdns <key>",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		name := args[0]
		key := args[1]
		var ext string
		if strings.HasPrefix(key, "-----BEGIN PGP PUBLIC KEY BLOCK-----") {
			ext = "asc"
		} else {
			ext = "gpg"
		}

		keyFile := filepath.Join(aptKeysBase, name+"."+ext)
		existingB, err := os.ReadFile(keyFile)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Sprintf("Could not read existing key file: %s", err), nil, StatusError, "", ""
		}
		existing := string(existingB)
		if existing == key {
			return nil, fmt.Sprintf("Key %s already known", name), nil, StatusSuccess, "", ""
		}
		apply = func(out Output) bool {
			out.Runningf("Storing the apt key %s", name)
			if err := os.WriteFile(keyFile, []byte(key), 0o644); err != nil {
				out.Errorf("Could not write requested content to %s: %s", keyFile, err)
				return false
			}
			out.Successf("Apt key %s added", key)
			return true
		}
		return nil, fmt.Sprintf("Need to add apt key %s", name), apply, StatusInfo, "", key
	},
	nil,
}
