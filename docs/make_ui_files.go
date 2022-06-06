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
	//go:embed content/intro.md
	intro []byte
	//go:embed content/language.md
	language []byte
	//go:embed content/ui.md
	ui []byte

	darkStyle  = glamour.DarkStyleConfig
	lightStyle = glamour.LightStyleConfig

	imageRe       = regexp.MustCompile(`!\[[^]]*\]\([^)]*\)`)
	frontmatterRe = regexp.MustCompile("(?s)---\n.*\n---")
)

func init() {
	darkStyle.Document.Margin = nil
	lightStyle.Document.Margin = nil
}

func makeUIFiles() {
	makeUIFile("intro", intro)
	makeUIFile("language", language)
	makeUIFile("ui", ui)
}

func makeUIFile(name string, content []byte) {
	content = imageRe.ReplaceAll(content, []byte{})
	content = frontmatterRe.ReplaceAll(content, []byte{})
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
