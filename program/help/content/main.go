package main

//go:generate go run .
import (
	"bytes"
	"os"
	"regexp"
	"text/template"

	"github.com/charmbracelet/glamour"
	"github.com/gosimple/slug"
	"github.com/muesli/ansi/compressor"
	"github.com/tiramiseb/quickonf/commands"
	"github.com/tiramiseb/quickonf/embeddedcookbook"
	"github.com/tiramiseb/quickonf/instructions"
)

var (
	darkStyle  = glamour.DarkStyleConfig
	lightStyle = glamour.LightStyleConfig

	imageRe       = regexp.MustCompile(`!\[[^]]*\]\([^)]*\)`)
	frontmatterRe = regexp.MustCompile("(?s)---\n.*\n---")
)

func main() {
	darkStyle.Document.Margin = nil
	lightStyle.Document.Margin = nil

	// Static content pages
	md, err := os.ReadFile("../../../docs/content/intro.md")
	if err != nil {
		panic(err)
	}
	makeUIFile("intro", md)

	md, err = os.ReadFile("../../../docs/content/language.md")
	if err != nil {
		panic(err)
	}
	makeUIFile("language", md)

	md, err = os.ReadFile("../../../docs/content/ui.md")
	if err != nil {
		panic(err)
	}
	makeUIFile("ui", md)

	// Commands
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

	// Cookbook
	tmpl, err = template.ParseFiles("recipe.md.tmpl")
	if err != nil {
		panic(err)
	}
	if err := embeddedcookbook.ForEach(func(recipe *instructions.Group) error {
		buf.Reset()
		if err := tmpl.Execute(&buf, recipe); err != nil {
			return err
		}
		makeUIFile("cookbook/"+slug.Make(recipe.Name), buf.Bytes())
		return nil
	}); err != nil {
		panic(err)
	}
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
