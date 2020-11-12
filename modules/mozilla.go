package modules

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"

	"github.com/tiramiseb/quickonf/modules/input"
	"github.com/tiramiseb/quickonf/output"
)

func init() {
	Register("firefox-extension", FirefoxExtension)
}

func FirefoxExtension(in interface{}, out output.Output, store map[string]interface{}) error {
	out.ModuleName("Install Firefox extension")
	data, err := input.SliceString(in, store)
	if err != nil {
		return err
	}
	// cfg, err := ini.Load(helper.Path(".mozilla/firefox/profiles.ini"))
	// if err != nil {
	// 	return err
	// }
	// path := cfg.Section("Profile0").Key("Path").String()
	for _, extension := range data {
		resp, err := http.Get("https://addons.mozilla.org/firefox/addon/" + extension)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		page := html.NewTokenizer(resp.Body)
		var url string
		for {
			tokenType := page.Next()
			if tokenType == html.ErrorToken {
				err := page.Err()
				if err != io.EOF {
					return err
				}
				break
			}
			if tokenType != html.StartTagToken {
				continue
			}
			tagName, hasAttr := page.TagName()
			if string(tagName) != "a" || !hasAttr {
				continue
			}
			url = getFirefoxExtensionURL(page)
			if url == "" {
				continue
			}
			fmt.Print("---")
		}
	}
	return nil
}

// MARCHE PAS :'(
func getFirefoxExtensionURL(page *html.Tokenizer) string {
	for {
		key, val, hasMore := page.TagAttr()
		if !bytes.Equal(key, []byte("class")) {
			return ""
		}
		// fmt.Print(".")
		fmt.Println(string(val))
		// if bytes.Contains(val, []byte("AMInstallButton-button")) {
		// 	fmt.Println(string(val))
		// }
		if !hasMore {
			break
		}
	}
	return ""
}
