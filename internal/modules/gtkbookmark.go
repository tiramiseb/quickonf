package modules

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

const gtkBookmarkPath = ".config/gtk-3.0/bookmarks"

func init() {
	Register("gtk-bookmarks", GtkBookmarks)
}

// GtkBookmarks sets the Gtk bookmarks
func GtkBookmarks(in interface{}, out output.Output) error {
	out.InstructionTitle("Set Gtk Bookmarks")
	data, err := helper.SliceString(in)
	if err != nil {
		return err
	}
	var w io.Writer
	if Dryrun {
		w = ioutil.Discard

	} else {
		f, err := os.Create(helper.Path(gtkBookmarkPath))
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}
	result := make([][2]string, len(data))
	for i, bookmark := range data {
		splat := strings.SplitN(bookmark, "=", 2)
		if len(splat) == 0 {
			continue
		}
		name := ""
		path := ""
		if len(splat) == 1 {
			path = splat[0]
		} else {
			name = splat[0]
			path = splat[1]
		}
		if !strings.Contains(path, "://") {
			path = "file://" + helper.Path(path)
		}
		if _, err := w.Write([]byte(path)); err != nil {
			return err
		}
		if name != "" {
			if _, err := w.Write([]byte(" " + name)); err != nil {
				return err
			}
		} else {
			name = filepath.Base(path)
		}
		result[i] = [2]string{name, path}
		if _, err := w.Write([]byte{'\n'}); err != nil {
			return err
		}
	}
	for _, d := range result {
		if Dryrun {
			out.Infof("Would set bookmark %s → %s", d[0], d[1])
			continue
		}
		out.Successf("Bookmark %s → %s", d[0], d[1])
	}

	return nil
}
