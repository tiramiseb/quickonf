package main

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"os"
	"regexp"
	"text/template"

	"github.com/charmbracelet/glamour"
	"github.com/muesli/ansi/compressor"
	"github.com/tiramiseb/quickonf/commands"
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
	makeHelpCommands()
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

	f, err := os.Create(name + ".dark.msg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	g, _ := gzip.NewWriterLevel(f, gzip.BestCompression)
	defer g.Close()
	if _, err := g.Write(resp); err != nil {
		panic(err)
	}

	resp, err = lightRender.RenderBytes(content)
	if err != nil {
		panic(err)
	}
	resp = compressor.Bytes(resp)

	f, err = os.Create(name + ".light.msg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	g, _ = gzip.NewWriterLevel(f, gzip.BestCompression)
	defer g.Close()
	if _, err := g.Write(resp); err != nil {
		panic(err)
	}
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

func makeHelpCommands() {
	tmpl, err := template.ParseFiles("command.md.tmpl")
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	for _, cmd := range commands.GetAll() {
		buf.Reset()
		if err := tmpl.Execute(&buf, cmd); err != nil {
			panic(err)
		}
		makeUIFile("commands/"+cmd.Name, buf.Bytes())
	}
}
