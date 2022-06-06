package commands

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tiramiseb/quickonf/internal/commands/helper"
)

func init() {
	register(gnomeExtensionInstall)
}

var gnomeExtensionInstall = Command{
	"gnome.extension.install",
	"Install a GNOME Shell extension from extensions.gnome.org",
	[]string{
		"Extension UUID",
	},
	nil,
	"Dash to dock\n  gnome.extension.install dash-to-dock@micxgx.gmail.com\n  gnome.extension.enable dash-to-dock@micxgx.gmail.com",
	func(args []string) (result []string, msg string, apply Apply, status Status, before, after string) {
		uuid := args[0]
		dest := filepath.Join("/usr/share/gnome-shell/extensions", uuid)

		// Version of GNOME Shell
		var buf bytes.Buffer
		helper.Exec(nil, &buf, "gnome-shell", "--version")
		gnomeVerLine := strings.Fields(buf.String())
		if len(gnomeVerLine) != 3 {
			return nil, fmt.Sprintf(`GNOME version invalid, should be "GNOME Shell X.Y.Z": %s`, buf.String()), nil, StatusError, "", ""
		}

		// Latest version of the extension
		gnomeVersion := string(gnomeVerLine[2])
		extInfo := struct {
			Version     int
			DownloadURL string `json:"download_url"`
		}{}
		if err := helper.DownloadJSON("https://extensions.gnome.org/extension-info/?uuid="+uuid+"&shell_version="+gnomeVersion, &extInfo); err != nil {
			return nil, fmt.Sprintf("Could not get information about extension %s: %s", uuid, err), nil, StatusError, "", ""
		}

		// Current installed version
		var localVersion int
		localMetadataFile, err := os.Open(filepath.Join(dest, "metadata.json"))
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Sprintf("Could not query local extension %s: %s", uuid, err), nil, StatusError, "", ""
			}
		} else {
			extMeta := struct {
				Version int
			}{}
			jsondec := json.NewDecoder(localMetadataFile)
			if err := jsondec.Decode(&extMeta); err != nil {
				localMetadataFile.Close()
				return nil, fmt.Sprintf("Could not read local extension %s data: %s", uuid, err), nil, StatusError, "", ""
			}
			localMetadataFile.Close()
			localVersion = extMeta.Version
		}
		if extInfo.Version <= localVersion {
			return nil, fmt.Sprintf("%s already installed in version %d", uuid, localVersion), nil, StatusSuccess, strconv.Itoa(localVersion), strconv.Itoa(extInfo.Version)
		}

		apply = func(out Output) (success bool) {
			out.Runningf("Downloading %s", uuid)
			resp, err := http.Get("https://extensions.gnome.org" + extInfo.DownloadURL)
			if err != nil {
				out.Errorf("Could not download %s: %s", uuid, err)
				return false
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				out.Errorf("Could not read downloaded data for %s: %s", uuid, err)
				return false
			}
			out.Runningf("Extracting %s", uuid)
			reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
			if err != nil {
				out.Errorf("Could not extract downloaded data for %s: %s", uuid, err)
				return false
			}
			for _, f := range reader.File {
				fpath := filepath.Join(dest, f.Name)

				// Zip Slip vulnerability
				if !strings.HasPrefix(fpath, dest) {
					out.Errorf("Could not extract file %s for %s: illegal path", f.Name, uuid)
					return false
				}

				if f.FileInfo().IsDir() {
					if err := os.MkdirAll(fpath, 0o755); err != nil {
						out.Errorf("Could not create %s: %s", fpath, err)
						return false
					}
					continue
				}

				if err = os.MkdirAll(filepath.Dir(fpath), 0o755); err != nil {
					out.Errorf("Could not create directory for %s: %s", fpath, err)
					return false
				}

				dest, err := os.Create(fpath)
				if err != nil {
					out.Errorf("Could not create %s: %s", fpath, err)
					return false
				}
				src, err := f.Open()
				if err != nil {
					dest.Close()
					out.Errorf("Could not read %s in archive: %s", f.Name, err)
					return false
				}
				if _, err := io.Copy(dest, src); err != nil {
					dest.Close()
					src.Close()
					out.Errorf("Could not write %s: %s", fpath, err)
					return false
				}
				dest.Close()
				src.Close()
			}
			out.Successf("Installed %s", uuid)
			return true
		}
		return nil, fmt.Sprintf("Need to install %s", uuid), apply, StatusInfo, strconv.Itoa(localVersion), strconv.Itoa(extInfo.Version)
	},
	nil,
}
