package details

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/program/global"
)

func (m *Model) View() string {
	var view string
	group := global.GetSelectedGroup()
	if group == nil {
		return ""
	}
	if len(group.Reports) == 0 {
		for _, ins := range group.Instructions {
			view += global.Styles[commands.StatusNotRun].Render(
				global.MakeWidth(ins.Name(), m.width),
			) + "\n"
		}
	}
	for _, rep := range group.Reports {
		content := fmt.Sprintf("[%s] %s", rep.Name, rep.Message)
		view += global.Styles[rep.Status].Render(
			global.MakeWidth(content, m.width),
		) + "\n"
	}
	m.viewport.SetContent(view)
	return m.viewport.View()
}
