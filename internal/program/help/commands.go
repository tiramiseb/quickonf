package help

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
)

var (
	darkStyle  = glamour.DarkStyleConfig
	lightStyle = glamour.LightStyleConfig
)

func init() {
	darkStyle.Document.Margin = nil
	lightStyle.Document.Margin = nil
}

func (m *Model) commandsDoc(dark bool) string {
	content, ok := m.filteredCommandsDoc[m.commandFilter]
	if !ok {
		var result strings.Builder
		for _, cmd := range m.commands {
			if !strings.Contains(cmd.Name, m.commandFilter) {
				continue
			}
			fmt.Fprint(&result, "\n# ", cmd.Name, "\n\n", cmd.Action, "\n")
			if len(cmd.Arguments) > 0 {
				fmt.Fprintln(&result, "\nArguments:")
				for _, arg := range cmd.Arguments {
					fmt.Fprintln(&result, "* ", arg)
				}
			}
			if len(cmd.Outputs) > 0 {
				fmt.Fprintln(&result, "\nOutputs:")
				for _, out := range cmd.Outputs {
					fmt.Fprintln(&result, "* ", out)
				}
			}
			if len(cmd.Example) > 0 {
				fmt.Fprintln(&result, "\nExample:\n\n```")
				fmt.Fprintln(&result, cmd.Example, "\n```")
			}
		}
		var renderer *glamour.TermRenderer
		var err error
		if dark {
			renderer, err = glamour.NewTermRenderer(
				glamour.WithStyles(darkStyle),
				glamour.WithWordWrap(m.width),
			)
		} else {
			renderer, err = glamour.NewTermRenderer(
				glamour.WithStyles(lightStyle),
				glamour.WithWordWrap(m.width),
			)
		}
		if err != nil {
			panic(err)
		}
		content, err = renderer.Render(result.String())
		if err != nil {
			content = "Could not render documentation: " + err.Error() + "\n-----\n" + content
		}
		m.filteredCommandsDoc[m.commandFilter] = content
	}
	return content
}
