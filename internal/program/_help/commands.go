package help

import (
	"fmt"
	"strings"
)

func (m *model) commandsDoc() string {
	content, ok := m.filteredCommandsDoc[m.filter]
	if !ok {
		var result strings.Builder
		for _, cmd := range m.commands {
			if !strings.Contains(cmd.Name, m.filter) {
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
		content = m.render(result.String())
		m.filteredCommandsDoc[m.filter] = content
	}
	return content
}
