package main

import (
	"bytes"
	_ "embed"
	"os"
	"regexp"

	"github.com/charmbracelet/glamour"
	"github.com/muesli/ansi/compressor"
)

var (
	//go:embed intro.md
	intro []byte
	//go:embed language.md
	language []byte
	//go:embed ui.md
	ui []byte

	darkStyle  = glamour.DarkStyleConfig
	lightStyle = glamour.LightStyleConfig

	spacesRe = regexp.MustCompile("  +")
)

func init() {
	darkStyle.Document.Margin = nil
	lightStyle.Document.Margin = nil
}

func main() {
	make("intro", intro)
	make("language", language)
	make("ui", ui)
}

func make(name string, content []byte) {
	width := maxWidth(content)

	darkRender, err := glamour.NewTermRenderer(
		glamour.WithStyles(darkStyle),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		panic(err)
	}
	lightRender, err := glamour.NewTermRenderer(
		glamour.WithStyles(lightStyle),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		panic(err)
	}

	resp, err := darkRender.RenderBytes(content)
	if err != nil {
		panic(err)
	}
	resp = compressor.Bytes(resp)
	os.WriteFile(name+".dark.msg", resp, 0o644)

	resp, err = lightRender.RenderBytes(content)
	if err != nil {
		panic(err)
	}
	resp = compressor.Bytes(resp)
	os.WriteFile(name+".light.msg", resp, 0o644)
}

func maxWidth(content []byte) int {
	var max int
	lines := bytes.Split(content, []byte("\n"))
	for _, l := range lines {
		if len(l) > max {
			max = len(l)
		}
	}
	return max
}
