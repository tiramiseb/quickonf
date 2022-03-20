package details

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/instructions"
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
	if global.Toggles["details"] {
		for _, rep := range group.Reports {
			view += m.reportView(rep)
			if rep.Status == commands.StatusInfo {
				view += m.detailsView(rep)
			}
		}
	} else {
		for _, rep := range group.Reports {
			view += m.reportView(rep)
		}
	}
	m.viewport.SetContent(view)
	return m.viewport.View()
}

func (m *Model) reportView(rep *instructions.CheckReport) string {
	content := fmt.Sprintf("[%s] %s", rep.Name, rep.Message)
	return global.Styles[rep.Status].Render(
		global.MakeWidth(content, m.width),
	) + "\n"
}

func (m *Model) detailsView(rep *instructions.CheckReport) string {
	switch {
	case rep.Before == "" && rep.After == "":
		return ""
	case rep.Before != "" && rep.After != "" && rep.Before != rep.After:
		leftWidth := (m.width - 1) / 2
		rightWidth := m.width - leftWidth - 1
		left := lipgloss.NewStyle().Width(leftWidth).Render(rep.Before)
		right := lipgloss.NewStyle().Width(rightWidth).Render(rep.After)
		height := lipgloss.Height(left)
		if rightHeight := lipgloss.Height(right); rightHeight > height {
			height = rightHeight
		}
		var separator string
		for i := 0; i < height-1; i++ {
			separator += "│\n"
		}
		separator += "│"
		return lipgloss.JoinHorizontal(lipgloss.Top, left, separator, right) + "\n"
	case rep.Before == "":
		return lipgloss.NewStyle().Width(m.width).Render(rep.After) + "\n"
	default:
		return lipgloss.NewStyle().Width(m.width).Render(rep.Before) + "\n"
	}
}
