package modules

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
)

func init() {
	// Register("browse-web", BrowseWeb)
	Register("parse-web-page", ParseWebPage)
}

// // BrowseWeb browses a website
// func BrowseWeb(in interface{}, out output.Output) error {
// 	out.InstructionTitle("Browse website")
// 	data, err := helper.MapStringString(in)
// 	if err != nil {
// 		return err
// 	}
// 	url, ok := data["url"]
// 	if !ok {
// 		return errors.New("Missing url")
// 	}
// 	driver := agouti.NewWebDriver("http://{{.Address}}", []string{"geckodriver", "--port={{.Port}}"})
// 	// driver := agouti.GeckoDriver()
// 	if err := driver.Start(); err != nil {
// 		return err
// 	}
// 	page, err := driver.NewPage(agouti.Browser("firefox"))
// 	if err != nil {
// 		return err
// 	}
// 	if err := page.Navigate(url); err != nil {
// 		return err
// 	}
// 	for {
// 		time.Sleep(500 * time.Millisecond)
// 		links := page.All("a")
// 		nb, err := links.Count()
// 		if err != nil {
// 			return err
// 		}
// 		if nb > 0 {
// 			fmt.Println("///")
// 			time.Sleep(10 * time.Second)
// 			fmt.Println("...")
// 			fmt.Println(links)
// 			for i := 0; i < nb; i++ {
// 				fmt.Println(links.At(i).Text())
// 			}
// 			break
// 		}
// 		// page.
// 	}
// 	return nil
// }

// ParseWebPage parses a web page...
func ParseWebPage(in interface{}, out output.Output) error {
	out.InstructionTitle("Parse web page")
	data, err := helper.MapStringString(in)
	if err != nil {
		return err
	}
	url, ok := data["url"]
	if !ok {
		return errors.New("Missing url")
	}

	reg, ok := data["regexp"]
	if !ok {
		return errors.New("Missing regexp")
	}

	page, err := helper.Download(url)
	if err != nil {
		return err
	}

	re, err := regexp.Compile(reg)
	if err != nil {
		return err
	}

	matches := re.FindStringSubmatch(string(page))

	if len(matches) == 0 {
		return errors.New("No match for " + re.String() + " in " + url)
	}

	out.Info("Match " + re.String())
	store, ok := data["store"]
	if ok {
		helper.Store(store, matches[0])
	}
	for i, name := range re.SubexpNames() {
		if name != "" {
			out.Info(fmt.Sprintf("%s value is %s", name, matches[i]))
			for key, val := range data {
				if key == "store-"+name {
					helper.Store(val, matches[i])
				}
			}
		}
	}

	// fmt.Println(string(page))

	return nil
}
